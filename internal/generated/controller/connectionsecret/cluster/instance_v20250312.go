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

package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/data"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// ClusterInstance_v20250312 is the instance type that implements ConnectionTargetInstance
type ClusterInstance_v20250312 struct {
	Cluster *generatedv1.Cluster
}

func (e *ClusterInstance_v20250312) GetConnectionTargetType() string {
	return "cluster"
}

func (e *ClusterInstance_v20250312) GetName() string {
	if e.Cluster != nil && e.Cluster.Spec.V20250312 != nil && e.Cluster.Spec.V20250312.Entry != nil && e.Cluster.Spec.V20250312.Entry.Name != nil {
		return *e.Cluster.Spec.V20250312.Entry.Name
	}
	return ""
}

func (e *ClusterInstance_v20250312) IsReady() bool {
	ready := meta.FindStatusCondition(e.Cluster.GetConditions(), state.ReadyCondition)
	return ready != nil && ready.Status == metav1.ConditionTrue
}

func (e *ClusterInstance_v20250312) GetScopeType() string {
	return "CLUSTER"
}

func (e *ClusterInstance_v20250312) GetProjectID(ctx context.Context) string {
	if e.Cluster == nil || e.Cluster.Status.V20250312 == nil || e.Cluster.Status.V20250312.GroupId == nil {
		return ""
	}
	return *e.Cluster.Status.V20250312.GroupId
}

func (e *ClusterInstance_v20250312) BuildConnectionData(ctx context.Context) *data.ConnectionSecret {
	if e.Cluster == nil || e.Cluster.Status.V20250312 == nil || e.Cluster.Status.V20250312.ConnectionStrings == nil {
		return nil // no data available
	}

	result := &data.ConnectionSecret{}
	conn := e.Cluster.Status.V20250312.ConnectionStrings

	if e.Cluster.Status.V20250312.ConnectionStrings != nil {
		result.ConnectionURL = ptr.Deref[string](conn.Standard, "")
		result.SrvConnectionURL = ptr.Deref[string](conn.StandardSrv, "")
	}

	return result
}
