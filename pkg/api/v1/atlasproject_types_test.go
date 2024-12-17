package v1

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"sigs.k8s.io/yaml"

	internalcmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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

func TestLastSpecFrom(t *testing.T) {
	tests := map[string]struct {
		annotations      map[string]string
		expectedLastSpec *AtlasProjectSpec
		expectedError    string
	}{

		"should return nil when there is no last spec": {},
		"should return error when last spec annotation is wrong": {
			annotations: map[string]string{"mongodb.com/last-applied-configuration": "{wrong}"},
			expectedError: "error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]:" +
				" invalid character 'w' looking for beginning of object key string",
		},
		"should return last spec": {
			annotations: map[string]string{"mongodb.com/last-applied-configuration": "{\"name\": \"my-project\"}"},
			expectedLastSpec: &AtlasProjectSpec{
				Name: "my-project",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &AtlasProject{}
			p.WithAnnotations(tt.annotations)
			lastSpec, err := p.LastSpecFrom("mongodb.com/last-applied-configuration")
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedLastSpec, lastSpec)
		})
	}
}
