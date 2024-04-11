package main

import (
	"net/http"
	"fmt"
	"context"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUsers(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := listUsers(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing users:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, users)
}

func listUsers(identityClient identity.IdentityClient, compartmentID string) ([]identity.User, error) {
    request := identity.ListUsersRequest{
        CompartmentId: &compartmentID,
    }

    response, err := identityClient.ListUsers(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}