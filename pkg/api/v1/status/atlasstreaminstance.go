package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

type AtlasStreamInstanceStatus struct {
	Common `json:",inline"`
	// Unique 24-hexadecimal character string that identifies the instance
	ID string `json:"id,omitempty"`
	// List that contains the hostnames assigned to the stream instance.
	Hostnames []string `json:"hostnames,omitempty"`
	// List of connections configured in the stream instance.
	Connections []StreamConnection `json:"connections"`
}

type StreamConnection struct {
	// Human-readable label that uniquely identifies the stream connection
	Name string `json:"name,omitempty"`
	// Reference for the resource that contains connection configuration
	ResourceRef common.ResourceRefNamespaced `json:"resourceRef,omitempty"`
}

// +kubebuilder:object:generate=false
type AtlasStreamInstanceStatusOption func(s *AtlasStreamInstanceStatus)

func AtlasStreamInstanceDetails(ID string, hostnames []string) AtlasStreamInstanceStatusOption {
	return func(s *AtlasStreamInstanceStatus) {
		s.ID = ID
		s.Hostnames = hostnames
	}
}
