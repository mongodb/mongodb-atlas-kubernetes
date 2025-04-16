// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// AtlasIPAccessListStatus is the most recent observed status of the AtlasIPAccessList cluster. Read-only.
type AtlasIPAccessListStatus struct {
	api.Common `json:",inline"`
	//
	// Status is the state of the ip access list
	Entries []IPAccessEntryStatus `json:"entries,omitempty"`
}

type IPAccessEntryStatus struct {
	// Entry is the ip access Atlas is managing
	Entry string `json:"entry"`
	// Status is the correspondent state of the entry
	Status string `json:"status"`
}

// +kubebuilder:object:generate=false

type AtlasIPAccessListStatusOption func(s *AtlasIPAccessListStatus)

func AddIPAccessListEntryStatus(entry, entryStatus string) AtlasIPAccessListStatusOption {
	return func(s *AtlasIPAccessListStatus) {
		for ix, ipEntryStatus := range s.Entries {
			if ipEntryStatus.Entry == entry {
				s.Entries[ix].Status = entryStatus

				return
			}
		}

		s.Entries = append(
			s.Entries,
			IPAccessEntryStatus{
				Entry:  entry,
				Status: entryStatus,
			},
		)
	}
}

func RemoveIPAccessListEntryStatus(entry string) AtlasIPAccessListStatusOption {
	return func(s *AtlasIPAccessListStatus) {
		for ix, ipEntryStatus := range s.Entries {
			if ipEntryStatus.Entry == entry {
				s.Entries = append(s.Entries[:ix], s.Entries[ix+1:]...)

				return
			}
		}
	}
}
