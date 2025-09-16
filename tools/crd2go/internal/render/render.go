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
