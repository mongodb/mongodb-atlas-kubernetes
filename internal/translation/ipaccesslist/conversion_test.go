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

package ipaccesslist

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

func TestParseIPNetwork(t *testing.T) {
	test := map[string]struct {
		ip       string
		cidr     string
		expected string
		err      error
	}{
		"should return empty when no ip or cidr is defined": {
			ip:       "",
			cidr:     "",
			expected: "",
		},
		"should return ip network when ip is set": {
			ip:       "192.168.100.150",
			cidr:     "",
			expected: "192.168.100.150/32",
		},
		"should return ip network when cidr is set": {
			ip:       "",
			cidr:     "192.168.100.0/24",
			expected: "192.168.100.0/24",
		},
		"should return ip network when both ip and cidr are set": {
			ip:       "192.168.100.150",
			cidr:     "192.168.100.0/24",
			expected: "192.168.100.150/32",
		},
		"should return error for invalid ip": {
			ip:       "wrong-ip",
			cidr:     "",
			expected: "",
			err:      errors.New("ip wrong-ip is invalid"),
		},
		"should return error for invalid cdir": {
			ip:       "",
			cidr:     "wrong-cidr",
			expected: "",
			err:      fmt.Errorf("cidr wrong-cidr is invalid: %w", &net.ParseError{Type: "CIDR address", Text: "wrong-cidr"}),
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			ipNet, err := parseIPNetwork(tt.ip, tt.cidr)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, ipNet)
		})
	}
}

func TestNewIPAccessEntry(t *testing.T) {
	deleteAfter := timeutil.MustParseISO8601("2024-06-27T13:30:15.999Z")

	test := map[string]struct {
		ipAccessList project.IPAccessList
		expected     *IPAccessEntry
		err          error
	}{
		"should create an ip entry": {
			ipAccessList: project.IPAccessList{
				IPAddress: "192.168.100.150",
			},
			expected: &IPAccessEntry{
				CIDR: "192.168.100.150/32",
			},
		},
		"should create a cidr entry": {
			ipAccessList: project.IPAccessList{
				CIDRBlock: "192.168.100.0/24",
			},
			expected: &IPAccessEntry{
				CIDR: "192.168.100.0/24",
			},
		},
		"should create an aws security group entry": {
			ipAccessList: project.IPAccessList{
				AwsSecurityGroup: "sg-12345",
			},
			expected: &IPAccessEntry{
				AWSSecurityGroup: "sg-12345",
			},
		},
		"should create a temporary entry": {
			ipAccessList: project.IPAccessList{
				IPAddress:       "192.168.100.150",
				DeleteAfterDate: "2024-06-27T13:30:15.999Z",
			},
			expected: &IPAccessEntry{
				CIDR:            "192.168.100.150/32",
				DeleteAfterDate: &deleteAfter,
			},
		},
		"should create an entry with comment": {
			ipAccessList: project.IPAccessList{
				IPAddress: "192.168.100.150",
				Comment:   "My Private IP Address",
			},
			expected: &IPAccessEntry{
				CIDR:    "192.168.100.150/32",
				Comment: "My Private IP Address",
			},
		},
		"should fail when network data is wrong": {
			ipAccessList: project.IPAccessList{
				IPAddress: "192.168.100.300",
			},
			err: errors.New("ip 192.168.100.300 is invalid"),
		},
		"should fail when date time string is wrong": {
			ipAccessList: project.IPAccessList{
				IPAddress:       "192.168.100.128",
				DeleteAfterDate: "2024-07-27 14:02",
			},
			err: &time.ParseError{Layout: "2006-01-02T15:04:05.999Z", Value: "2024-07-27 14:02", LayoutElem: "T", ValueElem: " 14:02", Message: ""},
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			entry, err := newIPAccessEntry(tt.ipAccessList)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, entry)
		})
	}
}

func TestNewIPAccessEntries(t *testing.T) {
	deleteAfter := timeutil.MustParseISO8601("2024-06-27T13:30:15.999Z")

	test := map[string]struct {
		ipAccessList []project.IPAccessList
		expected     IPAccessEntries
		err          error
	}{
		"should convert ip access list": {
			ipAccessList: []project.IPAccessList{
				{
					IPAddress:       "192.168.100.150",
					DeleteAfterDate: "2024-06-27T13:30:15.999Z",
				},
				{
					CIDRBlock: "192.168.1.0/24",
					Comment:   "My network",
				},
				{
					AwsSecurityGroup: "sg-12345",
				},
			},
			expected: IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &deleteAfter,
				},
				"192.168.1.0/24": {
					CIDR:    "192.168.1.0/24",
					Comment: "My network",
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			},
		},
		"should fail to convert ip access list": {
			ipAccessList: []project.IPAccessList{
				{
					IPAddress:       "192.168.100.700",
					DeleteAfterDate: "2024-06-27T13:30:15.999Z",
				},
				{
					CIDRBlock: "192.168.1.0/24",
					Comment:   "My network",
				},
				{
					AwsSecurityGroup: "sg-12345",
				},
			},
			err: errors.New("ip 192.168.100.700 is invalid"),
		},
	}
	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			entry, err := NewIPAccessEntries(tt.ipAccessList)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, entry)
		})
	}
}

