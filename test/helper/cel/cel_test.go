/*
Copyright 2024 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cel_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel/fake"
)

func TestCEL(t *testing.T) {
	testCases := []struct {
		name         string
		current, old runtime.Object
		wantErrs     []string
	}{
		{
			// Note: It would be desirable if this case failed as well.
			// This will become possible with CRD ratcheting and "optionalOldSelf: true" in the CRD declaration.
			// See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-optional-oldself.
			name: "creating a fake Resource with deprecated set values succeeds",
			old:  nil,
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo"},
				},
			},
		},
		{
			name: "updating a fake Resource and adding an empty deprecated set field succeeds",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{},
				},
			},
		},
		{
			name: "updating a fake Resource and adding a deprecated set with values fails",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo"},
				},
			},
			wantErrs: []string{
				`spec: Invalid value: "object": setting new deprecated set values is invalid: use the NewThing CRD instead.`,
			},
		},
		{
			name: "updating a fake Resource with an empty deprecated set field and adding values to it fails",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{},
				},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo"},
				},
			},
			wantErrs: []string{
				`spec: Invalid value: "object": setting new deprecated set values is invalid: use the NewThing CRD instead.`,
			},
		},
		{
			name: "updating a fake Resource with existing deprecated set values and adding a custom role succeeds",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo"},
				},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo", "bar"},
				},
			},
		},
		{
			name: "updating a fake Resource with existing deprecated set values and removing a custom role succeeds",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo", "bar"},
				},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo"},
				},
			},
		},
		{
			name: "updating a fake Resource with existing deprecated set values and removing all deprecated set values succeeds",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo", "bar"},
				},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{},
				},
			},
		},
		{
			name: "updating a fake Resource with existing deprecated set values and removing the deprecated set values field succeeds",
			old: &fake.Resource{
				Spec: fake.ResourceSpec{
					DeprecatedSet: []string{"foo", "bar"},
				},
			},
			current: &fake.Resource{
				Spec: fake.ResourceSpec{},
			},
		},
	}

	validator, err := cel.VersionValidatorFromFile(t, "./fake/test.mongodb.com_resources.yaml", "fake")
	require.NoError(t, err)
	_ = validator // why?
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				current, old map[string]interface{}
			)
			if tc.current != nil {
				current, err = runtime.DefaultUnstructuredConverter.ToUnstructured(tc.current)
				require.NoError(t, err)
			}
			if tc.old != nil {
				old, err = runtime.DefaultUnstructuredConverter.ToUnstructured(tc.old)
				require.NoError(t, err)
			}
			errs := validator(current, old)

			if got := len(errs); got != len(tc.wantErrs) {
				t.Errorf("expected errors %v, got %v", len(tc.wantErrs), len(errs))
				return
			}

			for i := range tc.wantErrs {
				got := errs[i].Error()
				if got != tc.wantErrs[i] {
					t.Errorf("want error %q, got %q", tc.wantErrs[i], got)
				}
			}
		})
	}
}
