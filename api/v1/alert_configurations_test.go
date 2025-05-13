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

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestServerlessMetricThreshold(t *testing.T) {
	tests := []struct {
		name      string
		akoData   *MetricThreshold
		atlasData *admin.FlexClusterMetricThreshold
		equal     bool
	}{
		{
			name: "Should be able to parse float Theshold",
			atlasData: &admin.FlexClusterMetricThreshold{
				MetricName: "test",
				Mode:       pointer.MakePtr("test"),
				Operator:   pointer.MakePtr("IN"),
				Threshold:  pointer.MakePtr(3.14),
				Units:      pointer.MakePtr("test"),
			},
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "IN",
				Threshold:  "3.14",
				Units:      "test",
				Mode:       "test",
			},
			equal: true,
		},
		{
			name: "Should be able to parse int Theshold",
			atlasData: &admin.FlexClusterMetricThreshold{
				MetricName: "test",
				Mode:       pointer.MakePtr("test"),
				Operator:   pointer.MakePtr("IN"),
				Threshold:  pointer.MakePtr[float64](3),
				Units:      pointer.MakePtr("test"),
			},
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "IN",
				Threshold:  "3",
				Units:      "test",
				Mode:       "test",
			},
			equal: true,
		},
		{
			name: "Should be false if Theshold is not a number",
			atlasData: &admin.FlexClusterMetricThreshold{
				MetricName: "test",
				Mode:       pointer.MakePtr("test"),
				Operator:   pointer.MakePtr("IN"),
				Threshold:  pointer.MakePtr(3.14),
				Units:      pointer.MakePtr("test"),
			},
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "IN",
				Threshold:  "13InvalidFloat",
				Units:      "test",
				Mode:       "test",
			},
			equal: false,
		},
		{
			name:      "Should be false input is nil",
			atlasData: nil,
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "IN",
				Threshold:  "3.14",
				Units:      "test",
				Mode:       "test",
			},
			equal: false,
		},
		{
			name: "Should be false if operator mismatched",
			atlasData: &admin.FlexClusterMetricThreshold{
				MetricName: "test",
				Mode:       pointer.MakePtr("test"),
				Operator:   pointer.MakePtr("IN"),
				Threshold:  pointer.MakePtr(3.14),
				Units:      pointer.MakePtr("test"),
			},
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "LOWER",
				Threshold:  "3.14",
				Units:      "test",
				Mode:       "test",
			},
			equal: false,
		},
		{
			name: "Should fail if Threshold mismatched",
			atlasData: &admin.FlexClusterMetricThreshold{
				MetricName: "test",
				Mode:       pointer.MakePtr("test"),
				Operator:   pointer.MakePtr("IN"),
				Threshold:  pointer.MakePtr(3.14),
				Units:      pointer.MakePtr("test"),
			},
			akoData: &MetricThreshold{
				MetricName: "test",
				Operator:   "IN",
				Threshold:  "2.718",
				Units:      "test",
				Mode:       "test",
			},
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.akoData.IsEqual(tt.atlasData))
		})
	}
}
