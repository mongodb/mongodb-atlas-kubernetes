package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

type AtlasStreamInstanceStatus struct {
	api.Common `json:",inline"`
	// Unique 24-hexadecimal character string that identifies the instance
	ID string `json:"id,omitempty"`
	// List that contains the hostnames assigned to the stream instance.
	Hostnames []string `json:"hostnames,omitempty"`
	// List of connections configured in the stream instance.
	Connections []StreamConnection `json:"connections,omitempty"`
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

func AtlasStreamInstanceAddConnection(name string, ref common.ResourceRefNamespaced) AtlasStreamInstanceStatusOption {
	return func(s *AtlasStreamInstanceStatus) {
		for i := range s.Connections {
			if s.Connections[i].Name == name {
				s.Connections[i].ResourceRef = ref

				return
			}
		}

		s.Connections = append(s.Connections, StreamConnection{Name: name, ResourceRef: ref})
	}
}

func AtlasStreamInstanceRemoveConnection(name string) AtlasStreamInstanceStatusOption {
	return func(s *AtlasStreamInstanceStatus) {
		for i := range s.Connections {
			if s.Connections[i].Name == name {
				s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
				break
			}
		}
	}
}
