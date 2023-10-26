package main

import (
	"math/rand"
	"time"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/core"
)

type Configuration struct {
	AppPort          int    `json:"app_port"`
	DBUsername       string `json:"DB_USERNAME"`
	DBPassword       string `json:"DB_PASSWORD"`
	DBHost           string `json:"DB_HOST"`
	DBPort           string `json:"DB_PORT"`
	DBDatabase       string `json:"DB_DATABASE"`
	TenancyID        string `json:"tenancy_id"`
	UserID           string `json:"user_id"`
	FingerprintPath  string `json:"fingerprint_path"`
	PrivateKeyPath   string `json:"private_key_path"`
	HomeRegion       string `json:"home_region"`
	JWTSecret        string `json:"jwt_secret"`
}

type APIServer struct {
	ListenAddr     string
	Store          Storage
	Config         common.ConfigurationProvider
	IdentityClient identity.IdentityClient
	ComputeClients []core.ComputeClient
	TenancyID      string
	Configuration  Configuration
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}