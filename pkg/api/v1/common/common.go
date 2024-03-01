package common

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

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
	// +optional
	Namespace string `json:"namespace"`
}

func (in ResourceRefNamespaced) Key() string {
	return in.Name + "|" + in.Namespace
}

// LabelSpec contains key-value pairs that tag and categorize the Cluster/DBUser
type LabelSpec struct {
	// +kubebuilder:validation:MaxLength:=255
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (rn *ResourceRefNamespaced) GetObject(parentNamespace string) *client.ObjectKey {
	if rn == nil {
		return nil
	}

	ns := SelectNamespace(rn.Namespace, parentNamespace)
	key := kube.ObjectKey(ns, rn.Name)
	return &key
}

func (rn *ResourceRefNamespaced) ReadPassword(ctx context.Context, kubeClient client.Client, parentNamespace string) (string, error) {
	if rn != nil {
		secret := &v1.Secret{}
		if err := kubeClient.Get(ctx, *rn.GetObject(parentNamespace), secret); err != nil {
			return "", err
		}
		p, exist := secret.Data["password"]
		switch {
		case !exist:
			return "", fmt.Errorf("secret %s is invalid: it doesn't contain 'password' field", secret.Name)
		case len(p) == 0:
			return "", fmt.Errorf("secret %s is invalid: the 'password' field is empty", secret.Name)
		default:
			return string(p), nil
		}
	}
	return "", nil
}

// SelectNamespace returns first non-empty namespace from the list
// "", "", "first", "second" => "first"
func SelectNamespace(namespaces ...string) string {
	for _, namespace := range namespaces {
		if namespace != "" {
			return namespace
		}
	}

	return ""
}
