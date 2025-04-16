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

package customresource_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
)

func sampleResource() *akov2.AtlasDatabaseUser {
	return &akov2.AtlasDatabaseUser{
		Spec: akov2.AtlasDatabaseUserSpec{},
	}
}

func taggedResource(tag, value string) *akov2.AtlasDatabaseUser {
	dbUser := sampleResource()
	annot := map[string]string{}
	annot[tag] = value
	dbUser.SetAnnotations(annot)
	return dbUser
}

func testOpChecker(reply bool) customresource.OperatorChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		return reply, nil
	}
}

func testAtlasChecker(reply bool) customresource.AtlasChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		return reply, nil
	}
}

var ErrOpChecker = fmt.Errorf("operator checker failed")

func failedOpChecker(err error) customresource.OperatorChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		return false, err
	}
}

var ErrAtlasChecker = fmt.Errorf("atlas checker failed")

func failedAtlasChecker(err error) customresource.AtlasChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		return false, err
	}
}

func TestWithoutProtectionIsOwned(t *testing.T) {
	owned, err := customresource.IsOwner(sampleResource(), false, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, owned, true)
}

func TestProtected(t *testing.T) {
	tests := []struct {
		title         string
		opChecker     customresource.OperatorChecker
		atlasChecker  customresource.AtlasChecker
		expectedOwned bool
	}{
		{"managed is owned", testOpChecker(true), nil, true},
		{"unmanaged but not in Atlas is owned", testOpChecker(false), testAtlasChecker(false), true},
		{"unmanaged but in Atlas is NOT owned", testOpChecker(false), testAtlasChecker(true), false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Protected and %s", tc.title), func(t *testing.T) {
			owned, err := customresource.IsOwner(sampleResource(), true, tc.opChecker, tc.atlasChecker)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOwned, owned)
		})
	}
}

func TestProtectedFailures(t *testing.T) {
	tests := []struct {
		title           string
		opChecker       customresource.OperatorChecker
		atlasChecker    customresource.AtlasChecker
		expectedFailure error
	}{
		{"When all checkers fail, operator checker fails first", failedOpChecker(ErrOpChecker), failedAtlasChecker(ErrAtlasChecker), ErrOpChecker},
		{"When unamanaged and atlas checker fails we get that its failure", testOpChecker(false), failedAtlasChecker(ErrAtlasChecker), ErrAtlasChecker},
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			_, err := customresource.IsOwner(sampleResource(), true, tc.opChecker, tc.atlasChecker)
			assert.Equal(t, tc.expectedFailure, err)
		})
	}
}

func TestIsResourceProtected(t *testing.T) {
	tests := []struct {
		title             string
		protectionFlag    bool
		resource          api.AtlasCustomResource
		expectedProtected bool
	}{
		{"Resource without tags with the flag set is protected", true, sampleResource(), true},
		{"Resource without tags with the flag unset isn't protected", false, sampleResource(), false},
		{
			"Resource with keep tag is protected",
			false,
			taggedResource(customresource.ResourcePolicyAnnotation, customresource.ResourcePolicyKeep),
			true,
		},
		{
			"Resource with delete tag and protected flag set is NOT protected",
			true,
			taggedResource(customresource.ResourcePolicyAnnotation, customresource.ResourcePolicyDelete),
			false,
		},
		{
			"Resource with delete tag and protected flag unset isn't protected",
			false,
			taggedResource(customresource.ResourcePolicyAnnotation, customresource.ResourcePolicyDelete),
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.expectedProtected, customresource.IsResourcePolicyKeepOrDefault(tc.resource, tc.protectionFlag))
		})
	}
}

