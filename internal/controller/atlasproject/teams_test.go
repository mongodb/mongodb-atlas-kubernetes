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

package atlasproject

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
)

func TestSyncAssignedTeams(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		teamsToAssign   map[string]*akov2.Team
		teamServiceMock func() teams.TeamsService
		expectedErr     error
	}{
		"should error to list teams": {
			teamServiceMock: func() teams.TeamsService {
				s := translation.NewTeamsServiceMock(t)
				s.EXPECT().
					ListProjectTeams(ctx, "projectID").
					Return(nil, errors.New("error to list teams"))

				return s
			},
			expectedErr: errors.New("error to list teams"),
		},
		"should error to unassign teams": {
			teamServiceMock: func() teams.TeamsService {
				s := translation.NewTeamsServiceMock(t)
				s.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_1",
						TeamName: "teamName_1",
					},
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_2",
						TeamName: "teamName_2",
					},
					{
						Roles:    []string{"GROUP_READ_ONLY"},
						TeamID:   "teamID_3",
						TeamName: "teamName_3",
					},
				}, nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_1").Return(errors.New("error to unassign team 1"))
				s.EXPECT().Unassign(ctx, "projectID", "teamID_2").Return(errors.New("error to unassign team 2"))
				s.EXPECT().Unassign(ctx, "projectID", "teamID_3").Return(errors.New("error to unassign team 3"))

				return s
			},
			expectedErr: errors.Join(
				errors.Join(
					errors.Join(
						nil,
						errors.New("error to unassign team 1"),
					),
					errors.New("error to unassign team 2"),
				),
				errors.New("error to unassign team 3"),
			),
		},
		"should error to assign teams": {
			teamsToAssign: map[string]*akov2.Team{
				"teamID_1": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_1",
					},
					Roles: []akov2.TeamRole{"GROUP_OWNER"},
				},
				"teamID_2": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_2",
					},
					Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
				},
			},
			teamServiceMock: func() teams.TeamsService {
				s := translation.NewTeamsServiceMock(t)
				s.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_1",
						TeamName: "teamName_1",
					},
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_2",
						TeamName: "teamName_2",
					},
					{
						Roles:    []string{"GROUP_READ_ONLY"},
						TeamID:   "teamID_3",
						TeamName: "teamName_3",
					},
				}, nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_2").Return(nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_3").Return(nil)

				s.EXPECT().Assign(
					ctx,
					&[]teams.AssignedTeam{
						{
							Roles:  []string{"GROUP_READ_ONLY"},
							TeamID: "teamID_2",
						},
					},
					"projectID",
				).Return(errors.New("error to assign team 2"))

				return s
			},
			expectedErr: errors.Join(
				nil,
				errors.New("error to assign team 2"),
			),
		},
		"should sync teams assigned": {
			teamsToAssign: map[string]*akov2.Team{
				"teamID_1": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_1",
					},
					Roles: []akov2.TeamRole{"GROUP_OWNER"},
				},
				"teamID_2": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_2",
					},
					Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
				},
			},
			teamServiceMock: func() teams.TeamsService {
				s := translation.NewTeamsServiceMock(t)
				s.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_1",
						TeamName: "teamName_1",
					},
					{
						Roles:    []string{"GROUP_OWNER"},
						TeamID:   "teamID_2",
						TeamName: "teamName_2",
					},
					{
						Roles:    []string{"GROUP_READ_ONLY"},
						TeamID:   "teamID_3",
						TeamName: "teamName_3",
					},
				}, nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_2").Return(nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_3").Return(nil)
				s.EXPECT().Assign(
					ctx,
					&[]teams.AssignedTeam{
						{
							Roles:  []string{"GROUP_READ_ONLY"},
							TeamID: "teamID_2",
						},
					},
					"projectID",
				).
					Return(nil)

				return s
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name: "projectName",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "projectName",
					Teams: []akov2.Team{
						{
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_1"},
							Roles:   []akov2.TeamRole{"GROUP_OWNER"},
						},
						{
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_2"},
							Roles:   []akov2.TeamRole{"GROUP_READ_ONLY"},
						},
					},
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
					Teams: []status.ProjectTeamStatus{
						{
							ID: "teamID_1",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_1",
							},
						},
						{
							ID: "teamID_2",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_2",
							},
						},
						{
							ID: "teamID_3",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_3",
							},
						},
					},
				},
			}
			team1 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_1",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_1",
				},
				Status: status.TeamStatus{
					ID: "teamID_1",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}
			team2 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_2",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_2",
				},
				Status: status.TeamStatus{
					ID: "teamID_2",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}
			team3 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_3",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_3",
				},
				Status: status.TeamStatus{
					ID: "teamID_3",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, team1, team2, team3).
				WithStatusSubresource(project, team1, team2, team3).
				Build()

			atlasClient := &admin.APIClient{}
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Log: logger,
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312: atlasClient,
				},
				Context: context.Background(),
			}
			r := &AtlasProjectReconciler{
				Client:        k8sClient,
				EventRecorder: record.NewFakeRecorder(10),
				Log:           logger,
			}
			err := r.syncAssignedTeams(ctx, tt.teamServiceMock(), "projectID", project, tt.teamsToAssign)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestEnsureAssignedTeamsRemovesTeams(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
	}{
		{
			name: "missing last-applied annotation",
		},
		{
			name: "empty last-applied teams",
			annotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{Name: "projectName"}),
			},
		},
		{
			name: "malformed last-applied JSON",
			annotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{not-valid-json}",
			},
		},
		{
			name: "previously applied teams",
			annotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{
					Name: "projectName",
					Teams: []akov2.Team{{
						TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
						Roles:   []akov2.TeamRole{"GROUP_OWNER"},
					}},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newEnsureAssignedTeamsFixture(t, tt.annotations, nil)
			ctx := fixture.workflowCtx.Context
			fixture.teamsService.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{{
				TeamID:   "teamID_1",
				TeamName: "teamName_1",
				Roles:    []string{"GROUP_OWNER"},
			}}, nil)
			fixture.teamsService.EXPECT().Unassign(ctx, "projectID", "teamID_1").Return(nil)

			result := fixture.reconciler.ensureAssignedTeams(fixture.workflowCtx, fixture.teamsService, fixture.project)

			assert.True(t, result.IsOk())
			fixture.teamsService.AssertExpectations(t)
			fixture.teamsService.AssertNotCalled(t, "Assign", mock.Anything, mock.Anything, mock.Anything)

			gotProject := &akov2.AtlasProject{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.project), gotProject))
			assert.Empty(t, gotProject.Status.Teams)

			gotTeam := &akov2.AtlasTeam{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.team), gotTeam))
			assert.Empty(t, gotTeam.Status.Projects)

			_, found := fixture.workflowCtx.GetCondition(api.ProjectTeamsReadyType)
			assert.False(t, found, "ProjectTeamsReady condition should be removed when the desired team list is empty")
		})
	}
}

