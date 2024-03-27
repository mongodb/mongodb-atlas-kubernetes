package v1

import (
	"reflect"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/google/go-cmp/cmp"
	internalcmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"sigs.k8s.io/yaml"
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

	t.Log(mustMarshal(t, ref))
	internalcmp.Normalize(ref)
	t.Log(mustMarshal(t, ref))
	for i := 0; i < 1_000; i++ {
		perm := ref.DeepCopy()
		internalcmp.PermuteOrder(perm)
		internalcmp.Normalize(perm)

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
