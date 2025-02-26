package reconciler

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

const (
	orgIDKey      = "orgId"
	publicAPIKey  = "publicApiKey"
	privateAPIKey = "privateApiKey"
)

func (r *AtlasReconciler) ResolveCredentials(ctx context.Context, referrer project.ProjectReferrerObject) (*atlas.Credentials, error) {
	connectionSecret := r.connectionSecretRef(referrer)
	if connectionSecret != nil && connectionSecret.Name != "" {
		creds, err := GetAtlasCredentials(ctx, r.Client, connectionSecret, &r.GlobalSecretRef)
		if err != nil {
			return nil, fmt.Errorf("error getting credentials from connection secret: %w", err)
		}
		return creds, nil
	}

	prj, err := r.fetchProject(ctx, referrer)
	if err != nil {
		return nil, fmt.Errorf("error resolving project reference: %w", err)
	}

	var projectSecret *client.ObjectKey
	if prj != nil {
		projectSecret = prj.ConnectionSecretObjectKey()
	}

	creds, err := GetAtlasCredentials(ctx, r.Client, projectSecret, &r.GlobalSecretRef)
	if err != nil {
		return nil, fmt.Errorf("error getting credentials from project reference: %w", err)
	}
	return creds, nil
}

func (r *AtlasReconciler) connectionSecretRef(pro project.ProjectReferrerObject) *client.ObjectKey {
	key := client.ObjectKeyFromObject(pro)
	pdr := pro.ProjectDualRef()
	if pdr.ConnectionSecret == nil {
		return nil
	}
	key.Name = pdr.ConnectionSecret.Name
	return &key
}

func GetAtlasCredentials(ctx context.Context, k8sClient client.Client, secretRef, fallbackRef *client.ObjectKey) (*atlas.Credentials, error) {
	if secretRef == nil {
		secretRef = fallbackRef
	}

	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, *secretRef, secret); err != nil {
		return nil, fmt.Errorf("failed to read Atlas API credentials from the secret %s: %w", secretRef.String(), err)
	}

	apiKeys := atlas.APIKeys{
		OrgID:      string(secret.Data[orgIDKey]),
		PublicKey:  string(secret.Data[publicAPIKey]),
		PrivateKey: string(secret.Data[privateAPIKey]),
	}

	if missingFields, valid := validate(&apiKeys); !valid {
		return nil, fmt.Errorf("the following fields are missing in the secret %v: %v", secretRef, missingFields)
	}

	return &atlas.Credentials{APIKeys: &apiKeys}, nil
}

func validate(apiKeys *atlas.APIKeys) ([]string, bool) {
	missingFields := make([]string, 0, 3)

	if apiKeys.OrgID == "" {
		missingFields = append(missingFields, orgIDKey)
	}

	if apiKeys.PublicKey == "" {
		missingFields = append(missingFields, publicAPIKey)
	}

	if apiKeys.PrivateKey == "" {
		missingFields = append(missingFields, privateAPIKey)
	}

	if len(missingFields) > 0 {
		return missingFields, false
	}

	return nil, true
}
