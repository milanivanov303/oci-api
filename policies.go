package main

import (
	"net/http"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"context"
)

func (s *APIServer) handlePolicies(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPolicies(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetPolicies(w http.ResponseWriter, r *http.Request) error {
	policies, err := listPolicies(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing policies:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, policies)
}

func listPolicies(identityClient identity.IdentityClient, compartmentID string) ([]identity.Policy, error) {
	request := identity.ListPoliciesRequest{
		CompartmentId: &compartmentID,
	}

	response, err := identityClient.ListPolicies(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}