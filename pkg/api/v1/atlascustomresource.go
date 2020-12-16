package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"k8s.io/apimachinery/pkg/runtime"
)

//+k8s:deepcopy-gen=false
// AtlasCustomResource is the interface common for all Atlas entities
type AtlasCustomResource interface {
	runtime.Object
	status.Reader
	status.Writer
}

var _ AtlasCustomResource = &AtlasProject{}
var _ AtlasCustomResource = &AtlasCluster{}
