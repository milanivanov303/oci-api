package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"fmt"
	"context"
)

func (s *APIServer) handleCompartments(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetCompartments(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetCompartments(w http.ResponseWriter, r *http.Request) error {
	compartments, err := listCompartments(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing compartments:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, compartments)
}

func listCompartments(identityClient identity.IdentityClient, compartmentID string) ([]identity.Compartment, error) {
    request := identity.ListCompartmentsRequest{
        CompartmentId: &compartmentID,
    }

    response, err := identityClient.ListCompartments(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}