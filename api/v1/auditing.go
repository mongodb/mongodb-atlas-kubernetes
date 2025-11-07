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

package v1

// Auditing represents MongoDB Maintenance Windows
type Auditing struct {
	// Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event.
	// For more information, see auditAuthorizationSuccess.
	// +optional
	AuditAuthorizationSuccess bool `json:"auditAuthorizationSuccess,omitempty"`
	// JSON-formatted audit filter used by the project.
	// +optional
	AuditFilter string `json:"auditFilter,omitempty"`
	// Denotes whether the project associated with the {GROUP-ID} has database auditing enabled.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}
