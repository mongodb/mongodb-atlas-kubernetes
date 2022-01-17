package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/timeutil"
)

// atlasProjectIPAccessList is a synonym of Atlas object as we need to implement 'Identifier' (and we cannot modify
// their object)
type atlasProjectIPAccessList mongodbatlas.ProjectIPAccessList

func (i atlasProjectIPAccessList) Identifier() interface{} {
	// hack: Atlas adds the CIDRBlock in case IPAddress is specified in the response.
	// This doesn't conform to the "update" contract (one field per List) and doesn't allow to "merge" lists.
	// So we ignore the CIDRblock in this case.
	// Note, this used to have "&& strings.HasPrefix(i.CIDRBlock, i.IPAddress" check as well but according to the example:
	// https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/#example-body the IP may
	// be not a prefix of CIDR!
	if i.CIDRBlock != "" && i.IPAddress != "" {
		return i.AwsSecurityGroup + i.IPAddress
	}
	return i.CIDRBlock + i.AwsSecurityGroup + i.IPAddress
}

// ensureIPAccessList ensures that the state of the Atlas IP Access List matches the
// state of the IP Access list specified in the project CR. Any Access Lists which exist
// in Atlas but are not specified in the CR are deleted.
func ensureIPAccessList(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	if err := validateIPAccessLists(project.Spec.ProjectIPAccessList); err != nil {
		return workflow.Terminate(workflow.ProjectIPAccessInvalid, err.Error())
	}
	active, expired := filterActiveIPAccessLists(project.Spec.ProjectIPAccessList)

	if result := createOrDeleteInAtlas(ctx.Client, projectID, active, ctx.Log); !result.IsOk() {
		return result
	}
	ctx.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(expired))
	return workflow.OK()
}

func validateIPAccessLists(ipAccessList []project.IPAccessList) error {
	for _, list := range ipAccessList {
		if err := validateSingleIPAccessList(list); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleIPAccessList performs validation of the IP access list. Note, that we intentionally don't validate
// IP addresses or CIDR blocks - this will be done by Atlas. But we need to validate the timestamp as we use it to filter
// active and expired ip access lists.
func validateSingleIPAccessList(list project.IPAccessList) error {
	if list.DeleteAfterDate != "" {
		_, err := timeutil.ParseISO8601(list.DeleteAfterDate)
		if err != nil {
			return err
		}
	}
	onlyOneSpecified := onlyOneSpecified(list.AwsSecurityGroup, list.CIDRBlock, list.IPAddress)
	allSpecified := isNotEmpty(list.AwsSecurityGroup) && isNotEmpty(list.CIDRBlock) && isNotEmpty(list.IPAddress)
	if !onlyOneSpecified || allSpecified {
		return errors.New("only one of the 'awsSecurityGroup', 'cidrBlock' or 'ipAddress' is required be specified")
	}
	return nil
}

func createOrDeleteInAtlas(client mongodbatlas.Client, projectID string, operatorIPAccessLists []project.IPAccessList, log *zap.SugaredLogger) workflow.Result {
	atlasAccess, _, err := client.ProjectIPAccessList.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}
	// Making a new slice with synonyms as Atlas IP Access list to enable usage of 'Identifiable'
	atlasAccessLists := make([]atlasProjectIPAccessList, len(atlasAccess.Results))
	for i, r := range atlasAccess.Results {
		atlasAccessLists[i] = atlasProjectIPAccessList(r)
	}

	accessListsToDelete := set.Difference(atlasAccessLists, operatorIPAccessLists)

	if err := deleteIPAccessFromAtlas(client, projectID, accessListsToDelete, log); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}

	if result := createIPAccessListsInAtlas(client, projectID, operatorIPAccessLists); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

// operatorToAtlasIPAccessList converts the ipAccessList specified in the project CR to the format
// expected by the Atlas API.
func operatorToAtlasIPAccessList(ipAccessLists []project.IPAccessList) ([]*mongodbatlas.ProjectIPAccessList, workflow.Result) {
	operatorAccessLists := make([]*mongodbatlas.ProjectIPAccessList, len(ipAccessLists))
	for i, list := range ipAccessLists {
		atlasFormat, err := list.ToAtlas()
		if err != nil {
			return nil, workflow.Terminate(workflow.Internal, err.Error())
		}
		operatorAccessLists[i] = atlasFormat
	}
	return operatorAccessLists, workflow.OK()
}

func createIPAccessListsInAtlas(client mongodbatlas.Client, projectID string, ipAccessLists []project.IPAccessList) workflow.Result {
	operatorAccessLists, status := operatorToAtlasIPAccessList(ipAccessLists)
	if !status.IsOk() {
		return status
	}

	if _, _, err := client.ProjectIPAccessList.Create(context.Background(), projectID, operatorAccessLists); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}
	return workflow.OK()
}

