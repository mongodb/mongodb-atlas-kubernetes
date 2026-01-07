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

package flexcluster

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/data"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type FlexClusterInstance_v20250312 struct {
	FlexCluster *generatedv1.FlexCluster
}

func (e *FlexClusterInstance_v20250312) GetConnectionTargetType() string {
	return "flexcluster"
}

func (e *FlexClusterInstance_v20250312) GetName() string {
	if e.FlexCluster != nil && e.FlexCluster.Spec.V20250312 != nil && e.FlexCluster.Spec.V20250312.Entry != nil {
		return e.FlexCluster.Spec.V20250312.Entry.Name
	}
	return ""
}

func (e *FlexClusterInstance_v20250312) IsReady() bool {
	ready := meta.FindStatusCondition(e.FlexCluster.GetConditions(), state.ReadyCondition)
	return ready != nil && ready.Status == metav1.ConditionTrue
}

func (e *FlexClusterInstance_v20250312) GetScopeType() string {
	return "CLUSTER"
}

func (e *FlexClusterInstance_v20250312) GetProjectID(ctx context.Context) string {
	if e.FlexCluster == nil || e.FlexCluster.Status.V20250312 == nil || e.FlexCluster.Status.V20250312.GroupId == nil {
		return ""
	}
	return *e.FlexCluster.Status.V20250312.GroupId
}

func (e *FlexClusterInstance_v20250312) BuildConnectionData(ctx context.Context) *data.ConnectionSecret {
	if e.FlexCluster == nil || e.FlexCluster.Status.V20250312 == nil || e.FlexCluster.Status.V20250312.ConnectionStrings == nil {
		return nil // no data available
	}

	result := &data.ConnectionSecret{}
	conn := e.FlexCluster.Status.V20250312.ConnectionStrings

	if e.FlexCluster.Status.V20250312.ConnectionStrings != nil {
		result.ConnectionURL = ptr.Deref[string](conn.Standard, "")
		result.SrvConnectionURL = ptr.Deref[string](conn.StandardSrv, "")
	}

	return result
}
