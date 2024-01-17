package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

//+k8s:deepcopy-gen=false

// AtlasCustomResource is the interface common for all Atlas entities
type AtlasCustomResource interface {
	metav1.Object
	runtime.Object
	api.Reader
	api.Writer
}

var _ AtlasCustomResource = &AtlasProject{}

var _ AtlasCustomResource = &AtlasTeam{}
var _ AtlasCustomResource = &AtlasDeployment{}
var _ AtlasCustomResource = &AtlasDatabaseUser{}
var _ AtlasCustomResource = &AtlasDataFederation{}
var _ AtlasCustomResource = &AtlasBackupSchedule{}
var _ AtlasCustomResource = &AtlasBackupPolicy{}
var _ AtlasCustomResource = &AtlasFederatedAuth{}
