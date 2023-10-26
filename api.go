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
	jwt "github.com/golang-jwt/jwt/v4"
)

// to implement db add store Storage
func NewAPIServer(listenAddr string, configuration Configuration) (*APIServer, error) {
	config, err := getConfig(configuration)
	if err != nil {
		return nil, err
	}

	identityClient, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		return nil, err
	}

	tenancyID, err := config.TenancyOCID()
	if err != nil {
		return nil, err
	}

	regionSubscriptions, err := listRegions(identityClient, configuration.TenancyID)
	if err != nil {
		fmt.Println("Error getting region subscriptions:", err)
        return nil, err
	}

    var computeClients []core.ComputeClient
	for _, region := range regionSubscriptions {
		regionConfig, err := getConfigForRegion(configuration, *region.RegionName)
		if err != nil {
			return nil, err
		}

		client, err := core.NewComputeClientWithConfigurationProvider(regionConfig)
		if err != nil {
			return nil, err
		}

		computeClients = append(computeClients, client)
	}

	return &APIServer{
		ListenAddr:     listenAddr,
		Config:         config,
		ComputeClients: computeClients,
		IdentityClient: identityClient,
		TenancyID:      tenancyID,
		Configuration:  configuration,
	}, nil
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/compartments", withJWTAuth(makeHTTPHandleFunc(s.handleCompartments), s.Configuration))
	router.HandleFunc("/instances/{compartment_id}", withJWTAuth(makeHTTPHandleFunc(s.handleInstances), s.Configuration))
	router.HandleFunc("/regions", withJWTAuth(makeHTTPHandleFunc(s.handleRegions), s.Configuration))
	// router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON APi server running on port: ", s.ListenAddr)

	http.ListenAndServe(s.ListenAddr, router)
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

	provider := common.NewRawConfigurationProvider(configuration.TenancyID, configuration.UserID, configuration.HomeRegion, string(fingerprint), string(privateKey), nil)

    return provider, nil
}

func getConfigForRegion(configuration Configuration, region string) (common.ConfigurationProvider, error) {
	fingerprint, err := ioutil.ReadFile(configuration.FingerprintPath)
    if err != nil {
        return nil, err
    }

    privateKey, err := ioutil.ReadFile(configuration.PrivateKeyPath)
    if err != nil {
        return nil, err
    }

	provider := common.NewRawConfigurationProvider(configuration.TenancyID, configuration.UserID, region, string(fingerprint), string(privateKey), nil)

    return provider, nil
}

// Helper functions
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt" : 9999999999,
		"type"      : "auth",
		"user"      : "ea_auto",
	}

	secret := "somestupidsecret123"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, configuration Configuration) http.HandlerFunc {
	secret := configuration.JWTSecret

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString, secret)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if (claims["user"] != "ea_auto") || (claims["type"] != "auth") {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// Accounts; not used; example for db implementation
// func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method == "GET" {
// 		return s.handleGetAccount(w, r)
// 	}
// 	if r.Method == "POST" {
// 		return s.handleCreateAccount(w, r)
// 	}
// 	if r.Method == "DELETE" {
// 		return s.handleDeleteAccount(w, r)
// 	}

// 	return fmt.Errorf("method not allowed %s", r.Method)
// }

// func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
// 	accounts, err := s.Store.GetAccounts()
// 	if err != nil {
// 		return err
// 	}

// 	return WriteJSON(w, http.StatusOK, accounts)
// }

// func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
// 	id := mux.Vars(r)["id"]

// 	fmt.Println(id)

// 	return WriteJSON(w, http.StatusOK, &Account{})
// }

// func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
// 	createAccountReq := new(CreateAccountRequest)
// 	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
// 		return err
// 	}

// 	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
// 	if err := s.Store.CreateAccount(account); err != nil {
// 		return err
// 	}

// 	// tokenString, err := createJWT(account)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// fmt.Println("JWT token: ", tokenString)

// 	return WriteJSON(w, http.StatusOK, account)
// }

// func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }