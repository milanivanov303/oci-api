package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"fmt"
	"context"
)

func (s *APIServer) handleDynamicGroups(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetDynamicGroups(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetDynamicGroups(w http.ResponseWriter, r *http.Request) error {
	groups, err := listDynamicGroups(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing dynamic groups:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, groups)
}

func listDynamicGroups(identityClient identity.IdentityClient, compartmentID string) ([]identity.DynamicGroup, error) {
    request := identity.ListDynamicGroupsRequest{
        CompartmentId: &compartmentID,
    }

    response, err := identityClient.ListDynamicGroups(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}