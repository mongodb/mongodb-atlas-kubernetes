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

package objmap_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin2025 "go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/objmap"
)

type testStruct struct {
	ID        int
	YesNo     bool
	Text      string
	Data      float32
	Timestamp time.Time
	SubStruct []subStruct
}

type subStruct struct {
	SubID   int
	YesNo   bool
	Text    string
	Data    float32
	TheTime time.Time
}

var sample = testStruct{
	ID:        12,
	YesNo:     true,
	Text:      "some text",
	Data:      2.4,
	Timestamp: time.Now().UTC().Round(time.Nanosecond),
	SubStruct: []subStruct{
		{
			SubID:   15,
			YesNo:   true,
			Text:    "more text",
			Data:    0.67,
			TheTime: time.Now().Add(-24 * 7 * time.Hour).UTC().Round(time.Nanosecond),
		},
		{
			SubID:   14,
			YesNo:   false,
			Text:    "and even more text",
			Data:    1.67,
			TheTime: time.Now().Add(-24 * 15 * time.Hour).UTC().Round(time.Nanosecond),
		},
	},
}

func TestParamsFill(t *testing.T) {
	objMapSample := map[string]any{
		"groupid": "62b6e34b3d91647abb20e7b8",
		"groupalertsconfig": map[string]any{
			"enabled":       true,
			"eventtypename": "some-event",
		},
	}
	result := admin2025.CreateAlertConfigApiParams{}
	require.NoError(t, objmap.FromObjectMap(&result, objMapSample))
	assert.Equal(t, admin2025.CreateAlertConfigApiParams{
		GroupId: "62b6e34b3d91647abb20e7b8",
		GroupAlertsConfig: &admin2025.GroupAlertsConfig{
			Enabled:       pointer.MakePtr(true),
			EventTypeName: pointer.MakePtr("some-event"),
		},
	}, result)
}

func TestToAndFromAndCopyObjectMap(t *testing.T) {
	sampleMap, err := objmap.ToObjectMap(sample)
	require.NoError(t, err)
	clone := map[string]any{}
	objmap.CopyFields(clone, sampleMap)
	result := testStruct{}
	require.NoError(t, objmap.FromObjectMap(&result, sampleMap))
	assert.Equal(t, sample, result)
}
