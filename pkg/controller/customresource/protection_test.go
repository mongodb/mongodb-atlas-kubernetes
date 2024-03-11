package customresource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
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
	return func(resource akov2.AtlasCustomResource) (bool, error) {
		return reply, nil
	}
}

func testAtlasChecker(reply bool) customresource.AtlasChecker {
	return func(resource akov2.AtlasCustomResource) (bool, error) {
		return reply, nil
	}
}

var ErrOpChecker = fmt.Errorf("operator checker failed")

func failedOpChecker(err error) customresource.OperatorChecker {
	return func(resource akov2.AtlasCustomResource) (bool, error) {
		return false, err
	}
}

var ErrAtlasChecker = fmt.Errorf("atlas checker failed")

func failedAtlasChecker(err error) customresource.AtlasChecker {
	return func(resource akov2.AtlasCustomResource) (bool, error) {
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
		resource          akov2.AtlasCustomResource
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
	resource.Spec.Username = "test-user"

	// ignore the error due to not configuring the fake client
	// we are not checking that, we are only interested on a new annotation in resource
	_ = customresource.ApplyLastConfigApplied(context.Background(), resource, fake.NewClientBuilder().Build())

	annot := resource.GetAnnotations()
	assert.NotEmpty(t, annot)
	expectedConfig := `{"projectRef":{"name":"","namespace":""},"roles":null,"username":"test-user"}`
	assert.Equal(t, annot[customresource.AnnotationLastAppliedConfiguration], expectedConfig)
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
