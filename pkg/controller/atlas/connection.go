package atlas

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Connection encapsulates Atlas connectivity information that is necessary to perform API requests
type Connection struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

// ReadConnection reads Atlas API connection parameters from AtlasProject Secret or from the default Operator one if the
// former is not specified
func ReadConnection(kubeClient client.Client, operatorName string, projectOverrideSecretRef *client.ObjectKey, log *zap.SugaredLogger) (Connection, error) {
	if projectOverrideSecretRef != nil {
		// TODO is it possible that part of connection (like orgID is still in the Operator level secret and needs to get merged?)
		log.Infof("Reading Atlas API credentials from the AtlasProject Secret %s", projectOverrideSecretRef)
		return readAtlasConnectionFromSecret(kubeClient, *projectOverrideSecretRef)
	}
	// TODO check the default "Operator level" Secret
	// return readAtlasConnectionFromSecret(operatorName + "-connection")
	return Connection{}, errors.New("The API keys are not configured")
}

func readAtlasConnectionFromSecret(kubeClient client.Client, secretRef client.ObjectKey) (Connection, error) {
	secret := &corev1.Secret{}
	if err := kubeClient.Get(context.Background(), secretRef, secret); err != nil {
		return Connection{}, err
	}
	secretData := make(map[string]string)
	for k, v := range secret.Data {
		secretData[k] = string(v)
	}
	if _, ok := secretData["orgId"]; !ok {
		return Connection{}, fmt.Errorf("Missing field with key 'orgId' in the secret %v", secretRef)
	}
	if _, ok := secretData["publicApiKey"]; !ok {
		return Connection{}, fmt.Errorf("Missing field with key 'publicApiKey' in the secret %v", secretRef)
	}
	if _, ok := secretData["privateApiKey"]; !ok {
		return Connection{}, fmt.Errorf("Missing field with key 'privateApiKey' in the secret %v", secretRef)
	}
	return Connection{
		OrgID:      secretData["orgId"],
		PublicKey:  secretData["publicApiKey"],
		PrivateKey: secretData["privateApiKey"],
	}, nil
}
