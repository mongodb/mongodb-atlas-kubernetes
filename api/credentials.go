package api

import "sigs.k8s.io/controller-runtime/pkg/client"

type LocalRef string

// +k8s:deepcopy-gen=false

// CredentialsProvider gives access to custom local credentials
type CredentialsProvider interface {
	Credentials() *LocalObjectReference
}

// +k8s:deepcopy-gen=false

// ObjectWithCredentials is a Kubernetes Object interface with credentials
type ObjectWithCredentials interface {
	client.Object
	CredentialsProvider
}