func TestApplyLastConfigApplied(t *testing.T) {
	resource := sampleResource()
	resource.Name = "foo"
	resource.Spec.Username = "test-user"

	scheme := runtime.NewScheme()
	utilruntime.Must(akov2.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithObjects(resource).WithScheme(scheme).Build()
	assert.NoError(t, customresource.ApplyLastConfigApplied(context.Background(), resource, c))

	annot := resource.GetAnnotations()
	assert.NotEmpty(t, annot)
	expectedConfig := `{"roles":null,"username":"test-user"}`
	assert.Equal(t, expectedConfig, annot[customresource.AnnotationLastAppliedConfiguration])
}

func TestIsResourceManagedByOperator(t *testing.T) {
	testCases := []struct {
		title         string
		annotated     bool
		expectManaged bool
	}{
		{
			title:         "If the resource is annotated with last applied config, then it is managed",
			annotated:     true,
			expectManaged: true,
		},
		{
			title:         "If the resource is NOT annotated with last applied config, then it is NOT managed",
			annotated:     true,
			expectManaged: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			resource := sampleResource()
			if tc.annotated {
				customresource.SetAnnotation(resource, customresource.AnnotationLastAppliedConfiguration, "")
			}

			managed, err := customresource.IsResourceManagedByOperator(resource)
			require.NoError(t, err)
			assert.Equal(t, tc.expectManaged, managed)
		})
	}
}

func TestPatchLastConfigApplied(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	utilruntime.Must(akov2.AddToScheme(scheme))
	for _, tc := range []struct {
		title  string
		object *akov2.AtlasProject
		spec   *akov2.AtlasProjectSpec
	}{
		{
			title: "spec without changes is a noop",
			object: &akov2.AtlasProject{
				TypeMeta:   metav1.TypeMeta{Kind: "AtlasProject", APIVersion: "atlas.mongodb.com"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns", Annotations: map[string]string{}},
				Spec: akov2.AtlasProjectSpec{
					Name:                "atlas-project-name",
					ConnectionSecret:    &common.ResourceRefNamespaced{Name: "secret-name"},
					CustomRoles:         []akov2.CustomRole{{}},
					ProjectIPAccessList: []project.IPAccessList{{}},
					PrivateEndpoints:    []akov2.PrivateEndpoint{{}},
					NetworkPeers:        []akov2.NetworkPeer{{}},
					Teams:               []akov2.Team{{}},
				},
				Status: status.AtlasProjectStatus{ID: "some-id"},
			},
			spec: &akov2.AtlasProjectSpec{ // same as object's spec, no changes
				Name:                "atlas-project-name",
				ConnectionSecret:    &common.ResourceRefNamespaced{Name: "secret-name"},
				CustomRoles:         []akov2.CustomRole{{}},
				ProjectIPAccessList: []project.IPAccessList{{}},
				PrivateEndpoints:    []akov2.PrivateEndpoint{{}},
				NetworkPeers:        []akov2.NetworkPeer{{}},
				Teams:               []akov2.Team{{}},
			},
		},

		{
			title: "cleared spec is applied with no other changes",
			object: &akov2.AtlasProject{
				TypeMeta:   metav1.TypeMeta{Kind: "AtlasProject", APIVersion: "atlas.mongodb.com"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns", Annotations: map[string]string{}},
				Spec: akov2.AtlasProjectSpec{
					Name:                "atlas-project-name",
					ConnectionSecret:    &common.ResourceRefNamespaced{Name: "secret-name"},
					CustomRoles:         []akov2.CustomRole{{}},
					ProjectIPAccessList: []project.IPAccessList{{}},
					PrivateEndpoints:    []akov2.PrivateEndpoint{{}},
					NetworkPeers:        []akov2.NetworkPeer{{}},
					Teams:               []akov2.Team{{}},
				},
				Status: status.AtlasProjectStatus{ID: "some-id"},
			},
			spec: &akov2.AtlasProjectSpec{ // clear applied
				Name:             "atlas-project-name",
				ConnectionSecret: &common.ResourceRefNamespaced{Name: "secret-name"},
				Teams:            []akov2.Team{{}},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().WithObjects(tc.object).WithScheme(scheme).Build()

			require.NoError(t, customresource.PatchLastConfigApplied(ctx, k8sClient, tc.object, tc.spec))

			result := akov2.AtlasProject{}
			require.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(tc.object), &result))

			resultSpec, err := customresource.ParseLastConfigApplied[akov2.AtlasProjectSpec](&result)
			require.NoError(t, err)
			assert.Equal(t, tc.spec, resultSpec)

			want := tc.object
			assert.Equal(t, clearProjectToCompare(want), clearProjectToCompare(&result))
		})
	}
}