func deleteIPAccessFromAtlas(client mongodbatlas.Client, projectID string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) error {
	for _, l := range listsToRemove {
		if _, err := client.ProjectIPAccessList.Delete(context.Background(), projectID, l.Identifier().(string)); err != nil {
			return err
		}
		log.Debugw("Removed IPAccessList from Atlas as it's not specified in current AtlasProject", "id", l.Identifier())
	}
	return nil
}

type IPAccessListStatusType string

const (
	IPAccessListActive  IPAccessListStatusType = "ACTIVE"
	IPAccessListFailed  IPAccessListStatusType = "FAILED"
	IPAccessListPending IPAccessListStatusType = "PENDING"
)

type IPAccessListStatus struct {
	Status string `json:"STATUS"`
}

// getAccessListEntry returns the identifier for the accessList. It should be exactly one of IPAddress, CIDRBlock
// or AwsSecurityGroup. This function assumes that the accessList is already validated and has only one of these
// fields populated.
func getAccessListEntry(accessList mongodbatlas.ProjectIPAccessList) string {
	if accessList.IPAddress != "" {
		return url.QueryEscape(accessList.IPAddress)
	}
	if accessList.CIDRBlock != "" {
		return url.QueryEscape(accessList.CIDRBlock)
	}
	return url.QueryEscape(accessList.AwsSecurityGroup)
}

// GetIPAccessListStatus returns the status of an individual project ip access list. The documentation can be found
// here https://docs.atlas.mongodb.com/reference/api/ip-access-list/get-one-access-list-entry-status/
func GetIPAccessListStatus(client mongodbatlas.Client, accessList mongodbatlas.ProjectIPAccessList) (IPAccessListStatus, error) {
	urlStr := fmt.Sprintf("groups/%s/accessList/%s/status", accessList.GroupID, getAccessListEntry(accessList))
	req, err := client.NewRequest(context.Background(), http.MethodGet, urlStr, nil)
	if err != nil {
		return IPAccessListStatus{}, err
	}
	ipAccessListStatus := IPAccessListStatus{}
	_, err = client.Do(context.Background(), req, &ipAccessListStatus)
	if err != nil {
		return IPAccessListStatus{}, err
	}
	return ipAccessListStatus, nil
}

func filterActiveIPAccessLists(accessLists []project.IPAccessList) ([]project.IPAccessList, []project.IPAccessList) {
	active := make([]project.IPAccessList, 0)
	expired := make([]project.IPAccessList, 0)
	for _, list := range accessLists {
		if list.DeleteAfterDate != "" {
			// We are ignoring the error as it will never happen due to validation check before
			iso8601, _ := timeutil.ParseISO8601(list.DeleteAfterDate)
			if iso8601.Before(time.Now()) {
				expired = append(expired, list)
				continue
			}
		}
		// Either 'deleteAfterDate' field is not specified or it's higher than the current time
		active = append(active, list)
	}
	return active, expired
}

func isNotEmpty(s string) bool {
	return s != ""
}

func onlyOneSpecified(values ...string) bool {
	found := false
	for _, v := range values {
		if v == "" {
			continue
		}

		if found {
			return false
		}

		found = true
	}

	return found
}
