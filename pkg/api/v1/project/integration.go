package project

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS;PROMETHEUS
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	LicenseKeyRef common.ResourceRefNamespaced `json:"licenseKeyRef,omitempty"`
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// +optional
	WriteTokenRef common.ResourceRefNamespaced `json:"writeTokenRef,omitempty"`
	// +optional
	ReadTokenRef common.ResourceRefNamespaced `json:"readTokenRef,omitempty"`
	// +optional
	APIKeyRef common.ResourceRefNamespaced `json:"apiKeyRef,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// +optional
	RoutingKeyRef common.ResourceRefNamespaced `json:"routingKeyRef,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// +optional
	URL string `json:"url,omitempty"`
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// +optional
	UserName string `json:"username,omitempty"`
	// +optional
	PasswordRef common.ResourceRefNamespaced `json:"passwordRef,omitempty"`
	// +optional
	ServiceDiscovery string `json:"serviceDiscovery,omitempty"`
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

// ToAtlas converts the Third Party Integration to native Atlas client format. Reads the password from the Secret
func (i Integration) ToAtlas(c client.Client, defaultNS string) (result *mongodbatlas.ThirdPartyIntegration, err error) {
	result = &mongodbatlas.ThirdPartyIntegration{
		Type:                     i.Type,
		AccountID:                i.AccountID,
		Region:                   i.Region,
		TeamName:                 i.TeamName,
		ChannelName:              i.ChannelName,
		FlowName:                 i.FlowName,
		OrgName:                  i.OrgName,
		URL:                      i.URL,
		Name:                     i.Name,
		MicrosoftTeamsWebhookURL: i.MicrosoftTeamsWebhookURL,
		UserName:                 i.UserName,
		ServiceDiscovery:         i.ServiceDiscovery,
		Scheme:                   i.Scheme,
		Enabled:                  i.Enabled,
	}

	readPassword := func(passwordField common.ResourceRefNamespaced, target *string, errors *[]error) {
		if passwordField.Name == "" {
			return
		}

		*target, err = passwordField.ReadPassword(c, defaultNS)
		storeError(err, errors)
	}

	errorList := make([]error, 0)
	readPassword(i.LicenseKeyRef, &result.LicenseKey, &errorList)
	readPassword(i.WriteTokenRef, &result.WriteToken, &errorList)
	readPassword(i.ReadTokenRef, &result.ReadToken, &errorList)
	readPassword(i.APIKeyRef, &result.APIKey, &errorList)
	readPassword(i.ServiceKeyRef, &result.ServiceKey, &errorList)
	readPassword(i.APITokenRef, &result.APIToken, &errorList)
	readPassword(i.RoutingKeyRef, &result.RoutingKey, &errorList)
	readPassword(i.SecretRef, &result.Secret, &errorList)
	readPassword(i.PasswordRef, &result.Password, &errorList)

	if len(errorList) != 0 {
		firstError := (errorList)[0]
		return nil, firstError
	}

	return result, nil
}

func IntegrationFromAtlas(atlasIntegration *mongodbatlas.ThirdPartyIntegration, kubeClient client.Client, nameSpace string, projectID string) (*Integration, error) {
	result := &Integration{
		Type:                     atlasIntegration.Type,
		LicenseKeyRef:            common.ResourceRefNamespaced{},
		AccountID:                atlasIntegration.AccountID,
		WriteTokenRef:            common.ResourceRefNamespaced{},
		ReadTokenRef:             common.ResourceRefNamespaced{},
		APIKeyRef:                common.ResourceRefNamespaced{},
		Region:                   atlasIntegration.Region,
		ServiceKeyRef:            common.ResourceRefNamespaced{},
		APITokenRef:              common.ResourceRefNamespaced{},
		TeamName:                 atlasIntegration.TeamName,
		ChannelName:              atlasIntegration.ChannelName,
		RoutingKeyRef:            common.ResourceRefNamespaced{},
		FlowName:                 atlasIntegration.FlowName,
		OrgName:                  atlasIntegration.OrgName,
		URL:                      atlasIntegration.URL,
		SecretRef:                common.ResourceRefNamespaced{},
		Name:                     atlasIntegration.Name,
		MicrosoftTeamsWebhookURL: atlasIntegration.MicrosoftTeamsWebhookURL,
		UserName:                 atlasIntegration.UserName,
		PasswordRef:              common.ResourceRefNamespaced{},
		ServiceDiscovery:         atlasIntegration.ServiceDiscovery,
		Scheme:                   atlasIntegration.Scheme,
		Enabled:                  atlasIntegration.Enabled,
	}

	writePassword := func(password string, resourceRef *common.ResourceRefNamespaced, errors *[]error, integrationType string, secretName string) {
		if password == "" {
			return
		}
		// To ensure the resource name is unique in the cluster, we combine the projectID, the integration type and the name of the field
		// There is only one integration of each type in a project
		//https://www.mongodb.com/docs/atlas/reference/api/third-party-integration-settings-get-one/
		passwordSecretName := projectID + integrationType + secretName

		resourceRef.Name = passwordSecretName
		resourceRef.Namespace = nameSpace

		data := map[string][]byte{
			"password": []byte(password),
		}
		object := metav1.ObjectMeta{Name: passwordSecretName}
		secret := &corev1.Secret{Data: data, ObjectMeta: object}

		if err := kubeClient.Create(context.Background(), secret); err != nil {
			storeError(err, errors)
		}
	}

	// Store any secret in kubernetes and store reference in operator Integration struct
	errorList := make([]error, 0)
	writePassword(atlasIntegration.LicenseKey, &result.LicenseKeyRef, &errorList, atlasIntegration.Type, "LicenseKey")
	writePassword(atlasIntegration.WriteToken, &result.WriteTokenRef, &errorList, atlasIntegration.Type, "WriteToken")
	writePassword(atlasIntegration.ReadToken, &result.ReadTokenRef, &errorList, atlasIntegration.Type, "ReadToken")
	writePassword(atlasIntegration.APIKey, &result.APIKeyRef, &errorList, atlasIntegration.Type, "APIKey")
	writePassword(atlasIntegration.ServiceKey, &result.ServiceKeyRef, &errorList, atlasIntegration.Type, "ServiceKey")
	writePassword(atlasIntegration.APIToken, &result.APITokenRef, &errorList, atlasIntegration.Type, "APIToken")
	writePassword(atlasIntegration.RoutingKey, &result.RoutingKeyRef, &errorList, atlasIntegration.Type, "RoutingKey")
	writePassword(atlasIntegration.Secret, &result.SecretRef, &errorList, atlasIntegration.Type, "Secret")
	writePassword(atlasIntegration.Password, &result.PasswordRef, &errorList, atlasIntegration.Type, "Password")

	if len(errorList) != 0 {
		firstError := (errorList)[0]
		return nil, firstError
	}

	return result, nil
}

func (i Integration) Identifier() interface{} {
	return i.Type
}

func storeError(err error, errors *[]error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}
