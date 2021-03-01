package atlasproject

import (
	"context"
	"errors"
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

func (r *AtlasProjectReconciler) ensureIPAccessList(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
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

	difference := set.Difference(atlasAccessLists, operatorIPAccessLists)

	if err := deleteIPAccessFromAtlas(client, projectID, difference, log); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}

	if result := createIPAccessListsInAtlas(client, projectID, operatorIPAccessLists); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

func createIPAccessListsInAtlas(client mongodbatlas.Client, projectID string, ipAccessLists []project.IPAccessList) workflow.Result {
	operatorAccessLists := make([]*mongodbatlas.ProjectIPAccessList, len(ipAccessLists))
	for i, list := range ipAccessLists {
		atlasFormat, err := list.ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		operatorAccessLists[i] = atlasFormat
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
