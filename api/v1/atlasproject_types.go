// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"go.uber.org/zap/zapcore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
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
	// The name length must not exceed 64 characters. The name must contain only letters, numbers, spaces, dashes, and underscores.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Name cannot be modified after project creation"
	Name string `json:"name"`

	// RegionUsageRestrictions designate the project's AWS region when using Atlas for Government.
	// This parameter should not be used with commercial Atlas.
	// In Atlas for Government, not setting this field (defaulting to NONE) means the project is restricted to COMMERCIAL_FEDRAMP_REGIONS_ONLY.
	// +kubebuilder:validation:Enum=NONE;GOV_REGIONS_ONLY;COMMERCIAL_FEDRAMP_REGIONS_ONLY
	// +kubebuilder:default:=NONE
	// +optional
	RegionUsageRestrictions string `json:"regionUsageRestrictions,omitempty"`

	// ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
	// Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.
	// +optional
	ConnectionSecret *common.ResourceRefNamespaced `json:"connectionSecretRef,omitempty"`

	// ProjectIPAccessList allows the use of the IP Access List for a Project. See more information at
	// https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
	// Deprecated: Migrate to the AtlasIPAccessList Custom Resource in accordance with the migration guide
	// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
	// +optional
	ProjectIPAccessList []project.IPAccessList `json:"projectIpAccessList,omitempty"`

	// MaintenanceWindow allows to specify a preferred time in the week to run maintenance operations. See more
	// information at https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/
	// +optional
	MaintenanceWindow project.MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// PrivateEndpoints is a list of Private Endpoints configured for the current Project.
	// Deprecated: Migrate to the AtlasPrivateEndpoint Custom Resource in accordance with the migration guide
	// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
	// +optional
	PrivateEndpoints []PrivateEndpoint `json:"privateEndpoints,omitempty"`

	// CloudProviderAccessRoles is a list of Cloud Provider Access Roles configured for the current Project.
	// Deprecated: This configuration was deprecated in favor of CloudProviderIntegrations
	// +optional
	CloudProviderAccessRoles []CloudProviderAccessRole `json:"cloudProviderAccessRoles,omitempty"`

	// CloudProviderIntegrations is a list of Cloud Provider Integration configured for the current Project.
	// +optional
	CloudProviderIntegrations []CloudProviderIntegration `json:"cloudProviderIntegrations,omitempty"`

	// AlertConfiguration is a list of Alert Configurations configured for the current Project.
	// If you use this setting, you must also set spec.alertConfigurationSyncEnabled to true for Atlas Kubernetes
	// Operator to modify project alert configurations.
	// If you omit or leave this setting empty, Atlas Kubernetes Operator doesn't alter the project's alert
	// configurations. If creating a project, Atlas applies the default project alert configurations.
	AlertConfigurations []AlertConfiguration `json:"alertConfigurations,omitempty"`

	// AlertConfigurationSyncEnabled is a flag that enables/disables Alert Configurations sync for the current Project.
	// If true - project alert configurations will be synced according to AlertConfigurations.
	// If not - alert configurations will not be modified by the operator. They can be managed through the API, CLI, and UI.
	//kubebuilder:default:=false
	// +optional
	AlertConfigurationSyncEnabled bool `json:"alertConfigurationSyncEnabled,omitempty"`

	// NetworkPeers is a list of Network Peers configured for the current Project.
	// Deprecated: Migrate to the AtlasNetworkPeering and AtlasNetworkContainer custom resources in accordance with
	// the migration guide at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
	// +optional
	NetworkPeers []NetworkPeer `json:"networkPeers,omitempty"`

	// Flag that indicates whether Atlas Kubernetes Operator creates a project with the default alert configurations.
	// If you use this setting, you must also set spec.alertConfigurationSyncEnabled to true for Atlas Kubernetes
	// Operator to modify project alert configurations.
	// If you set this parameter to false when you create a project, Atlas doesn't add the default alert configurations
	// to your project.
	// This setting has no effect on existing projects.
	// +kubebuilder:default:=true
	// +optional
	WithDefaultAlertsSettings bool `json:"withDefaultAlertsSettings,omitempty"`

	// X509CertRef is a reference to the Kubernetes Secret which contains PEM-encoded CA certificate.
	// Atlas Kubernetes Operator watches secrets only with the label atlas.mongodb.com/type=credentials to avoid
	// watching unnecessary secrets.
	// +optional
	X509CertRef *common.ResourceRefNamespaced `json:"x509CertRef,omitempty"`

	// Integrations is a list of MongoDB Atlas integrations for the project.
	// Deprecated: Migrate to the AtlasThirdPartyIntegration custom resource in accordance with the migration guide
	// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
	// +optional
	Integrations []project.Integration `json:"integrations,omitempty"`

	// EncryptionAtRest allows to set encryption for AWS, Azure and GCP providers.
	// +optional
	EncryptionAtRest *EncryptionAtRest `json:"encryptionAtRest,omitempty"`

	// Auditing represents MongoDB Maintenance Windows.
	// +optional
	Auditing *Auditing `json:"auditing,omitempty"`

	// Settings allows the configuration of the Project Settings.
	// +optional
	Settings *ProjectSettings `json:"settings,omitempty"`

	// CustomRoles lets you create and change custom roles in your cluster.
	// Use custom roles to specify custom sets of actions that the Atlas built-in roles can't describe.
	// Deprecated: Migrate to the AtlasCustomRoles custom resource in accordance with the migration guide
	// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
	// +optional
	CustomRoles []CustomRole `json:"customRoles,omitempty"`

	// Teams enable you to grant project access roles to multiple users.
	// +optional
	Teams []Team `json:"teams,omitempty"`

	// BackupCompliancePolicyRef is a reference to the backup compliance custom resource.
	// +optional
	BackupCompliancePolicyRef *common.ResourceRefNamespaced `json:"backupCompliancePolicyRef,omitempty"`
}

