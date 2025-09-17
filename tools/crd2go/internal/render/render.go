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

package render

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

type CRDRenderRequest struct {
	gotype.Request
	Filename string
	Version  string
	Kind     string
	Type     *gotype.GoType
}

type CRD2GoRenderer interface {
	// RenderDoc generates the doc.go file from the request, version and group inputs
	RenderDoc(req *gotype.Request, group, version string) error

	// RenderSchema generates the schema.go file from the request, version and group inputs
	RenderSchema(req *gotype.Request, group, version string) error

	// RenderCRD renders each of the CRD Go files form the rewuqest and versioned CRD
	RenderCRD(req *CRDRenderRequest) error
}

var Default = JenRenderer{}
