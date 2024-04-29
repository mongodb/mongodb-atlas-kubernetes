package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

type AtlasStreamConnectionSpec struct {
	// Human-readable label that uniquely identifies the stream connection
	Name string `json:"name"`
	// Type of the connection. Can be either Cluster or Kafka
	// +kubebuilder:validation:Enum:=Kafka;Cluster;Sample
	ConnectionType string `json:"type"`
	// The configuration to be used to connect to a Atlas Cluster
	ClusterConfig *ClusterConnectionConfig `json:"clusterConfig,omitempty"`
	// The configuration to be used to connect to a Kafka Cluster
	KafkaConfig *StreamsKafkaConnection `json:"kafkaConfig,omitempty"`
}

type ClusterConnectionConfig struct {
	// Name of the cluster configured for this connection
	Name string `json:"name"`
	// The name of a Built in or Custom DB Role to connect to an Atlas Cluster
	Role StreamsClusterDBRole `json:"role"`
}

type StreamsClusterDBRole struct {
	// The name of the role to use. Can be a built in role or a custom role
	Name string `json:"name"`
	// Type of the DB role. Can be either BuiltIn or Custom
	// +kubebuilder:validation:Enum:=BUILT_IN;CUSTOM
	RoleType string `json:"type"`
}

type StreamsKafkaConnection struct {
	// User credentials required to connect to a Kafka Cluster. Includes the authentication type, as well as the parameters for that authentication mode
	Authentication StreamsKafkaAuthentication `json:"authentication"`
	// Comma separated list of server addresses
	BootstrapServers string `json:"bootstrapServers"`
	// Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use
	Security StreamsKafkaSecurity `json:"security"`
	// A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters
	Config map[string]string `json:"config,omitempty"`
}

type StreamsKafkaAuthentication struct {
	// Style of authentication. Can be one of PLAIN, SCRAM-256, or SCRAM-512
	// +kubebuilder:validation:Enum:=PLAIN;SCRAM-256;SCRAM-512
	Mechanism string `json:"mechanism"`
	// Reference to the secret containing th Username and Password of the account to connect to the Kafka cluster.
	Credentials common.ResourceRefNamespaced `json:"credentials"`
}

type StreamsKafkaSecurity struct {
	// Describes the transport type. Can be either PLAINTEXT or SSL
	// +kubebuilder:validation:Enum:=PLAINTEXT;SSL
	Protocol string `json:"protocol"`
	// A trusted, public x509 certificate for connecting to Kafka over SSL
	Certificate common.ResourceRefNamespaced `json:"certificate,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AtlasStreamConnection is the Schema for the atlasstreamconnections API
type AtlasStreamConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasStreamConnectionSpec          `json:"spec,omitempty"`
	Status status.AtlasStreamConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasStreamConnectionList contains a list of AtlasStreamConnection
type AtlasStreamConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasStreamConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasStreamConnection{}, &AtlasStreamConnectionList{})
}

// GetStatus implements status.Reader
func (f *AtlasStreamConnection) GetStatus() status.Status {
	return f.Status
}

func (f *AtlasStreamConnection) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	f.Status.Conditions = conditions
	f.Status.ObservedGeneration = f.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasStreamConnectionStatusOption)
		v(&f.Status)
	}
}
