package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DeleteAll  bool
	Lifetime   int // lifetime in hours
	PublicKey  string
	PrivateKey string
	ManagerUrl string
}

const DefaultManagerURL = "https://cloud-qa.mongodb.com/"

func NewConfig() (*Config, error) {
	deleteAllStr := os.Getenv("CLEAN_ALL")
	deleteAll := deleteAllStr != "false"

	var lifetime int
	if !deleteAll {
		lifetimeStr := os.Getenv("MAX_PROJECT_LIFETIME")
		if lifetimeStr == "" {
			return nil, fmt.Errorf("MAX_PROJECT_LIFETIME is not set")
		}
		var err error
		lifetime, err = strconv.Atoi(lifetimeStr)
		if err != nil {
			return nil, fmt.Errorf("MAX_PROJECT_LIFETIME must be an integer")
		}
	}

	publicKey := os.Getenv("MCLI_PUBLIC_API_KEY")
	if publicKey == "" {
		return nil, fmt.Errorf("MCLI_PUBLIC_API_KEY must be set")
	}
	privateKey := os.Getenv("MCLI_PRIVATE_API_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("MCLI_PRIVATE_API_KEY must be set")
	}
	managerUrl := os.Getenv("MCLI_OPS_MANAGER_URL")
	if managerUrl == "" {
		managerUrl = DefaultManagerURL
	}

	return &Config{
		DeleteAll:  deleteAll,
		Lifetime:   lifetime,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		ManagerUrl: managerUrl,
	}, nil
}