func TestEnsureAssignedTeamsOwnership(t *testing.T) {
	tests := []struct {
		name               string
		unmanagedPolicy    bool
		lastAppliedPresent bool
		lastAppliedTeams   bool
		wantManaged        bool
	}{
		{name: "managed without last-applied", wantManaged: true},
		{name: "managed with empty last-applied teams", lastAppliedPresent: true, wantManaged: true},
		{name: "managed with last-applied teams", lastAppliedPresent: true, lastAppliedTeams: true, wantManaged: true},
		{name: "unmanaged without last-applied", unmanagedPolicy: true},
		{name: "unmanaged with empty last-applied teams", unmanagedPolicy: true, lastAppliedPresent: true},
		{name: "unmanaged with last-applied teams", unmanagedPolicy: true, lastAppliedPresent: true, lastAppliedTeams: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			annotations := map[string]string{}
			if tt.unmanagedPolicy {
				annotations[customresource.ProjectTeamsPolicyAnnotation] = customresource.ProjectTeamsPolicyUnmanaged
			}
			if tt.lastAppliedPresent {
				lastApplied := akov2.AtlasProjectSpec{Name: "projectName"}
				if tt.lastAppliedTeams {
					lastApplied.Teams = []akov2.Team{{TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}}}
				}
				annotations[customresource.AnnotationLastAppliedConfiguration] = mustMarshalProjectSpec(t, lastApplied)
			}
			initialLastApplied := annotations[customresource.AnnotationLastAppliedConfiguration]
			fixture := newEnsureAssignedTeamsFixture(t, annotations, nil)
			ctx := fixture.workflowCtx.Context
			originalTeamsStatus := append([]status.ProjectTeamStatus(nil), fixture.project.Status.Teams...)
			originalProjectsStatus := append([]status.TeamProject(nil), fixture.team.Status.Projects...)

			if tt.wantManaged {
				fixture.teamsService.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{{
					TeamID:   "teamID_1",
					TeamName: "teamName_1",
					Roles:    []string{"GROUP_OWNER"},
				}}, nil)
				fixture.teamsService.EXPECT().Unassign(ctx, "projectID", "teamID_1").Return(nil)
			}

			result := fixture.reconciler.ensureAssignedTeams(fixture.workflowCtx, fixture.teamsService, fixture.project)

			assert.True(t, result.IsOk())
			fixture.teamsService.AssertNotCalled(t, "Assign", mock.Anything, mock.Anything, mock.Anything)
			if !tt.wantManaged {
				fixture.teamsService.AssertNotCalled(t, "ListProjectTeams", mock.Anything, mock.Anything)
				fixture.teamsService.AssertNotCalled(t, "Unassign", mock.Anything, mock.Anything, mock.Anything)
			}
			fixture.teamsService.AssertExpectations(t)

			gotProject := &akov2.AtlasProject{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.project), gotProject))
			gotTeam := &akov2.AtlasTeam{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.team), gotTeam))
			if tt.wantManaged {
				assert.Empty(t, gotProject.Status.Teams)
				assert.Empty(t, gotTeam.Status.Projects)
			} else {
				assert.Equal(t, originalTeamsStatus, gotProject.Status.Teams)
				assert.Equal(t, originalProjectsStatus, gotTeam.Status.Projects)
			}

			if !tt.lastAppliedPresent {
				_, found := gotProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
				assert.False(t, found, "last-applied annotation should remain absent")
			} else {
				lastApplied, err := customresource.ParseLastConfigApplied[akov2.AtlasProjectSpec](gotProject)
				require.NoError(t, err)
				require.NotNil(t, lastApplied)
				assert.Equal(t, "projectName", lastApplied.Name, "non-team last-applied fields should be preserved")
				if tt.unmanagedPolicy {
					assert.Empty(t, lastApplied.Teams)
				} else {
					assert.Equal(t, initialLastApplied, gotProject.Annotations[customresource.AnnotationLastAppliedConfiguration], "managed last-applied annotation should remain unchanged")
				}
			}

			_, found := fixture.workflowCtx.GetCondition(api.ProjectTeamsReadyType)
			assert.False(t, found, "ProjectTeamsReady condition should be removed when the desired team list is empty")
		})
	}
}

