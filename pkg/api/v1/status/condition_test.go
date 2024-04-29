package status

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestHasConditionType(t *testing.T) {
	for _, tc := range []struct {
		name   string
		typ    ConditionType
		source []Condition
		want   bool
	}{
		{
			name:   "nil source",
			typ:    ReadyType,
			source: nil,
			want:   false,
		},
		{
			name:   "empty source",
			typ:    ReadyType,
			source: []Condition{},
			want:   false,
		},
		{
			name: "want ready type",
			typ:  ReadyType,
			source: []Condition{
				{
					Type:   ReadyType,
					Status: corev1.ConditionTrue,
				},
			},
			want: true,
		},
		{
			name: "do not want ready type",
			typ:  ReadyType,
			source: []Condition{
				{
					Type:   ValidationSucceeded,
					Status: corev1.ConditionTrue,
				},
			},
			want: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := HasConditionType(tc.typ, tc.source); got != tc.want {
				t.Errorf("want HasConditionType %v, got %t", tc.want, got)
			}
		})
	}
}