func clearProjectToCompare(prj *akov2.AtlasProject) *akov2.AtlasProject {
	copy := prj.DeepCopy()
	// ignore the resourceVersion and the last applied config, compared separately
	delete(copy.Annotations, customresource.AnnotationLastAppliedConfiguration)
	copy.ResourceVersion = ""
	return copy
}

func TestPatchLastConfigAppliedErrors(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	utilruntime.Must(akov2.AddToScheme(scheme))
	for _, tc := range []struct {
		title        string
		object       *akov2.AtlasProject
		spec         *struct{}
		wantErrorMsg string
	}{
		{
			title:        "nil spec fails",
			object:       &akov2.AtlasProject{},
			spec:         nil,
			wantErrorMsg: "spec is nil",
		},
		{
			title:        "empty struct cannot be patched",
			object:       &akov2.AtlasProject{},
			spec:         &struct{}{},
			wantErrorMsg: "failed to patch resource:  \"\" is invalid: metadata.name: Required value: name is required",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			k8sClient := fake.NewClientBuilder().WithObjects(tc.object).WithScheme(scheme).Build()

			err := customresource.PatchLastConfigApplied(ctx, k8sClient, tc.object, tc.spec)
			assert.ErrorContains(t, err, tc.wantErrorMsg)
		})
	}
}

func TestParseLastConfigApplied(t *testing.T) {
	for _, tc := range []struct {
		title        string
		object       *akov2.AtlasProject
		want         *akov2.AtlasProjectSpec
		wantErrorMsg string
	}{
		{
			title:  "empty project without annotations renders nothing",
			object: &akov2.AtlasProject{},
		},
		{
			title: "broken JSON in annotation renders error",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: "bad-json",
					},
				},
			},
			wantErrorMsg: "error parsing JSON annotation value into a v1.AtlasProjectSpec",
		},
		{
			title: "proper but empty JSON renders empty spec",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: "{}",
					},
				},
			},
			want: &akov2.AtlasProjectSpec{},
		},
		{
			title: "sample JSON spec renders original",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: jsonize(
							akov2.AtlasProjectSpec{
								Name:                "atlas-project-name",
								ConnectionSecret:    &common.ResourceRefNamespaced{Name: "secret-name"},
								CustomRoles:         []akov2.CustomRole{{}},
								ProjectIPAccessList: []project.IPAccessList{{}},
								PrivateEndpoints:    []akov2.PrivateEndpoint{{}},
								NetworkPeers:        []akov2.NetworkPeer{{}},
								Teams:               []akov2.Team{{}},
							}),
					},
				},
			},
			want: &akov2.AtlasProjectSpec{
				Name:                "atlas-project-name",
				ConnectionSecret:    &common.ResourceRefNamespaced{Name: "secret-name"},
				CustomRoles:         []akov2.CustomRole{{}},
				ProjectIPAccessList: []project.IPAccessList{{}},
				PrivateEndpoints:    []akov2.PrivateEndpoint{{}},
				NetworkPeers:        []akov2.NetworkPeer{{}},
				Teams:               []akov2.Team{{}},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			spec, err := customresource.ParseLastConfigApplied[akov2.AtlasProjectSpec](tc.object)
			if tc.wantErrorMsg != "" {
				assert.ErrorContains(t, err, tc.wantErrorMsg)
			}
			assert.Equal(t, tc.want, spec)
		})
	}
}

func jsonize(obj any) string {
	js, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	return string(js)
}
