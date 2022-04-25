package atlasproject

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) ensureX509(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) (authmode.AuthModes, workflow.Result) {
	log := ctx.Log

	var specCert string
	var err error
	authModes := project.Status.AuthModes
	if key := project.X509SecretObjectKey(); key != nil {
		specCert, err = readX509CertFromSecret(r.Client, *key, log)
		if err != nil {
			return authModes, workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	if authModes.CheckAuthMode(authmode.X509) && specCert == "" {
		log.Infow("Disable x509 auth", "projectID", projectID)
		_, err := ctx.Client.X509AuthDBUsers.DisableCustomerX509(context.Background(), projectID)
		if err != nil {
			return authModes, workflow.Terminate(workflow.Internal, err.Error())
		}
		authModes.RemoveAuthMode(authmode.X509)
		return authModes, workflow.OK()
	}

	customer, _, err := ctx.Client.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), projectID)
	if err != nil {
		return authModes, workflow.Terminate(workflow.Internal, err.Error())
	}

	if specCert != customer.Cas {
		conf := mongodbatlas.CustomerX509{
			Cas: specCert,
		}
		log.Infow("Saving new x509 cert", "projectID", projectID)
		log.Debugw("New customer", "conf", conf)

		_, _, err := ctx.Client.X509AuthDBUsers.SaveConfiguration(context.Background(), projectID, &conf)
		if err != nil {
			return authModes, workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	if !authModes.CheckAuthMode(authmode.X509) && specCert != "" {
		log.Debugw("Adding new AuthMode to the status", "mode", authmode.X509)
		authModes.AddAuthMode(authmode.X509)
	}

	return authModes, workflow.OK()
}

func readX509CertFromSecret(kubeClient client.Client, secretRef client.ObjectKey, log *zap.SugaredLogger) (string, error) {
	secret := &corev1.Secret{}
	log.Debugw("reading X.509 certificate from the secret", "secretRef", secretRef)
	if err := kubeClient.Get(context.Background(), secretRef, secret); err != nil {
		return "", err
	}

	const defaultName = "ca.crt"
	certData, found := secret.Data[defaultName]
	if !found {
		if len(secret.Data) != 1 {
			errorMsg := fmt.Sprintf("the secret should have data entry with key \"%s\" or have a single data entry", defaultName)
			return "", errors.New(errorMsg)
		}

		singleKey, _ := getFirstMapItemKey(secret.Data)
		certData = secret.Data[singleKey]
	}

	cert := string(certData)
	if isNotPemEncoded(cert) {
		cert = base64.StdEncoding.EncodeToString(certData)
		if isNotPemEncoded(cert) {
			return "", errors.New("certificate has to be .pem encoded")
		}
	}

	return cert, nil
}

func isNotPemEncoded(cert string) bool {
	return !(strings.Contains(cert, "-----BEGIN") && strings.Contains(cert, "-----END"))
}

func getFirstMapItemKey(aMap map[string][]byte) (string, bool) {
	for key := range aMap {
		return key, true
	}

	return "", false
}