func TestEnsureAssignedTeamsUnmanagedRejectsNonemptySpec(t *testing.T) {
	annotations := map[string]string{
		customresource.ProjectTeamsPolicyAnnotation:       customresource.ProjectTeamsPolicyUnmanaged,
		customresource.AnnotationLastAppliedConfiguration: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{Name: "projectName"}),
	}
	desiredTeams := []akov2.Team{{TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}}}
	fixture := newEnsureAssignedTeamsFixture(t, annotations, desiredTeams)
	ctx := fixture.workflowCtx.Context
	originalTeamsStatus := append([]status.ProjectTeamStatus(nil), fixture.project.Status.Teams...)
	originalProjectsStatus := append([]status.TeamProject(nil), fixture.team.Status.Projects...)
	originalLastApplied := fixture.project.Annotations[customresource.AnnotationLastAppliedConfiguration]

	result := fixture.reconciler.ensureAssignedTeams(fixture.workflowCtx, fixture.teamsService, fixture.project)

	assert.False(t, result.IsOk())
	assert.ErrorContains(t, result.GetError(), "spec.teams must be empty")
	assert.ErrorContains(t, result.GetError(), "unmanaged teams policy")
	condition, found := fixture.workflowCtx.GetCondition(api.ProjectTeamsReadyType)
	require.True(t, found)
	assert.Equal(t, string(metav1.ConditionFalse), string(condition.Status))
	fixture.teamsService.AssertNotCalled(t, "ListProjectTeams", mock.Anything, mock.Anything)
	fixture.teamsService.AssertNotCalled(t, "Assign", mock.Anything, mock.Anything, mock.Anything)
	fixture.teamsService.AssertNotCalled(t, "Unassign", mock.Anything, mock.Anything, mock.Anything)
	fixture.teamsService.AssertExpectations(t)

	gotProject := &akov2.AtlasProject{}
	require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.project), gotProject))
	assert.Equal(t, originalTeamsStatus, gotProject.Status.Teams)
	assert.Equal(t, originalLastApplied, gotProject.Annotations[customresource.AnnotationLastAppliedConfiguration])

	gotTeam := &akov2.AtlasTeam{}
	require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.team), gotTeam))
	assert.Equal(t, originalProjectsStatus, gotTeam.Status.Projects)
}

