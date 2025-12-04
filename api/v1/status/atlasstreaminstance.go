// Copyright 2025 MongoDB Inc
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

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// AtlasStreamInstanceStatus defines the observed state of AtlasStreamInstance.
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
