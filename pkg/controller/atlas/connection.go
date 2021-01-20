package atlas

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	orgIDKey      = "orgId"
	publicAPIKey  = "publicApiKey"
	privateAPIKey = "privateApiKey"
)

// Connection encapsulates Atlas connectivity information that is necessary to perform API requests
type Connection struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

// ReadConnection reads Atlas API connection parameters from AtlasProject Secret or from the default Operator one if the
// former is not specified
func ReadConnection(log *zap.SugaredLogger, kubeClient client.Client, operatorName string, projectOverrideSecretRef *client.ObjectKey) (Connection, workflow.Result) {
	if projectOverrideSecretRef != nil {
		// TODO is it possible that part of connection (like orgID is still in the Operator level secret and needs to get merged?)
		log.Infof("Reading Atlas API credentials from the AtlasProject Secret %s", projectOverrideSecretRef)
		return readAtlasConnectionFromSecret(kubeClient, *projectOverrideSecretRef)
	}

	// TODO check the default "Operator level" Secret
	// return readAtlasConnectionFromSecret(operatorName + "-connection")
	return Connection{}, workflow.Terminate(workflow.AtlasCredentialsNotProvided, "the API keys are not configured")
}

func readAtlasConnectionFromSecret(kubeClient client.Client, secretRef client.ObjectKey) (Connection, workflow.Result) {
	secret := &corev1.Secret{}
	if err := kubeClient.Get(context.Background(), secretRef, secret); err != nil {
		return Connection{}, workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
	}
	secretData := make(map[string]string)
	for k, v := range secret.Data {
		secretData[k] = string(v)
	}

	if err := validateConnectionSecret(secretRef, secretData); err != nil {
		return Connection{}, workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
	}

	return Connection{
		OrgID:      secretData["orgId"],
		PublicKey:  secretData["publicApiKey"],
		PrivateKey: secretData["privateApiKey"],
	}, workflow.OK()
}

func validateConnectionSecret(secretRef client.ObjectKey, secretData map[string]string) error {
	var missingFields []string
	requiredKeys := []string{orgIDKey, publicAPIKey, privateAPIKey}

	for _, key := range requiredKeys {
		if _, ok := secretData[key]; !ok {
			missingFields = append(missingFields, key)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("the following fields are missing in the Secret %v: %v", secretRef, missingFields)
	}
	return nil
}
