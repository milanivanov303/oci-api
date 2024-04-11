package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"context"
	"fmt"
)

func (s *APIServer) handleApiKeys(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetApiKeys(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetApiKeys(w http.ResponseWriter, r *http.Request) error {
	users, err := listUsers(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing users:", err)
        return err
    }

	var allApiKeys []identity.ApiKey
	for _, user := range users {
		apiKeys, err := listApiKeys(s.IdentityClient, *user.Id)
		if err != nil {
			fmt.Printf("Error listing API keys for user %s: %v\n", *user.Name, err)
			continue
		}
		allApiKeys = append(allApiKeys, apiKeys...)
	}

	return WriteJSON(w, http.StatusOK, allApiKeys)
}

func listApiKeys(identityClient identity.IdentityClient, userID string) ([]identity.ApiKey, error) {
    request := identity.ListApiKeysRequest{
        UserId: &userID,
    }

    response, err := identityClient.ListApiKeys(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}