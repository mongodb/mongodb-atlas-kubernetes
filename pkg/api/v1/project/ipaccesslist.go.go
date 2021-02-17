package project

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
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

// Identifier returns the "id" of the ProjectIPAccessList. Note, that it's an error to specify more than one of these
// fields - the business layer must validate this beforehand
func (i IPAccessList) Identifier() interface{} {
	return i.CIDRBlock + i.AwsSecurityGroup + i.IPAddress
}
