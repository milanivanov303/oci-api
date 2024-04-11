package main

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	"net/http"
	"context"
	"fmt"
)

func (s *APIServer) handleSmtpCredentials(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetSmtpCredentials(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetSmtpCredentials(w http.ResponseWriter, r *http.Request) error {
	users, err := listUsers(s.IdentityClient, s.TenancyID)
    if err != nil {
        fmt.Println("Error listing users:", err)
        return err
    }

	var allCredentials []identity.SmtpCredentialSummary
	for _, user := range users {
		credentials, err := listSmtpCredentials(s.IdentityClient, *user.Id)
		if err != nil {
			fmt.Printf("Error listing SMTP credentials for user %s: %v\n", *user.Name, err)
			continue
		}
		allCredentials = append(allCredentials, credentials...)
	}

	return WriteJSON(w, http.StatusOK, allCredentials)
}

func listSmtpCredentials(identityClient identity.IdentityClient, userID string) ([]identity.SmtpCredentialSummary, error) {
    request := identity.ListSmtpCredentialsRequest{
        UserId: &userID,
    }

    response, err := identityClient.ListSmtpCredentials(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}