func TestEnsureAssignedTeamsUnmanagedErrors(t *testing.T) {
	patchErr := errors.New("sentinel patch error")
	tests := []struct {
		name               string
		lastApplied        string
		patchErr           error
		expectedErrorParts []string
	}{
		{
			name:               "malformed last-applied JSON",
			lastApplied:        "{not-valid-json}",
			expectedErrorParts: []string{"failed to parse last applied configuration"},
		},
		{
			name: "failed to patch last-applied configuration",
			lastApplied: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{
				Name:  "projectName",
				Teams: []akov2.Team{{TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}}},
			}),
			patchErr:           patchErr,
			expectedErrorParts: []string{"failed to patch last applied configuration", patchErr.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			annotations := map[string]string{
				customresource.ProjectTeamsPolicyAnnotation:       customresource.ProjectTeamsPolicyUnmanaged,
				customresource.AnnotationLastAppliedConfiguration: tt.lastApplied,
			}
			fixture := newEnsureAssignedTeamsFixture(t, annotations, nil)
			ctx := fixture.workflowCtx.Context
			originalTeamsStatus := append([]status.ProjectTeamStatus(nil), fixture.project.Status.Teams...)
			originalProjectsStatus := append([]status.TeamProject(nil), fixture.team.Status.Projects...)
			originalLastApplied := fixture.project.Annotations[customresource.AnnotationLastAppliedConfiguration]
			if tt.patchErr != nil {
				fixture.reconciler.Client = &patchErrorClient{Client: fixture.k8sClient, err: tt.patchErr}
			}

			result := fixture.reconciler.ensureAssignedTeams(fixture.workflowCtx, fixture.teamsService, fixture.project)

			assert.False(t, result.IsOk())
			for _, part := range tt.expectedErrorParts {
				assert.ErrorContains(t, result.GetError(), part)
			}
			condition, found := fixture.workflowCtx.GetCondition(api.ProjectTeamsReadyType)
			require.True(t, found)
			assert.Equal(t, string(metav1.ConditionFalse), string(condition.Status))
			fixture.teamsService.AssertNotCalled(t, "ListProjectTeams", mock.Anything, mock.Anything)
			fixture.teamsService.AssertNotCalled(t, "Assign", mock.Anything, mock.Anything, mock.Anything)
			fixture.teamsService.AssertNotCalled(t, "Unassign", mock.Anything, mock.Anything, mock.Anything)
			fixture.teamsService.AssertExpectations(t)

			gotProject := &akov2.AtlasProject{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.project), gotProject))
			assert.Equal(t, originalTeamsStatus, gotProject.Status.Teams)
			assert.Equal(t, originalLastApplied, gotProject.Annotations[customresource.AnnotationLastAppliedConfiguration])

			gotTeam := &akov2.AtlasTeam{}
			require.NoError(t, fixture.k8sClient.Get(ctx, client.ObjectKeyFromObject(fixture.team), gotTeam))
			assert.Equal(t, originalProjectsStatus, gotTeam.Status.Projects)
		})
	}
}

type patchErrorClient struct {
	client.Client
	err error
}

func (c *patchErrorClient) Patch(context.Context, client.Object, client.Patch, ...client.PatchOption) error {
	return c.err
}

func mustMarshalProjectSpec(t *testing.T, spec akov2.AtlasProjectSpec) string {
	t.Helper()
	raw, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal project spec: %v", err)
	}

	return string(raw)
}