const hiddenField = "*** redacted ***"

//nolint:errcheck
func (p AtlasProjectSpec) MarshalLogObject(e zapcore.ObjectEncoder) error {
	printable := p.DeepCopy()
	// cleanup AlertConfigurations
	for i := range printable.AlertConfigurations {
		for j := range printable.AlertConfigurations[i].Notifications {
			printable.AlertConfigurations[i].Notifications[j].SetAPIToken(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetDatadogAPIKey(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetFlowdockAPIToken(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetDatadogAPIKey(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].MobileNumber = hiddenField
			printable.AlertConfigurations[i].Notifications[j].SetOpsGenieAPIKey(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetServiceKey(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetVictorOpsAPIKey(hiddenField)
			printable.AlertConfigurations[i].Notifications[j].SetVictorOpsRoutingKey(hiddenField)
		}
	}

	e.AddReflected("AtlasProjectSpec", printable)
	return nil
}

var _ api.AtlasCustomResource = &AtlasProject{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Atlas Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Atlas ID",type=string,JSONPath=`.status.id`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=ap

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

func (p *AtlasProject) GetStatus() api.Status {
	return p.Status
}

func (p *AtlasProject) UpdateStatus(conditions []api.Condition, options ...api.Option) {
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

func (p *AtlasProject) WithAnnotations(annotations map[string]string) *AtlasProject {
	p.Annotations = annotations
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

func (p *AtlasProject) WithBackupCompliancePolicy(name string) *AtlasProject {
	if name != "" {
		p.Spec.BackupCompliancePolicyRef = &common.ResourceRefNamespaced{Name: name, Namespace: p.Namespace}
	}
	return p
}

func (p *AtlasProject) WithBackupCompliancePolicyNamespaced(name, namespace string) *AtlasProject {
	if name != "" {
		p.Spec.BackupCompliancePolicyRef = &common.ResourceRefNamespaced{Name: name, Namespace: namespace}
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
