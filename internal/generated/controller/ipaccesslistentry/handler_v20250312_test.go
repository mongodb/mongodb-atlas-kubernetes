// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipaccesslistentry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312018/admin"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestBuildNetworkPermissionEntry_DeleteAfterDate(t *testing.T) {
	t.Run("valid RFC3339 string is parsed to *time.Time", func(t *testing.T) {
		entry := buildNetworkPermissionEntry(ipal(func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
			e.CidrBlock = pointer.MakePtr("10.0.0.0/8")
			e.DeleteAfterDate = pointer.MakePtr("2025-07-01T12:30:00Z")
		}))
		require.NotNil(t, entry.DeleteAfterDate)
		want := time.Date(2025, 7, 1, 12, 30, 0, 0, time.UTC)
		assert.Equal(t, want, *entry.DeleteAfterDate)
	})

	t.Run("invalid string is silently ignored — DeleteAfterDate stays nil", func(t *testing.T) {
		entry := buildNetworkPermissionEntry(ipal(func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
			e.CidrBlock = pointer.MakePtr("10.0.0.0/8")
			e.DeleteAfterDate = pointer.MakePtr("not-a-date")
		}))
		assert.Nil(t, entry.DeleteAfterDate,
			"unparseable date must not produce a zero-value time — Atlas would reject the request")
	})

	t.Run("nil deleteAfterDate produces nil in output", func(t *testing.T) {
		entry := buildNetworkPermissionEntry(ipal(func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
			e.CidrBlock = pointer.MakePtr("10.0.0.0/8")
		}))
		assert.Nil(t, entry.DeleteAfterDate)
	})
}

func TestBuildNetworkPermissionEntry_CIDRNormalization(t *testing.T) {
	t.Run("host bits are masked to network address", func(t *testing.T) {
		entry := buildNetworkPermissionEntry(ipal(func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
			e.CidrBlock = pointer.MakePtr("192.168.1.5/24")
		}))
		assert.Equal(t, pointer.MakePtr("192.168.1.0/24"), entry.CidrBlock)
	})

	t.Run("already-normalized CIDR is unchanged", func(t *testing.T) {
		entry := buildNetworkPermissionEntry(ipal(func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
			e.CidrBlock = pointer.MakePtr("10.0.0.0/8")
		}))
		assert.Equal(t, pointer.MakePtr("10.0.0.0/8"), entry.CidrBlock)
	})
}

func TestBuildNetworkPermissionEntry_EntryType(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*akov2generated.IPAccessListEntrySpecV20250312Entry)
		check func(*testing.T, v20250312sdk.NetworkPermissionEntry)
	}{
		{
			name: "cidrBlock",
			setup: func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
				e.CidrBlock = pointer.MakePtr("10.0.0.0/8")
				e.Comment = pointer.MakePtr("office network")
			},
			check: func(t *testing.T, got v20250312sdk.NetworkPermissionEntry) {
				assert.Equal(t, pointer.MakePtr("10.0.0.0/8"), got.CidrBlock)
				assert.Equal(t, pointer.MakePtr("office network"), got.Comment)
				assert.Nil(t, got.IpAddress)
				assert.Nil(t, got.AwsSecurityGroup)
			},
		},
		{
			name: "ipAddress",
			setup: func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
				e.IpAddress = pointer.MakePtr("1.2.3.4")
			},
			check: func(t *testing.T, got v20250312sdk.NetworkPermissionEntry) {
				assert.Equal(t, pointer.MakePtr("1.2.3.4"), got.IpAddress)
				assert.Nil(t, got.CidrBlock)
				assert.Nil(t, got.AwsSecurityGroup)
			},
		},
		{
			name: "awsSecurityGroup",
			setup: func(e *akov2generated.IPAccessListEntrySpecV20250312Entry) {
				e.AwsSecurityGroup = pointer.MakePtr("sg-0123456789abcdef0")
			},
			check: func(t *testing.T, got v20250312sdk.NetworkPermissionEntry) {
				assert.Equal(t, pointer.MakePtr("sg-0123456789abcdef0"), got.AwsSecurityGroup)
				assert.Nil(t, got.CidrBlock)
				assert.Nil(t, got.IpAddress)
			},
		},
		{
			name: "nil entry returns empty struct",
			setup: func(_ *akov2generated.IPAccessListEntrySpecV20250312Entry) {
				// leave entry empty
			},
			check: func(t *testing.T, got v20250312sdk.NetworkPermissionEntry) {
				assert.Equal(t, v20250312sdk.NetworkPermissionEntry{}, got)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entry := buildNetworkPermissionEntry(ipal(tc.setup))
			tc.check(t, entry)
		})
	}
}

func ipal(customize func(*akov2generated.IPAccessListEntrySpecV20250312Entry)) *akov2generated.IPAccessListEntry {
	entry := &akov2generated.IPAccessListEntrySpecV20250312Entry{}
	customize(entry)
	return &akov2generated.IPAccessListEntry{
		Spec: akov2generated.IPAccessListEntrySpec{
			V20250312: &akov2generated.IPAccessListEntrySpecV20250312{
				Entry: entry,
			},
		},
	}
}