type ensureAssignedTeamsFixture struct {
	k8sClient    client.Client
	project      *akov2.AtlasProject
	team         *akov2.AtlasTeam
	teamsService *translation.TeamsServiceMock
	workflowCtx  *workflow.Context
	reconciler   *AtlasProjectReconciler
}

func newEnsureAssignedTeamsFixture(t *testing.T, annotations map[string]string, desiredTeams []akov2.Team) ensureAssignedTeamsFixture {
	t.Helper()

	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "projectName",
			Namespace:   "testNS",
			Annotations: annotations,
		},
		Spec: akov2.AtlasProjectSpec{
			Name:  "projectName",
			Teams: desiredTeams,
		},
		Status: status.AtlasProjectStatus{
			ID:    "projectID",
			Teams: []status.ProjectTeamStatus{{ID: "teamID_1", TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}}},
		},
	}
	team := &akov2.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{Name: "teamName_1", Namespace: "testNS"},
		Status: status.TeamStatus{
			ID:       "teamID_1",
			Projects: []status.TeamProject{{ID: "projectID", Name: "projectName"}},
		},
	}

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(project, team).
		WithStatusSubresource(project, team).
		Build()

	logger := zaptest.NewLogger(t).Sugar()
	workflowCtx := &workflow.Context{Context: context.Background(), Log: logger}
	workflowCtx.SetConditionTrue(api.ProjectTeamsReadyType)

	return ensureAssignedTeamsFixture{
		k8sClient:    k8sClient,
		project:      project,
		team:         team,
		teamsService: translation.NewTeamsServiceMock(t),
		workflowCtx:  workflowCtx,
		reconciler: &AtlasProjectReconciler{
			Client:        k8sClient,
			EventRecorder: record.NewFakeRecorder(10),
			Log:           logger,
		},
	}
}

func TestUpdateTeamState(t *testing.T) {
	tests := map[string]struct {
		team                     *akov2.AtlasTeam
		isRemoval                bool
		expectedAssignedProjects []status.TeamProject
	}{
		"should add project to team status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
			},
		},
		"should not duplicate projects already listed on status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
			},
		},
		"should remove project from team status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
				},
			},
			isRemoval:                true,
			expectedAssignedProjects: nil,
		},
		"should deduplicate pre-existing entries for the current project": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{Name: "testTeam", Namespace: "testNS"},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{ID: "projectID", Name: "projectName"},
						{ID: "projectID", Name: "projectName"},
						{ID: "projectID", Name: "projectName"},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{ID: "projectID", Name: "projectName"},
			},
		},
		"should deduplicate pre-existing entries for other projects": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{Name: "testTeam", Namespace: "testNS"},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{ID: "otherProjectID", Name: "otherProjectName"},
						{ID: "otherProjectID", Name: "otherProjectName"},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{ID: "projectID", Name: "projectName"},
				{ID: "otherProjectID", Name: "otherProjectName"},
			},
		},
		"should not modify status of other assigned projects": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{
							ID:   "existingProjectID",
							Name: "existingProjectName",
						},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
				{
					ID:   "existingProjectID",
					Name: "existingProjectName",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-secret",
				},
				Data: map[string][]byte{
					"orgId":         []byte("0987654321"),
					"publicApiKey":  []byte("api-pub-key"),
					"privateApiKey": []byte("api-priv-key"),
				},
				Type: "Opaque",
			}
			project := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name: "projectName",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name: "my-secret",
					},
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(secret, project, tt.team).
				WithStatusSubresource(tt.team).
				Build()

			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasProjectReconciler{
				Client:        k8sClient,
				EventRecorder: record.NewFakeRecorder(1),
				Log:           logger,
			}

			err := r.updateTeamState(ctx, project, &common.ResourceRefNamespaced{Name: tt.team.Name, Namespace: tt.team.Namespace}, tt.isRemoval)
			assert.NoError(t, err)

			actualTeam := &akov2.AtlasTeam{}
			assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(tt.team), actualTeam))
			assert.Equal(t, tt.expectedAssignedProjects, actualTeam.Status.Projects)
		})
	}
}

