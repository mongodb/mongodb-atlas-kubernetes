package api

import (
	corev1 "k8s.io/api/core/v1"
)

type LocalRef string

// +k8s:deepcopy-gen=false

// CredentialsProvider gives access to custom local credentials
type CredentialsProvider interface {
	Credentials() *corev1.LocalObjectReference
}

// +k8s:deepcopy-gen=false

// ResourceWithCredentials is to be implemented by all CRDs using custom local credentials
type ResourceWithCredentials interface {
	CredentialsProvider
	GetName() string
	GetNamespace() string
}

// LocalCredentialHolder is to be embedded by Specs of CRDs using custom local credentials
type LocalCredentialHolder struct {
	ConnectionSecret *corev1.LocalObjectReference `json:"connectionSecret,omitempty"`
}

func (ch *LocalCredentialHolder) Credentials() *corev1.LocalObjectReference {
	return ch.ConnectionSecret
}
