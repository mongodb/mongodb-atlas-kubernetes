package contract

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func DefaultIPAccessList() []admin.NetworkPermissionEntry {
	anyAccess := "0.0.0.0/0"
	return []admin.NetworkPermissionEntry{{CidrBlock: &anyAccess}}
}

func WithIPAccessList(ipAccessList []admin.NetworkPermissionEntry) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		log.Printf("Setting up IP Access list %s...", display(ipAccessList))
		apiClient, err := NewAPIClient()
		if err != nil {
			return nil, err
		}
		_, _, err = apiClient.ProjectIPAccessListApi.CreateProjectIpAccessList(
			ctx, resources.ProjectID, &ipAccessList).Execute()
		if err != nil {
			return nil, fmt.Errorf("failed to setup IP access list %s: %w", display(ipAccessList), err)
		}
		log.Printf("IP access list %s setup", display(ipAccessList))
		return resources, nil
	}
}

func clearIPAccessList(ctx context.Context, projectID string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	pagIPAccessList, _, err :=
		apiClient.ProjectIPAccessListApi.ListProjectIpAccessLists(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to list IP access list: %w", err)
	}
	if pagIPAccessList.TotalCount == nil || *pagIPAccessList.TotalCount == 0 {
		return nil
	}
	total := *pagIPAccessList.TotalCount
	if *pagIPAccessList.Results == nil {
		return fmt.Errorf("empty list of ip accesses but expected %d items", total)
	}
	ipAccessList := *pagIPAccessList.Results
	if len(ipAccessList) < total {
		return fmt.Errorf("expected %d but got %d items", total, len(ipAccessList))
	}
	for _, ipAccess := range ipAccessList {
		access := decode(ipAccess)
		_, _, err := apiClient.ProjectIPAccessListApi.DeleteProjectIpAccessList(ctx, projectID, access).Execute()
		if err != nil {
			return fmt.Errorf("failed to remove IP access %s: %w", access, err)
		}
		log.Printf("IP access %s removed", access)
	}
	log.Printf("IP access lists cleared")
	return nil
}

func display(ipAccessList []admin.NetworkPermissionEntry) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "[ ")
	for _, ipAccess := range ipAccessList {
		fmt.Fprintf(&buf, "%s ", decode(ipAccess))
	}
	fmt.Fprintf(&buf, " ]")
	return buf.String()
}

func decode(ipAccess admin.NetworkPermissionEntry) string {
	switch {
	case ipAccess.CidrBlock != nil && *ipAccess.CidrBlock != "":
		quads := strings.Split(*ipAccess.CidrBlock, "/")
		if quads[1] == "32" {
			return quads[0]
		}
		return *ipAccess.CidrBlock
	case ipAccess.IpAddress != nil && *ipAccess.IpAddress != "":
		ip := strings.Split(*ipAccess.IpAddress, "/")
		return ip[0]
	case ipAccess.AwsSecurityGroup != nil && *ipAccess.AwsSecurityGroup != "":
		return *ipAccess.AwsSecurityGroup
	default:
		return ""
	}
}
