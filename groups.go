package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"fmt"
	"context"
)

func (s *APIServer) handleGroups(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetGroups(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetGroups(w http.ResponseWriter, r *http.Request) error {
	groups, err := listGroups(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing groups:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, groups)
}

func listGroups(identityClient identity.IdentityClient, compartmentID string) ([]identity.Group, error) {
    request := identity.ListGroupsRequest{
        CompartmentId: &compartmentID,
    }

    response, err := identityClient.ListGroups(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}