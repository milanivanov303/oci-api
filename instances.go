package main

import (
	"net/http"
	"fmt"
	"context"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleInstances(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetInstances(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetInstances(w http.ResponseWriter, r *http.Request) error {
	compartmentID := mux.Vars(r)["compartment_id"]

	instances, err := s.listInstancesFromAllRegions(compartmentID)
    if err != nil {
        fmt.Println("Error listing instances:", err)
        return err
    }

	return WriteJSON(w, http.StatusOK, instances)
}

func (s *APIServer) listInstancesFromAllRegions(compartmentID string) ([]core.Instance, error) {
    var allInstances []core.Instance
    for _, client := range s.ComputeClients {
        instances, err := listInstances(client, compartmentID)
        if err != nil {
            fmt.Println("Error listing instances:", err)
            return nil, err
        }
        allInstances = append(allInstances, instances...)
    }

    return allInstances, nil
}

func listInstances(client core.ComputeClient, compartmentID string) ([]core.Instance, error) {
    request := core.ListInstancesRequest{
        CompartmentId: &compartmentID,
    }

    response, err := client.ListInstances(context.Background(), request)
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}
