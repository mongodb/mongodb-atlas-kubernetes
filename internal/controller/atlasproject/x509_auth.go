package atlasproject

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

func terminateX509(workflowCtx *workflow.Context, err error) workflow.Result {
	workflowCtx.SetConditionFalseMsg(api.X509AuthReadyType, err.Error())
	return workflow.Terminate(workflow.ProjectX509NotConfigured, err)
}

func emptyX509(workflowCtx *workflow.Context) workflow.Result {
	workflowCtx.UnsetCondition(api.X509AuthReadyType)
	return idleX509()
}

func idleX509() workflow.Result {
	return workflow.OK()
}

func (r *AtlasProjectReconciler) ensureX509(ctx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.Result {
	atlasProject.Status.AuthModes.AddAuthMode(authmode.Scram)

	hasAuthModesX509 := atlasProject.Status.AuthModes.CheckAuthMode(authmode.X509)
	hasX509Cert := atlasProject.X509SecretObjectKey() != nil

	switch {
	case hasX509Cert:
		return r.enableX509Authentication(ctx, atlasProject)
	case !hasX509Cert && hasAuthModesX509:
		return r.disableX509Authentication(ctx, atlasProject)
	default:
		ctx.EnsureStatusOption(status.AtlasProjectAuthModesOption(atlasProject.Status.AuthModes))
	}
	return idleX509()
}

func (r *AtlasProjectReconciler) enableX509Authentication(ctx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.Result {
	specCert, err := readX509CertFromSecret(ctx.Context, r.Client, *atlasProject.X509SecretObjectKey(), r.Log)
	if err != nil {
		return terminateX509(ctx, err)
	}

	ldapConfig, _, err := ctx.SdkClient.LDAPConfigurationApi.GetLDAPConfiguration(ctx.Context, atlasProject.ID()).Execute()
	if err != nil {
		return terminateX509(ctx, err)
	}

	customerX509 := ldapConfig.GetCustomerX509()
	if specCert != customerX509.GetCas() {
		conf := admin.UserSecurity{
			CustomerX509: &admin.DBUserTLSX509Settings{
				Cas: &specCert,
			},
		}
		r.Log.Infow("Saving new x509 cert", "projectID", atlasProject.ID())

		_, _, err = ctx.SdkClient.LDAPConfigurationApi.SaveLDAPConfiguration(ctx.Context, atlasProject.ID(), &conf).Execute()
		if err != nil {
			return terminateX509(ctx, err)
		}
	}

	atlasProject.Status.AuthModes.AddAuthMode(authmode.X509)
	ctx.EnsureStatusOption(status.AtlasProjectAuthModesOption(atlasProject.Status.AuthModes))

	return idleX509()
}

func (r *AtlasProjectReconciler) disableX509Authentication(ctx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.Result {
	r.Log.Infow("Disable x509 auth", "projectID", atlasProject.ID())
	_, _, err := ctx.SdkClient.X509AuthenticationApi.DisableCustomerManagedX509(ctx.Context, atlasProject.ID()).Execute()
	if err != nil {
		return terminateX509(ctx, err)
	}

	atlasProject.Status.AuthModes.RemoveAuthMode(authmode.X509)
	ctx.EnsureStatusOption(status.AtlasProjectAuthModesOption(atlasProject.Status.AuthModes))

	return emptyX509(ctx)
}

func readX509CertFromSecret(ctx context.Context, kubeClient client.Client, secretRef client.ObjectKey, log *zap.SugaredLogger) (string, error) {
	secret := &corev1.Secret{}
	log.Debugw("reading X.509 certificate from the secret", "secretRef", secretRef)
	if err := kubeClient.Get(ctx, secretRef, secret); err != nil {
		return "", err
	}

	const defaultName = "ca.crt"
	certData, found := secret.Data[defaultName]
	if !found {
		if len(secret.Data) != 1 {
			errorMsg := fmt.Sprintf("the secret should have data entry with key \"%s\" or have a single data entry, data: %v", defaultName, secret.Data)
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
