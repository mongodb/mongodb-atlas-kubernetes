package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
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
var _ AtlasCustomResource = &AtlasTeam{}
var _ AtlasCustomResource = &AtlasDeployment{}
var _ AtlasCustomResource = &AtlasDatabaseUser{}
var _ AtlasCustomResource = &AtlasDataFederation{}
var _ AtlasCustomResource = &AtlasBackupSchedule{}
var _ AtlasCustomResource = &AtlasBackupPolicy{}
var _ AtlasCustomResource = &AtlasFederatedAuth{}
var _ AtlasCustomResource = &AtlasStreamInstance{}
var _ AtlasCustomResource = &AtlasStreamConnection{}
var _ AtlasCustomResource = &AtlasSearchIndexConfig{}

// InitCondition initializes the underlying type of the given condition to the given default value
// if the underlying condition type is unset.
func InitCondition(resource AtlasCustomResource, defaultCondition status.Condition) []status.Condition {
	conditions := resource.GetStatus().GetConditions()
	if !status.HasConditionType(defaultCondition.Type, conditions) {
		return status.EnsureConditionExists(defaultCondition, conditions)
	}
	return conditions
}
