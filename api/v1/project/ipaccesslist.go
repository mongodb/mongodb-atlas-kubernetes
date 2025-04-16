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

package project

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
)

type IPAccessList struct {
	// Unique identifier of AWS security group in this access list entry.
	// +optional
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"`
	// Range of IP addresses in CIDR notation in this access list entry.
	// +optional
	CIDRBlock string `json:"cidrBlock,omitempty"`
	// Comment associated with this access list entry.
	// +optional
	Comment string `json:"comment,omitempty"`
	// Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry.
	// +optional
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
	// Entry using an IP address in this access list entry.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
}

// ToAtlas converts the ProjectIPAccessList to native Atlas client format.
func (i IPAccessList) ToAtlas() (*mongodbatlas.ProjectIPAccessList, error) {
	result := &mongodbatlas.ProjectIPAccessList{}
	err := compat.JSONCopy(result, i)
	return result, err
}

// ************************************ Builder methods *************************************************
// Note, that we don't use pointers here as the AtlasProject uses this without pointers

func NewIPAccessList() IPAccessList {
	return IPAccessList{}
}

func (i IPAccessList) WithComment(comment string) IPAccessList {
	i.Comment = comment
	return i
}

func (i IPAccessList) WithIP(ip string) IPAccessList {
	i.IPAddress = ip
	return i
}

func (i IPAccessList) WithCIDR(cidr string) IPAccessList {
	i.CIDRBlock = cidr
	return i
}

func (i IPAccessList) WithAWSGroup(group string) IPAccessList {
	i.AwsSecurityGroup = group
	return i
}

func (i IPAccessList) WithDeleteAfterDate(date string) IPAccessList {
	i.DeleteAfterDate = date
	return i
}
