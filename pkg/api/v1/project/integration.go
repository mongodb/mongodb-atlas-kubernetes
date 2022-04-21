package project

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
	v1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	LicenseKeyRef SecretReference `json:"licenseKeyRef,omitempty"`
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// +optional
	WriteTokenRef SecretReference `json:"writeTokenRef,omitempty"`
	// +optional
	ReadTokenRef SecretReference `json:"readTokenRef,omitempty"`
	// +optional
	APIKeyRef SecretReference `json:"apiKeyRef,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	ServiceKeyRef SecretReference `json:"serviceKeyRef,omitempty"`
	// +optional
	APITokenRef SecretReference `json:"apiTokenRef,omitempty"`
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// +optional
	RoutingKeyRef SecretReference `json:"routingKeyRef,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// +optional
	URL string `json:"url,omitempty"`
	// +optional
	SecretRef SecretReference `json:"secret,omitempty"`
}

type SecretReference struct {
	// Name is the name of the Kubernetes Resource
	Name string `json:"name"`

	// Namespace is the namespace of the Kubernetes Resource
	// +optional
	Namespace string `json:"namespace"`
}

func (sr *SecretReference) GetObject(parentNamespace string) *client.ObjectKey {
	if sr == nil {
		return nil
	}

	ns := sr.Namespace
	if sr.Namespace == "" {
		ns = parentNamespace
	}
	key := kube.ObjectKey(ns, sr.Name)
	return &key
}

func (sr *SecretReference) ReadPassword(kubeClient client.Client, parentNamespace string) (string, error) {
	if sr != nil {
		nsType := sr.GetObject(parentNamespace)
		if nsType == nil {
			return "", fmt.Errorf("object is empty")
		}

		secret := &v1.Secret{}
		if err := kubeClient.Get(context.Background(), *nsType, secret); err != nil {
			return "", fmt.Errorf("can not read secret (%w), value %v", err, nsType)
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

func (i Integration) ToAtlas(defaultNS string, c client.Client) *mongodbatlas.ThirdPartyIntegration {
	result := mongodbatlas.ThirdPartyIntegration{}
	result.Type = i.Type
	result.LicenseKey, _ = i.LicenseKeyRef.ReadPassword(c, defaultNS)
	result.AccountID = i.AccountID
	result.WriteToken, _ = i.WriteTokenRef.ReadPassword(c, defaultNS)
	result.ReadToken, _ = i.ReadTokenRef.ReadPassword(c, defaultNS)
	result.APIKey, _ = i.APIKeyRef.ReadPassword(c, defaultNS)
	result.Region = i.Region
	result.ServiceKey, _ = i.ServiceKeyRef.ReadPassword(c, defaultNS)
	result.APIToken, _ = i.APITokenRef.ReadPassword(c, defaultNS)
	result.TeamName = i.TeamName
	result.ChannelName = i.ChannelName
	result.RoutingKey, _ = i.RoutingKeyRef.ReadPassword(c, defaultNS)
	result.FlowName = i.FlowName
	result.OrgName = i.OrgName
	result.URL = i.URL
	result.Secret, _ = i.SecretRef.ReadPassword(c, defaultNS)
	return &result
}

func (i Integration) Identifier() interface{} {
	return i.Type
}
