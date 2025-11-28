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
//

package crapi

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/refs"
)

// Translator allows to translate back and forth between a CRD schema
// and SDK API structures of a certain version.
// A translator is an immutable configuration object, it can be safely shared
// across goroutines
type Translator interface {
	// MajorVersion returns the pinned SDK major version
	MajorVersion() string

	// Mappings returns all the OpenAPi custom reference extensions, or an error
	Mappings() ([]*refs.Mapping, error)

	// ToAPI translates a source Kubernetes object into a target API structure.
	// It uses the spec only to populate ethe API request, nothing from the status.
	// The target is set to a API request struct to be filled.
	// The source is set to the Kubernetes CR value. Only the spec data is used here.
	// The request includes the translator and the dependencies associated with the
	// source CR, usually Kubernetes secrets.
	ToAPI(target any, source client.Object, objs ...client.Object) error

	// FromAPI translates a source API structure into a Kubernetes object.
	// The API source is used to populate the Kubernetes spec, including the
	// spec.entry and status as well.
	// The target is set to CR value to be filled. Both spec and status are filled.
	// The source is set to API response.
	// The request includes the translator and any dependencies associated with the
	// source CR.
	// Returns any extra objects extracted from the response as separate Kubernetes
	// objects, such as Kubernetes secrets, for instance. This list does not include
	// the mutated target, and will be empty if nothing else was extracted off the ç
	// response.
	FromAPI(target client.Object, source any, objs ...client.Object) ([]client.Object, error)
}

// Request is deprecated do not use
//
// Deprecated: request is no longer used in the ToAPI and FromAPI calls
type Request struct {
	Translator   Translator
	Dependencies []client.Object
}
