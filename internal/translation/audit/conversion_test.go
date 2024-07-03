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

func TestToAtlas(t *testing.T) {
	TrueBool := true
	FalseBool := false
	EmptyAudit := "{}"
	testCases := []struct {
		title          string
		input          *AuditConfig
		expectedOutput *admin.AuditLog
	}{
		{
			title: "Just enabled",
			input: NewAuditConfig(
				&akov2.Auditing{
					Enabled: true,
				},
			),
			expectedOutput: &admin.AuditLog{
				Enabled:                   &TrueBool,
				AuditAuthorizationSuccess: &FalseBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "Auth success logs as well",
			input: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
				},
			),
			expectedOutput: &admin.AuditLog{
				AuditAuthorizationSuccess: &TrueBool,
				Enabled:                   &TrueBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "With a filter",
			input: NewAuditConfig(
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{"atype":"authenticate"}`,
				},
			),
			expectedOutput: &admin.AuditLog{
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
			input: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			expectedOutput: &admin.AuditLog{
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
			input: NewAuditConfig(
				&akov2.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			expectedOutput: &admin.AuditLog{
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
			input: NewAuditConfig(
				&akov2.Auditing{},
			),
			expectedOutput: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{}`
					return &s
				}(),
				AuditAuthorizationSuccess: &FalseBool,
				Enabled:                   &FalseBool,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualResult := toAtlas(tc.input)
			assert.Equal(t, tc.expectedOutput, actualResult)
		})
	}
}

func TestFromAtlas(t *testing.T) {
	TrueBool := true
	FalseBool := false
	EmptyAudit := "{}"
	testCases := []struct {
		title          string
		expectedOutput *AuditConfig
		input          *admin.AuditLog
	}{
		{
			title: "Just enabled",
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{
					Enabled: true,
				},
			),
			input: &admin.AuditLog{
				Enabled:                   &TrueBool,
				AuditAuthorizationSuccess: &FalseBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "Auth success logs as well",
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
				},
			),
			input: &admin.AuditLog{
				AuditAuthorizationSuccess: &TrueBool,
				Enabled:                   &TrueBool,
				AuditFilter:               &EmptyAudit,
			},
		},
		{
			title: "With a filter",
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{
					Enabled:     true,
					AuditFilter: `{"atype":"authenticate"}`,
				},
			),
			input: &admin.AuditLog{
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
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{
					Enabled:                   true,
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			input: &admin.AuditLog{
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
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate"}`,
				},
			),
			input: &admin.AuditLog{
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
			expectedOutput: NewAuditConfig(
				&akov2.Auditing{},
			),
			input: &admin.AuditLog{
				AuditFilter: func() *string {
					s := `{}`
					return &s
				}(),
				AuditAuthorizationSuccess: &FalseBool,
				Enabled:                   &FalseBool,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actualResult := fromAtlas(tc.input)
			assert.Equal(t, tc.expectedOutput, actualResult)
		})
	}
}
