package v1

/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

// SecretRef is a reference to a Kubernetes Secret
type SecretRef struct {
	// Name is the name of the Kubernetes Secret
	Name string `json:"name"`
}

// LabelSpec contains key-value pairs that tag and categorize the Cluster/DBUser
type LabelSpec struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