// TestSyncAssignedTeamsSteadyStateDeduplication is the regression test for the bug
// reported in https://github.com/mongodb/mongodb-atlas-kubernetes/issues/3359.
//
// Root cause: a prior operator bug (fixed in 33196ea9) caused updateTeamState to
// append rather than replace project entries, so each reconcile cycle that
// triggered a team assignment would grow status.projects by one entry. After
// upgrading to the version with that fix, the duplicates never go away because
// syncAssignedTeams skips calling updateTeamState for teams already in steady
// state (assigned with matching roles). Crucially, in steady state the count does
// not keep growing either — it is simply frozen.
//
// This test simulates two consecutive reconcile cycles in steady state and verifies
// that:
//   - the duplicate count does not increase between cycles (frozen without the fix),
//   - the list converges to exactly one entry on the first cycle (with the fix), and
//   - a second cycle leaves the list unchanged (fix is idempotent).
//
// Without the fix both assertions on Projects length fail because updateTeamState
// is never called for steady-state teams and the three duplicates remain.
func TestSyncAssignedTeamsSteadyStateDeduplication(t *testing.T) {
	ctx := context.Background()

	// Team carries three copies of the same project — the residue left by the
	// prior accumulation bug after an operator upgrade.
	team := &akov2.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{Name: "teamName_1", Namespace: "testNS"},
		Spec:       akov2.TeamSpec{Name: "teamName_1"},
		Status: status.TeamStatus{
			ID: "teamID_1",
			Projects: []status.TeamProject{
				{ID: "projectID", Name: "projectName"},
				{ID: "projectID", Name: "projectName"},
				{ID: "projectID", Name: "projectName"},
			},
		},
	}
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: "projectName", Namespace: "testNS"},
		Spec: akov2.AtlasProjectSpec{
			Name: "projectName",
			Teams: []akov2.Team{
				{
					TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
					Roles:   []akov2.TeamRole{"GROUP_OWNER"},
				},
			},
		},
		Status: status.AtlasProjectStatus{
			ID: "projectID",
			Teams: []status.ProjectTeamStatus{
				{ID: "teamID_1", TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}},
			},
		},
	}

	testScheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(testScheme))
	assert.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(project, team).
		WithStatusSubresource(project, team).
		Build()

	// Atlas already has the team assigned with the correct role — steady state.
	// The mock is registered twice: once per reconcile cycle.
	atlasAssignedTeams := []teams.AssignedTeam{
		{TeamID: "teamID_1", TeamName: "teamName_1", Roles: []string{"GROUP_OWNER"}},
	}
	teamsService := translation.NewTeamsServiceMock(t)
	teamsService.EXPECT().
		ListProjectTeams(ctx, "projectID").
		Return(atlasAssignedTeams, nil).
		Times(2)

	logger := zaptest.NewLogger(t).Sugar()
	r := &AtlasProjectReconciler{
		Client:        k8sClient,
		EventRecorder: record.NewFakeRecorder(10),
		Log:           logger,
	}
	workflowCtx := &workflow.Context{
		Log:          logger,
		Context:      ctx,
		SdkClientSet: &atlas.ClientSet{SdkClient20250312: &admin.APIClient{}},
	}
	newTeamsToAssign := func() map[string]*akov2.Team {
		return map[string]*akov2.Team{
			"teamID_1": {
				TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
				Roles:   []akov2.TeamRole{"GROUP_OWNER"},
			},
		}
	}

	// Cycle 1: duplicates must be cleaned up.
	assert.NoError(t, r.syncAssignedTeams(workflowCtx, teamsService, "projectID", project, newTeamsToAssign()))
	got := &akov2.AtlasTeam{}
	assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(team), got))
	assert.Equal(t,
		[]status.TeamProject{{ID: "projectID", Name: "projectName"}},
		got.Status.Projects,
		"cycle 1: steady-state reconcile must deduplicate status.projects",
	)

	// Cycle 2: a second reconcile must not re-introduce duplicates (fix is idempotent).
	assert.NoError(t, r.syncAssignedTeams(workflowCtx, teamsService, "projectID", project, newTeamsToAssign()))
	assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(team), got))
	assert.Equal(t,
		[]status.TeamProject{{ID: "projectID", Name: "projectName"}},
		got.Status.Projects,
		"cycle 2: a second steady-state reconcile must leave status.projects unchanged",
	)
}
