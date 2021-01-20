package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//+k8s:deepcopy-gen=false

// AtlasCustomResource is the interface common for all Atlas entities
type AtlasCustomResource interface {
	metav1.Object
	runtime.Object
	status.Reader
	status.Writer
}

var _ AtlasCustomResource = &AtlasProject{}

var _ AtlasCustomResource = &AtlasCluster{}
