package v1

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	internalcmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestSpecEquality(t *testing.T) {
	ref := &AtlasProjectSpec{
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
						Roles:       []string{"1", "2", "3"},
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
				Matchers:      []Matcher{},
			},
			{
				Enabled:       true,
				EventTypeName: "foo",
				Matchers:      nil,
			},
			{
				Enabled:       true,
				EventTypeName: "foo",
			},
		},
	}

	perm := ref.DeepCopy()
	for i := 0; i < 1_000; i++ {
		perm.AlertConfigurations = permute(perm.AlertConfigurations)

		for i := range perm.AlertConfigurations {
			perm.AlertConfigurations[i].Notifications = permute(perm.AlertConfigurations[i].Notifications)

			for j := range perm.AlertConfigurations[i].Notifications {
				perm.AlertConfigurations[i].Notifications[j].Roles = permute(perm.AlertConfigurations[i].Notifications[j].Roles)
			}
		}

		for i := range perm.AlertConfigurations {
			perm.AlertConfigurations[i].Matchers = permute(perm.AlertConfigurations[i].Matchers)
		}

		if !internalcmp.SemanticEqual(ref, perm) {
			jRef, _ := json.MarshalIndent(ref, "", "\t")
			jPermutedCopy, _ := json.MarshalIndent(perm, "", "\t")
			t.Errorf("expected reference:\n%v\nto be equal to the reordered copy:\n%v\nbut it isn't, diff:\n%v",
				string(jRef), string(jPermutedCopy), cmp.Diff(string(jRef), string(jPermutedCopy)),
			)
			return
		}
	}
}

func permute[T any](in []T) []T {
	if len(in) == 0 {
		return nil
	}
	result := make([]T, len(in))
	for i, j := range rand.Perm(len(in)) {
		result[i] = in[j]
	}
	return result
}
