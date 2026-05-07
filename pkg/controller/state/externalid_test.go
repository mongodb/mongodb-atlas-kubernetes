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

package state_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

func TestGetExternalID(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
		wantID      string
		wantErrMsg  string
	}{
		{
			name:        "returns id when annotation is present",
			annotations: map[string]string{ctrlstate.AnnotationExternalID: "abc123"},
			wantID:      "abc123",
		},
		{
			name:       "returns error when annotation is missing",
			wantErrMsg: ctrlstate.AnnotationExternalID,
		},
		{
			name:        "returns error when annotation value is empty",
			annotations: map[string]string{ctrlstate.AnnotationExternalID: ""},
			wantErrMsg:  ctrlstate.AnnotationExternalID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj := &metav1.ObjectMeta{Annotations: tc.annotations}
			id, err := ctrlstate.GetExternalID(obj)
			if tc.wantErrMsg != "" {
				require.ErrorContains(t, err, tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantID, id)
		})
	}
}
