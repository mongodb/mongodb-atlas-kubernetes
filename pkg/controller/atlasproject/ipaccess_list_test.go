package atlasproject

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestFilterActiveIPAccessLists(t *testing.T) {
	t.Run("One expired, one active", func(t *testing.T) {
		dateBefore := time.Now().UTC().Add(time.Hour * -1).Format("2006-01-02T15:04:05.999Z")
		dateAfter := time.Now().UTC().Add(time.Hour * 5).Format("2006-01-02T15:04:05.999Z")
		ipAccessExpired := project.IPAccessList{DeleteAfterDate: dateBefore}
		ipAccessActive := project.IPAccessList{DeleteAfterDate: dateAfter}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive, ipAccessExpired})
		assert.Equal(t, []project.IPAccessList{ipAccessActive}, active)
		assert.Equal(t, []project.IPAccessList{ipAccessExpired}, expired)
	})
	t.Run("Two active", func(t *testing.T) {
		dateAfter1 := time.Now().UTC().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		dateAfter2 := time.Now().UTC().Add(time.Hour * 5).Format("2006-01-02T15:04:05")
		ipAccessActive1 := project.IPAccessList{DeleteAfterDate: dateAfter1}
		ipAccessActive2 := project.IPAccessList{DeleteAfterDate: dateAfter2}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Equal(t, []project.IPAccessList{ipAccessActive2, ipAccessActive1}, active)
		assert.Empty(t, expired)
	})
}

func TestSyncIPAccessList(t *testing.T) {
	t.Run("should fail to perform deletion", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				CIDRBlock: "10.0.0.0/24",
			},
		}
		m := mockadmin.NewProjectIPAccessListApi(t)
		m.EXPECT().DeleteProjectIpAccessList(mock.Anything, mock.Anything, mock.Anything).Return(admin.DeleteProjectIpAccessListApiRequest{ApiService: m})
		m.EXPECT().DeleteProjectIpAccessListExecute(mock.Anything).Return(
			nil, nil, errors.New("failed"),
		)
		atlasClient := &admin.APIClient{ProjectIPAccessListApi: m}

		workflowCtx := &workflow.Context{
			SdkClient: atlasClient,
			Context:   context.Background(),
		}

		assert.ErrorContains(t, syncIPAccessList(workflowCtx, "projectID", current, desired), "failed")
	})

	t.Run("should fail to perform creation", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				CIDRBlock: "10.0.0.0/24",
			},
		}
		m := mockadmin.NewProjectIPAccessListApi(t)
		m.EXPECT().CreateProjectIpAccessList(mock.Anything, mock.Anything, mock.Anything).Return(admin.CreateProjectIpAccessListApiRequest{ApiService: m})
		m.EXPECT().CreateProjectIpAccessListExecute(mock.Anything).Return(
			nil, nil, errors.New("failed"),
		)
		m.EXPECT().DeleteProjectIpAccessList(mock.Anything, mock.Anything, mock.Anything).Return(admin.DeleteProjectIpAccessListApiRequest{ApiService: m})
		m.EXPECT().DeleteProjectIpAccessListExecute(mock.Anything).Return(
			nil, nil, nil,
		)
		atlasClient := &admin.APIClient{ProjectIPAccessListApi: m}
		workflowCtx := &workflow.Context{
			SdkClient: atlasClient,
			Context:   context.Background(),
		}

		assert.ErrorContains(t, syncIPAccessList(workflowCtx, "projectID", current, desired), "failed")
	})

	t.Run("should succeed when there are no changes", func(t *testing.T) {
		current := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		desired := []project.IPAccessList{
			{
				IPAddress: "10.0.0.1",
			},
		}
		workflowCtx := &workflow.Context{
			Context: context.Background(),
		}

		assert.NoError(t, syncIPAccessList(workflowCtx, "projectID", current, desired))
	})
}
