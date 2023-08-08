package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/core"
	"io/ioutil"
	"context"
)

// Server
type APIServer struct {
	listenAddr string
	store      Storage
	config     common.ConfigurationProvider
	client     identity.IdentityClient
	tenancyID  string
}

// to implement db add store Storage
func NewAPIServer(listenAddr string, configuration Configuration) (*APIServer, error) {
	config, err := getConfig(configuration)
	if err != nil {
		return nil, err
	}

	client, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		return nil, err
	}

	tenancyID, err := config.TenancyOCID()
	if err != nil {
		return nil, err
	}

	return &APIServer{
		listenAddr: listenAddr,
		config:     config,
		client:     client,
		tenancyID:  tenancyID,
	}, nil
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/compartments", makeHTTPHandleFunc(s.handleCompartments))
	router.HandleFunc("/instances/{compartment_id}", makeHTTPHandleFunc(s.handleInstances))
	// router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON APi server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func getConfig(configuration Configuration) (common.ConfigurationProvider, error) {
	fingerprint, err := ioutil.ReadFile(configuration.FingerprintPath)
    if err != nil {
        return nil, err
    }

    privateKey, err := ioutil.ReadFile(configuration.PrivateKeyPath)
    if err != nil {
        return nil, err
    }

	provider := common.NewRawConfigurationProvider(configuration.TenancyID, configuration.UserID, configuration.Region, string(fingerprint), string(privateKey), nil)

    return provider, nil
}

// Instances
func (s *APIServer) handleInstances(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetInstances(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetInstances(w http.ResponseWriter, r *http.Request) error {
	compartmentID := mux.Vars(r)["compartment_id"]

	instances, err := listInstances(s.config, compartmentID)
    if err != nil {
        fmt.Println("Error listing instances:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, instances)
}

func listInstances(config common.ConfigurationProvider, compartmentID string) ([]core.Instance, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(config)
    if err != nil {
        return nil, err
    }

    request := core.ListInstancesRequest{
        CompartmentId: &compartmentID,
    }

    response, err := client.ListInstances(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}  

// Compartments
func (s *APIServer) handleCompartments(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetCompartments(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetCompartments(w http.ResponseWriter, r *http.Request) error {
	compartments, err := listCompartments(s.client, s.tenancyID)
    if err != nil {
        fmt.Println("Error listing compartments:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, compartments)
}

func listCompartments(client identity.IdentityClient, compartmentID string) ([]identity.Compartment, error) {
    request := identity.ListCompartmentsRequest{
        CompartmentId: &compartmentID,
    }

    response, err := client.ListCompartments(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}

// Accounts
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Helper functions
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}