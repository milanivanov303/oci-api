package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"fmt"
	"context"
)

func (s *APIServer) handleRegions(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetRegions(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetRegions(w http.ResponseWriter, r *http.Request) error {
	regions, err := listRegions(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing regions:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, regions)
}

func listRegions(identityClient identity.IdentityClient, compartmentID string) ([]identity.RegionSubscription, error) {
    request := identity.ListRegionSubscriptionsRequest{
        TenancyId: &compartmentID,
    }

    response, err := identityClient.ListRegionSubscriptions(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}