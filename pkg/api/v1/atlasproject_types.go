/*
Copyright 2020 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// Important:
// The procedure working with this file:
// 1. Edit the file
// 1. Run "make generate" to regenerate code
// 2. Run "make manifests" to regenerate the CRD

// Dev note: this file should be placed in "v1" package (not the nested one) as 'make manifests' doesn't generate the proper
// CRD - this may be addressed later as having a subpackage may get a much nicer code

func init() {
	SchemeBuilder.Register(&AtlasProject{}, &AtlasProjectList{})
}

// AtlasProjectSpec defines the desired state of Project in Atlas
type AtlasProjectSpec struct {

	// Name is the name of the Project that is created in Atlas by the Operator if it doesn't exist yet.
	Name string `json:"name"`

	// ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
	// Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.
	// +optional
	ConnectionSecret *common.ResourceRefNamespaced `json:"connectionSecretRef,omitempty"`

	// ProjectIPAccessList allows to enable the IP Access List for the Project. See more information at
	// https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
	// +optional
	ProjectIPAccessList []project.IPAccessList `json:"projectIpAccessList,omitempty"`

	// MaintenanceWindow allows to specify a preferred time in the week to run maintenance operations. See more
	// information at https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/
	// +optional
	MaintenanceWindow project.MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// PrivateEndpoints is a list of Private Endpoints configured for the current Project.
	PrivateEndpoints []PrivateEndpoint `json:"privateEndpoints,omitempty"`
	// CloudProviderAccessRoles is a list of Cloud Provider Access Roles configured for the current Project.
	CloudProviderAccessRoles []CloudProviderAccessRole `json:"cloudProviderAccessRoles,omitempty"`

	// AlertConfiguration is a list of Alert Configurations configured for the current Project.
	AlertConfigurations []AlertConfiguration `json:"alertConfigurations,omitempty"`

	// AlertConfigurationSyncEnabled is a flag that enables/disables Alert Configurations sync for the current Project.
	// If true - project alert configurations will be synced according to AlertConfigurations.
	// If not - alert configurations will not be modified by the operator. They can be managed through API, cli, UI.
	//kubebuilder:default:=false
	// +optional
	AlertConfigurationSyncEnabled bool `json:"alertConfigurationSyncEnabled,omitempty"`

	// NetworkPeers is a list of Network Peers configured for the current Project.
	NetworkPeers []NetworkPeer `json:"networkPeers,omitempty"`

	// Flag that indicates whether to create the new project with the default alert settings enabled. This parameter defaults to true
	// +kubebuilder:default:=true
	// +optional
	WithDefaultAlertsSettings bool `json:"withDefaultAlertsSettings,omitempty"`

	// X509CertRef is the name of the Kubernetes Secret which contains PEM-encoded CA certificate
	X509CertRef *common.ResourceRefNamespaced `json:"x509CertRef,omitempty"`

	// Integrations is a list of MongoDB Atlas integrations for the project
	// +optional
	Integrations []project.Integration `json:"integrations,omitempty"`

	// EncryptionAtRest allows to set encryption for AWS, Azure and GCP providers
	// +optional
	EncryptionAtRest *EncryptionAtRest `json:"encryptionAtRest,omitempty"`

	// Auditing represents MongoDB Maintenance Windows
	// +optional
	Auditing *Auditing `json:"auditing,omitempty"`

	// Settings allow to set Project Settings for the project
	// +optional
	Settings *ProjectSettings `json:"settings,omitempty"`

	// The customRoles lets you create, and change custom roles in your cluster. Use custom roles to specify custom sets of actions that the Atlas built-in roles can't describe.
	// +optional
	CustomRoles []CustomRole `json:"customRoles,omitempty"`

	// Teams enable you to grant project access roles to multiple users.
	// +optional
	Teams []Team `json:"teams,omitempty"`
}

const hiddenField = "*** redacted ***"

//nolint:errcheck
func (p AtlasProjectSpec) MarshalLogObject(e zapcore.ObjectEncoder) error {
	printable := p.DeepCopy()
	// cleanup encryption at EncryptionAtRest
	if printable.EncryptionAtRest != nil {
		printable.EncryptionAtRest.AwsKms.AccessKeyID = hiddenField
		printable.EncryptionAtRest.AwsKms.CustomerMasterKeyID = hiddenField
		printable.EncryptionAtRest.AwsKms.SecretAccessKey = hiddenField
		printable.EncryptionAtRest.AwsKms.RoleID = hiddenField
		printable.EncryptionAtRest.AzureKeyVault.Secret = hiddenField
		printable.EncryptionAtRest.GoogleCloudKms.ServiceAccountKey = hiddenField
	}

	// cleanup AlertConfigurations
	for i := range printable.AlertConfigurations {
		for j := range printable.AlertConfigurations[i].Notifications {
			printable.AlertConfigurations[i].Notifications[j].APIToken = hiddenField
			printable.AlertConfigurations[i].Notifications[j].DatadogAPIKey = hiddenField
			printable.AlertConfigurations[i].Notifications[j].FlowdockAPIToken = hiddenField
			printable.AlertConfigurations[i].Notifications[j].DatadogAPIKey = hiddenField
			printable.AlertConfigurations[i].Notifications[j].MobileNumber = hiddenField
			printable.AlertConfigurations[i].Notifications[j].OpsGenieAPIKey = hiddenField
			printable.AlertConfigurations[i].Notifications[j].ServiceKey = hiddenField
			printable.AlertConfigurations[i].Notifications[j].VictorOpsAPIKey = hiddenField
			printable.AlertConfigurations[i].Notifications[j].VictorOpsRoutingKey = hiddenField
		}
	}

	e.AddReflected("AtlasProjectSpec", printable)
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com

// AtlasProject is the Schema for the atlasprojects API
type AtlasProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasProjectSpec          `json:"spec,omitempty"`
	Status status.AtlasProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasProjectList contains a list of AtlasProject
type AtlasProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasProject `json:"items"`
}

// ID is just a shortcut for ID from the status
func (p AtlasProject) ID() string {
	return p.Status.ID
}

func (p *AtlasProject) ConnectionSecretObjectKey() *client.ObjectKey {
	if p.Spec.ConnectionSecret != nil {
		var key client.ObjectKey
		if p.Spec.ConnectionSecret.Namespace != "" {
			key = kube.ObjectKey(p.Spec.ConnectionSecret.Namespace, p.Spec.ConnectionSecret.Name)
		} else {
			key = kube.ObjectKey(p.Namespace, p.Spec.ConnectionSecret.Name)
		}
		return &key
	}
	return nil
}

func (p *AtlasProject) GetStatus() status.Status {
	return p.Status
}

func (p *AtlasProject) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	p.Status.Conditions = conditions
	p.Status.ObservedGeneration = p.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasProjectStatusOption)
		v(&p.Status)
	}
}

func (p *AtlasProject) X509SecretObjectKey() *client.ObjectKey {
	return p.Spec.X509CertRef.GetObject(p.Namespace)
}

func (p *AtlasProject) GetCondition(condType status.ConditionType) *status.Condition {
	for _, cond := range p.Status.Conditions {
		if cond.Type == condType {
			return &cond
		}
	}
	return nil
}

// CheckConditions check AtlasProject conditions
// First check if ReadyType condition exists and is True
// Then check if a condition of ProjectReadyType, DeploymentReadyType, IPAccessListReadyType or PrivateEndpointReadyType condition exists and is False
// returns nil otherwise
func (p *AtlasProject) CheckConditions() *status.Condition {
	cond := p.GetCondition(status.ReadyType)
	if cond != nil && cond.Status == corev1.ConditionTrue {
		return cond
	}
	cond = p.GetCondition(status.ProjectReadyType)
	if cond != nil && cond.Status == corev1.ConditionFalse {
		return cond
	}
	cond = p.GetCondition(status.DeploymentReadyType)
	if cond != nil && cond.Status == corev1.ConditionFalse {
		return cond
	}
	cond = p.GetCondition(status.IPAccessListReadyType)
	if cond != nil && cond.Status == corev1.ConditionFalse {
		return cond
	}
	cond = p.GetCondition(status.PrivateEndpointReadyType)
	if cond != nil && cond.Status == corev1.ConditionFalse {
		return cond
	}
	return nil
}

// ************************************ Builder methods *************************************************

func NewProject(namespace, name, nameInAtlas string) *AtlasProject {
	return &AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasProjectSpec{
			Name: nameInAtlas,
		},
	}
}

func (p *AtlasProject) WithLabels(labels map[string]string) *AtlasProject {
	p.Labels = labels
	return p
}

func (p *AtlasProject) WithAnnotations(labels map[string]string) *AtlasProject {
	p.Labels = labels
	return p
}

func (p *AtlasProject) WithName(name string) *AtlasProject {
	p.Name = name
	return p
}

func (p *AtlasProject) WithAtlasName(name string) *AtlasProject {
	p.Spec.Name = name
	return p
}

func (p *AtlasProject) WithConnectionSecret(name string) *AtlasProject {
	if name != "" {
		p.Spec.ConnectionSecret = &common.ResourceRefNamespaced{Name: name, Namespace: p.Namespace}
	}
	return p
}

func (p *AtlasProject) WithConnectionSecretNamespaced(name, namespace string) *AtlasProject {
	if name != "" {
		p.Spec.ConnectionSecret = &common.ResourceRefNamespaced{Name: name, Namespace: namespace}
	}
	return p
}

func (p *AtlasProject) WithIPAccessList(ipAccess project.IPAccessList) *AtlasProject {
	if p.Spec.ProjectIPAccessList == nil {
		p.Spec.ProjectIPAccessList = []project.IPAccessList{}
	}
	p.Spec.ProjectIPAccessList = append(p.Spec.ProjectIPAccessList, ipAccess)
	return p
}

func (p *AtlasProject) WithMaintenanceWindow(window project.MaintenanceWindow) *AtlasProject {
	p.Spec.MaintenanceWindow = window
	return p
}

func DefaultProject(namespace, connectionSecretName string) *AtlasProject {
	return NewProject(namespace, "test-project", namespace).WithConnectionSecret(connectionSecretName)
}
