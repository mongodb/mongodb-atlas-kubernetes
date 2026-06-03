package core

import (
	"runtime"
)

const (
	// DefaultCloudURL is default base URL for the services.
	DefaultCloudURL = "https://cloud.mongodb.com"
	// ClientName of the v2 API client.
	ClientName = "go-atlas-sdk-admin"
	// DefaultUserAgent is default user agent header.
	DefaultUserAgent = ClientName + "/" + Version + " (" + runtime.GOOS + ";" + runtime.GOARCH + ")"
)
