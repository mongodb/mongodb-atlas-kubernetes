package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestConversionContructor(t *testing.T) {
	testCases := []struct {
		title          string
		input          *akov2.Auditing
		expectedOutput *AuditConfig
	}{
		{
			title: "Just enabled",
			input: &akov2.Auditing{
				Enabled: true,
			},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{}`,
				},
			},
		},
		{
			title: "Auth success logs as well",
			input: &akov2.Auditing{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
			},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{}`,
				},
			},
		},
		{
			title: "With a filter",
			input: &akov2.Auditing{
				Enabled:     true,
				AuditFilter: `{"atype":"authenticate"}`,
			},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{"atype":"authenticate"}`,
				},
			},
		},
		{
			title: "With a filter and success logs",
			input: &akov2.Auditing{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
				AuditFilter:               `{"atype":"authenticate"}`,
			},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			},
		},
		{
			title: "All set but disabled",
			input: &akov2.Auditing{
				AuditAuthorizationSuccess: true,
				AuditFilter:               `{"atype":"authenticate"}`,
			},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			},
		},
		{
			title: "Default (disabled) case",
			input: &akov2.Auditing{},
			expectedOutput: &AuditConfig{
				&akov2.Auditing{
					AuditFilter: `{}`,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualResult := NewAuditConfig(tc.input)
			assert.Equal(t, tc.expectedOutput, actualResult)
		})
	}
}

func TestToAndFromAtlas(t *testing.T) {
	TrueBool := true
	FalseBool := false
	EmptyAudit := "{}"
	testCases := []struct {
		title          string
		kubernetesSide *AuditConfig
		atlasSide      *admin.AuditLog
	}{
		{
			title: "Just enabled",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{
					Enabled: true,
				},
			),
			atlasSide: &admin.AuditLog{
				Enabled:                   &TrueBool,
				AuditAuthorizationSuccess: &FalseBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "Auth success logs as well",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
				},
			),
			atlasSide: &admin.AuditLog{
				AuditAuthorizationSuccess: &TrueBool,
				Enabled:                   &TrueBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "With a filter",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{"atype":"authenticate"}`,
				},
			),
			atlasSide: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{"atype":"authenticate"}`
					return &s
				}(),
				Enabled:                   &TrueBool,
				AuditAuthorizationSuccess: &FalseBool,
			},
		},
		{
			title: "With a filter and success logs",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			atlasSide: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{"atype":"authenticate"}`
					return &s
				}(),
				Enabled:                   &TrueBool,
				AuditAuthorizationSuccess: &TrueBool,
			},
		},
		{
			title: "All set but disabled",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			atlasSide: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{"atype":"authenticate"}`
					return &s
				}(),
				AuditAuthorizationSuccess: &TrueBool,
				Enabled:                   &FalseBool,
			},
		},
		{
			title: "Default (disabled) case",
			kubernetesSide: NewAuditConfig(
				&akov2.Auditing{},
			),
			atlasSide: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{}`
					return &s
				}(),
				AuditAuthorizationSuccess: &FalseBool,
				Enabled:                   &FalseBool,
			},
		},
	}
	//from Kuberenets to Atlas
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualResult := toAtlas(tc.kubernetesSide)
			assert.Equal(t, tc.atlasSide, actualResult)
		})
	}
	//from Atlas to Kubernetes
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualResult := fromAtlas(tc.atlasSide)
			assert.Equal(t, tc.kubernetesSide, actualResult)
		})
	}
}
