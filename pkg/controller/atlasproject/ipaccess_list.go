package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/timeutil"
)

type atlasProjectIPAccessList mongodbatlas.ProjectIPAccessList

func (i atlasProjectIPAccessList) Identifier() interface{} {
	return i.CIDRBlock + i.AwsSecurityGroup + i.IPAddress
}

func (r *AtlasProjectReconciler) ensureIPAccessList(ctx *workflow.Context, connection atlas.Connection, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	if err := validateIPAccessLists(project.Spec.ProjectIPAccessList); err != nil {
		return workflow.Terminate(workflow.ProjectIPAccessInvalid, err.Error())
	}
	active, expired, result := filterActiveIPAccessLists(project)
	if !result.IsOk() {
		return result
	}
	client, err := atlas.Client(r.AtlasDomain, connection, ctx.Log)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	if result := createOrDeleteInAtlas(client, projectID, active, ctx.Log); !result.IsOk() {
		return result
	}
	fmt.Println(active, expired, ctx, connection)
	return workflow.OK()
}

func validateIPAccessLists(ipAccessList []mdbv1.ProjectIPAccessList) error {
	for _, list := range ipAccessList {
		if list.DeleteAfterDate != "" {
			_, err := timeutil.ParseISO8601(list.DeleteAfterDate)
			if err != nil {
				return err
			}
		}
		// Go doesn't support XOR, but uses '!=' instead: https://stackoverflow.com/a/23025720/614239
		onlyOneSpecified := isNotEmpty(list.AwsSecurityGroup) != isNotEmpty(list.CIDRBlock) != isNotEmpty(list.IPAddress)
		if !onlyOneSpecified {
			return errors.New("only one of the 'awsSecurityGroup', 'cidrBlock' or 'ipAddress' is required be specified")
		}
	}
	return nil
}

func createOrDeleteInAtlas(client *mongodbatlas.Client, projectID string, ipAccessLists []mdbv1.ProjectIPAccessList, log *zap.SugaredLogger) workflow.Result {
	atlasAccessLists, _, err := client.ProjectIPAccessList.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}
	atlasAccess := make([]atlasProjectIPAccessList, len(atlasAccessLists.Results))
	for i, r := range atlasAccessLists.Results {
		atlasAccess[i] = atlasProjectIPAccessList(r)
	}

	operatorAccessLists := make([]*mongodbatlas.ProjectIPAccessList, len(ipAccessLists))
	for i, list := range ipAccessLists {
		atlasFormat, err := list.ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
		operatorAccessLists[i] = &atlasFormat
	}

	difference := set.Difference(atlasAccess, operatorAccessLists)

	if err := deleteIPAccessFromAtlas(client, projectID, difference, log); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}

	if _, _, err := client.ProjectIPAccessList.Create(context.Background(), projectID, operatorAccessLists); err != nil {
		return workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, err.Error())
	}
	return workflow.OK()
}

func deleteIPAccessFromAtlas(client *mongodbatlas.Client, projectID string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) error {
	for _, l := range listsToRemove {
		if _, err := client.ProjectIPAccessList.Delete(context.Background(), projectID, l.Identifier().(string)); err != nil {
			return err
		}
		log.Debugw("Removed IPAccessList from Atlas as it's not specified in current AtlasProject", "id", l.Identifier())
	}
	return nil
}

func filterActiveIPAccessLists(project *mdbv1.AtlasProject) ([]mdbv1.ProjectIPAccessList, []mdbv1.ProjectIPAccessList, workflow.Result) {
	active := make([]mdbv1.ProjectIPAccessList, 0)
	expired := make([]mdbv1.ProjectIPAccessList, 0)
	for _, list := range project.Spec.ProjectIPAccessList {
		if list.DeleteAfterDate != "" {
			iso8601, err := timeutil.ParseISO8601(list.DeleteAfterDate)
			if err != nil {
				// Bad formatting done by user
				return active, expired, workflow.Terminate(workflow.ProjectIPAccessInvalid, err.Error())
			}
			if iso8601.Before(time.Now()) {
				expired = append(expired, list)
				continue
			}
		}
		// Either 'deleteAfterDate' field is not specified or it's higher than the current time
		active = append(active, list)
	}
	return active, expired, workflow.OK()
}

func isNotEmpty(s string) bool {
	return s != ""
}
