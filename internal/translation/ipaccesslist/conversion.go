package ipaccesslist

import (
	"fmt"
	"net"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

type IPAccessEntry struct {
	CIDR             string
	AWSSecurityGroup string
	DeleteAfterDate  *time.Time
	Comment          string
}

func (i *IPAccessEntry) ID() string {
	if i.CIDR != "" {
		return i.CIDR
	}

	return i.AWSSecurityGroup
}

func (i *IPAccessEntry) IsExpired(at time.Time) bool {
	if i.DeleteAfterDate == nil {
		return false
	}

	return i.DeleteAfterDate.Before(at)
}

type IPAccessEntries map[string]*IPAccessEntry

func (i IPAccessEntries) GetByStatus(expired bool) IPAccessEntries {
	entries := make(IPAccessEntries, len(i))

	for ix, entry := range i {
		if entry.IsExpired(time.Now()) == expired {
			entries[ix] = entry
		}
	}

	return entries
}

func FromInternal(ipAccessEntries IPAccessEntries) []project.IPAccessList {
	list := make([]project.IPAccessList, 0, len(ipAccessEntries))

	for _, entry := range ipAccessEntries {
		ipAccessList := project.IPAccessList{
			AwsSecurityGroup: entry.AWSSecurityGroup,
			CIDRBlock:        entry.CIDR,
			Comment:          entry.Comment,
		}

		if strings.HasSuffix(entry.CIDR, "/32") {
			ipAccessList.IPAddress = strings.TrimRight(entry.CIDR, "/32")
		}

		if entry.DeleteAfterDate != nil {
			ipAccessList.DeleteAfterDate = timeutil.FormatISO8601(*entry.DeleteAfterDate)
		}

		list = append(list, ipAccessList)
	}

	return list
}

func NewIPAccessEntries(ipAccessList []project.IPAccessList) (IPAccessEntries, error) {
	entries := make(IPAccessEntries, len(ipAccessList))
	for i := range ipAccessList {
		entry, err := newIPAccessEntry(ipAccessList[i])
		if err != nil {
			return nil, err
		}

		entries[entry.ID()] = entry
	}

	return entries, nil
}

func newIPAccessEntry(ipAccessList project.IPAccessList) (*IPAccessEntry, error) {
	cidr, err := parseIPNetwork(ipAccessList.IPAddress, ipAccessList.CIDRBlock)
	if err != nil {
		return nil, err
	}

	var deleteAfterDate *time.Time
	if ipAccessList.DeleteAfterDate != "" {
		dt, err := timeutil.ParseISO8601(ipAccessList.DeleteAfterDate)
		if err != nil {
			return nil, err
		}

		deleteAfterDate = &dt
	}

	return &IPAccessEntry{
		CIDR:             cidr,
		AWSSecurityGroup: ipAccessList.AwsSecurityGroup,
		DeleteAfterDate:  deleteAfterDate,
		Comment:          ipAccessList.Comment,
	}, nil
}

func toAtlas(ipAccessEntries IPAccessEntries) *[]admin.NetworkPermissionEntry {
	netPermissions := make([]admin.NetworkPermissionEntry, 0, len(ipAccessEntries))
	for i := range ipAccessEntries {
		entry := ipAccessEntries[i]
		netPerm := admin.NetworkPermissionEntry{
			DeleteAfterDate: entry.DeleteAfterDate,
		}

		if entry.AWSSecurityGroup != "" {
			netPerm.SetAwsSecurityGroup(entry.AWSSecurityGroup)
		}

		if entry.CIDR != "" {
			netPerm.SetCidrBlock(entry.CIDR)
		}

		if entry.Comment != "" {
			netPerm.SetComment(entry.Comment)
		}

		netPermissions = append(netPermissions, netPerm)
	}

	return &netPermissions
}

func fromAtlas(netPermissions []admin.NetworkPermissionEntry) IPAccessEntries {
	entries := make(IPAccessEntries, len(netPermissions))
	for i := range netPermissions {
		netPermission := netPermissions[i]
		entry := &IPAccessEntry{
			CIDR:             netPermission.GetCidrBlock(),
			AWSSecurityGroup: netPermission.GetAwsSecurityGroup(),
			DeleteAfterDate:  netPermission.DeleteAfterDate,
			Comment:          netPermission.GetComment(),
		}

		entries[entry.ID()] = entry
	}

	return entries
}

func parseIPNetwork(ip, cidr string) (string, error) {
	if ip != "" {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return "", fmt.Errorf("ip %s is invalid", ip)
		}

		ipNet := net.IPNet{
			IP:   parsedIP,
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}
		return ipNet.String(), nil
	}

	if cidr != "" {
		var err error
		_, parsedNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return "", fmt.Errorf("cidr %s is invalid: %w", cidr, err)
		}

		return parsedNet.String(), nil
	}

	return "", nil
}
