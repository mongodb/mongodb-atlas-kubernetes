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
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	internalcmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestSpecEquality(t *testing.T) {
	ref := &AtlasProjectSpec{
		PrivateEndpoints: []PrivateEndpoint{
			{
				Endpoints: GCPEndpoints{
					{
						EndpointName: "foo",
						IPAddress:    "bar",
					},
					{
						EndpointName: "123",
						IPAddress:    "456",
					},
				},
			},
		},
		AlertConfigurations: []AlertConfiguration{
			{
				Enabled:       true,
				EventTypeName: "foo",
				Notifications: []Notification{
					{
						APITokenRef: common.ResourceRefNamespaced{
							Name: "foo",
						},
						ChannelName: "bar",
						DelayMin:    admin.PtrInt(1),
					},
					{
						ChannelName: "foo",
						DelayMin:    admin.PtrInt(2),
						Roles:       []string{"2", "3", "1"},
					},
					{
						ChannelName: "foo",
						DelayMin:    admin.PtrInt(2),
					},
					{
						APITokenRef: common.ResourceRefNamespaced{
							Name: "bar",
						},
						ChannelName: "bar",
						DelayMin:    admin.PtrInt(1),
					},
				},
			},
			{
				Enabled:       true,
				EventTypeName: "foo",
				Matchers: []Matcher{
					{
						FieldName: "foo",
					},
					{
						FieldName: "bar",
						Operator:  "foo",
					},
					{
						FieldName: "bar",
						Operator:  "bar",
					},
					{
						FieldName: "baz",
						Operator:  "foo",
					},
				},
			},
			{
				Enabled:       true,
				EventTypeName: "foo",
			},
			{
				Enabled:       true,
				EventTypeName: "foo",
			},
		},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	err := internalcmp.Normalize(ref)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100_000; i++ {
		perm := ref.DeepCopy()
		internalcmp.PermuteOrder(perm, r)
		err := internalcmp.Normalize(perm)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ref, perm) {
			jRef := mustMarshal(t, ref)
			jPermutedCopy := mustMarshal(t, perm)
			t.Errorf("expected reference:\n%v\nto be equal to the reordered copy:\n%v\nbut it isn't, diff:\n%v",
				jRef, jPermutedCopy, cmp.Diff(jRef, jPermutedCopy),
			)
			return
		}
	}
}

func mustMarshal(t *testing.T, what any) string {
	t.Helper()
	result, err := yaml.Marshal(what)
	if err != nil {
		t.Fatal(err)
	}
	return string(result)
}

func TestProjectCELChecks(t *testing.T) {
	for _, tc := range []struct {
		title          string
		old, obj       *AtlasProject
		expectedErrors []string
	}{
		{
			title: "Cannot rename a project",
			old: &AtlasProject{
				Spec: AtlasProjectSpec{
					Name: "name-old",
				},
			},
			obj: &AtlasProject{
				Spec: AtlasProjectSpec{
					Name: "name-new",
				},
			},
			expectedErrors: []string{"spec.name: Invalid value: \"string\": Name cannot be modified after project creation"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.old)
			require.NoError(t, err)
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasprojects.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, unstructuredOldObject)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
