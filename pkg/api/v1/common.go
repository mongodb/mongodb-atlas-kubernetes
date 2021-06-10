package v1

/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

// ResourceRef is a reference to a Kubernetes Resource
type ResourceRef struct {
	// Name is the name of the Kubernetes Resource
	Name string `json:"name"`
}

// ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace
type ResourceRefNamespaced struct {
	// Name is the name of the Kubernetes Resource
	Name string `json:"name"`

	// Namespace is the namespace of the Kubernetes Resource
	Namespace string `json:"namespace"`
}

// LabelSpec contains key-value pairs that tag and categorize the Cluster/DBUser
type LabelSpec struct {
	// +kubebuilder:validation:MaxLength:=255
	Key   string `json:"key"`
	Value string `json:"value"`
}
