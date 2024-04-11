package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"context"
	"fmt"
)

func (s *APIServer) handleAuthTokens(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAuthTokens(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAuthTokens(w http.ResponseWriter, r *http.Request) error {
	users, err := listUsers(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing users:", err)
        return err
    }

	var allAuthTokens []identity.AuthToken
	for _, user := range users {
		authTokens, err := listAuthTokens(s.IdentityClient, *user.Id)
		if err != nil {
			fmt.Printf("Error listing auth tokens for user %s: %v\n", *user.Name, err)
			continue
		}
		allAuthTokens = append(allAuthTokens, authTokens...)
	}

	return WriteJSON(w, http.StatusOK, allAuthTokens)
}

func listAuthTokens(identityClient identity.IdentityClient, userID string) ([]identity.AuthToken, error) {
    request := identity.ListAuthTokensRequest{
        UserId: &userID,
    }

    response, err := identityClient.ListAuthTokens(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}