func TestIPAccessEntry_ID(t *testing.T) {
	test := map[string]struct {
		ipAccessEntry *IPAccessEntry
		expected      string
	}{
		"should return cidr as identifier": {
			ipAccessEntry: &IPAccessEntry{
				CIDR: "192.168.1.0/24",
			},
			expected: "192.168.1.0/24",
		},
		"should return aws security group as identifier": {
			ipAccessEntry: &IPAccessEntry{
				AWSSecurityGroup: "sg-12345",
			},
			expected: "sg-12345",
		},
		"should return cidr over aws security group as identifier when both are set": {
			ipAccessEntry: &IPAccessEntry{
				CIDR:             "192.168.1.0/24",
				AWSSecurityGroup: "sg-12345",
			},
			expected: "192.168.1.0/24",
		},
		"should return empty when no identifier field is set": {
			ipAccessEntry: &IPAccessEntry{
				Comment: "my config",
			},
			expected: "",
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.ipAccessEntry.ID())
		})
	}
}

func TestIPAccessEntry_IsExpired(t *testing.T) {
	expired := time.Now().UTC().Add(time.Minute * -1)
	active := time.Now().UTC().Add(time.Minute * 1)

	test := map[string]struct {
		ipAccessEntry *IPAccessEntry
		expected      bool
	}{
		"should be expired": {
			ipAccessEntry: &IPAccessEntry{
				CIDR:            "192.168.1.0/24",
				DeleteAfterDate: &expired,
			},
			expected: true,
		},
		"should be active": {
			ipAccessEntry: &IPAccessEntry{
				CIDR:            "192.168.1.0/24",
				DeleteAfterDate: &active,
			},
			expected: false,
		},
		"should be active when no expiration is set": {
			ipAccessEntry: &IPAccessEntry{
				CIDR: "192.168.1.0/24",
			},
			expected: false,
		},
	}

	for name, tt := range test {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.ipAccessEntry.IsExpired(time.Now()))
		})
	}
}

func TestIPAccessEntries(t *testing.T) {
	expired := time.Now().UTC().Add(time.Minute * -1)
	active := time.Now().UTC().Add(time.Minute * 1)
	tests := map[string]struct {
		expired  bool
		expected IPAccessEntries
	}{
		"should filter expired entries": {
			expired: true,
			expected: IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &expired,
				},
			},
		},
		"should filter active entries": {
			expired: false,
			expected: IPAccessEntries{
				"192.168.1.0/24": {
					CIDR:            "192.168.1.0/24",
					DeleteAfterDate: &active,
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			entries := IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &expired,
				},
				"192.168.1.0/24": {
					CIDR:            "192.168.1.0/24",
					DeleteAfterDate: &active,
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			}
			assert.Equal(t, tt.expected, entries.GetByStatus(tt.expired))
		})
	}
}

func TestToAKO(t *testing.T) {
	t.Run("should convert to AKO", func(t *testing.T) {
		expired := time.Now().UTC().Add(time.Minute * -1)
		entries := IPAccessEntries{
			"192.168.100.150/32": {
				CIDR:            "192.168.100.150/32",
				DeleteAfterDate: &expired,
			},
			"192.168.1.0/24": {
				CIDR:    "192.168.1.0/24",
				Comment: "My Network",
			},
			"sg-12345": {
				AWSSecurityGroup: "sg-12345",
			},
		}

		assert.ElementsMatch(
			t,
			[]project.IPAccessList{
				{
					IPAddress:       "192.168.100.150",
					CIDRBlock:       "192.168.100.150/32",
					DeleteAfterDate: timeutil.FormatISO8601(expired),
				},
				{
					CIDRBlock: "192.168.1.0/24",
					Comment:   "My Network",
				},
				{
					AwsSecurityGroup: "sg-12345",
				},
			},
			FromInternal(entries),
		)
	})
}
