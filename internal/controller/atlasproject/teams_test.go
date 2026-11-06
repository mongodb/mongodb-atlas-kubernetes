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
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

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

func TestEnsureAssignedTeams(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		buildObjects func(t *testing.T) (*akov2.AtlasProject, *akov2.AtlasTeam)
		setupService func(s *translation.TeamsServiceMock)
		assertions   func(t *testing.T, k8sClient client.Client, teamsService *translation.TeamsServiceMock, project *akov2.AtlasProject, team *akov2.AtlasTeam, result workflow.DeprecatedResult)
	}{
		"no ownership skips sync and status changes": {
			buildObjects: func(t *testing.T) (*akov2.AtlasProject, *akov2.AtlasTeam) {
				project := &akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "projectName",
						Namespace: "testNS",
						Annotations: map[string]string{
							customresource.AnnotationLastAppliedConfiguration: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{Name: "projectName"}),
						},
					},
					Spec: akov2.AtlasProjectSpec{
						Name:  "projectName",
						Teams: nil,
					},
					Status: status.AtlasProjectStatus{
						ID: "projectID",
						Teams: []status.ProjectTeamStatus{{
							ID:      "teamID_1",
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
						}},
					},
				}

				team := &akov2.AtlasTeam{
					ObjectMeta: metav1.ObjectMeta{Name: "teamName_1", Namespace: "testNS"},
					Status: status.TeamStatus{
						ID: "teamID_1",
						Projects: []status.TeamProject{{
							ID:   "projectID",
							Name: "projectName",
						}},
					},
				}

				return project, team
			},
			assertions: func(t *testing.T, k8sClient client.Client, teamsService *translation.TeamsServiceMock, project *akov2.AtlasProject, team *akov2.AtlasTeam, result workflow.DeprecatedResult) {
				assert.True(t, result.IsOk())
				teamsService.AssertNotCalled(t, "ListProjectTeams", mock.Anything, mock.Anything)
				teamsService.AssertNotCalled(t, "Unassign", mock.Anything, mock.Anything, mock.Anything)

				gotProject := &akov2.AtlasProject{}
				assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(project), gotProject))
				assert.Equal(t, []status.ProjectTeamStatus{{ID: "teamID_1", TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"}}}, gotProject.Status.Teams)

				gotTeam := &akov2.AtlasTeam{}
				assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(team), gotTeam))
				assert.Equal(t, []status.TeamProject{{ID: "projectID", Name: "projectName"}}, gotTeam.Status.Projects)
			},
		},
		"no ownership malformed last applied terminates": {
			buildObjects: func(_ *testing.T) (*akov2.AtlasProject, *akov2.AtlasTeam) {
				project := &akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "projectName",
						Namespace: "testNS",
						Annotations: map[string]string{
							customresource.AnnotationLastAppliedConfiguration: "{not-valid-json}",
						},
					},
					Spec: akov2.AtlasProjectSpec{Name: "projectName"},
					Status: status.AtlasProjectStatus{
						ID: "projectID",
					},
				}

				return project, nil
			},
			assertions: func(t *testing.T, _ client.Client, teamsService *translation.TeamsServiceMock, _ *akov2.AtlasProject, _ *akov2.AtlasTeam, result workflow.DeprecatedResult) {
				assert.False(t, result.IsOk())
				assert.ErrorContains(t, result.GetError(), "failed to get last applied configuration")
				teamsService.AssertNotCalled(t, "ListProjectTeams", mock.Anything, mock.Anything)
				teamsService.AssertNotCalled(t, "Unassign", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		"with previously applied teams still removes": {
			buildObjects: func(t *testing.T) (*akov2.AtlasProject, *akov2.AtlasTeam) {
				project := &akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "projectName",
						Namespace: "testNS",
						Annotations: map[string]string{
							customresource.AnnotationLastAppliedConfiguration: mustMarshalProjectSpec(t, akov2.AtlasProjectSpec{
								Name: "projectName",
								Teams: []akov2.Team{{
									TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
									Roles:   []akov2.TeamRole{"GROUP_OWNER"},
								}},
							}),
						},
					},
					Spec: akov2.AtlasProjectSpec{Name: "projectName", Teams: nil},
					Status: status.AtlasProjectStatus{
						ID: "projectID",
						Teams: []status.ProjectTeamStatus{{
							ID:      "teamID_1",
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_1", Namespace: "testNS"},
						}},
					},
				}

				team := &akov2.AtlasTeam{
					ObjectMeta: metav1.ObjectMeta{Name: "teamName_1", Namespace: "testNS"},
					Status: status.TeamStatus{
						ID: "teamID_1",
						Projects: []status.TeamProject{{
							ID:   "projectID",
							Name: "projectName",
						}},
					},
				}

				return project, team
			},
			setupService: func(s *translation.TeamsServiceMock) {
				s.EXPECT().ListProjectTeams(ctx, "projectID").Return([]teams.AssignedTeam{{
					TeamID:   "teamID_1",
					TeamName: "teamName_1",
					Roles:    []string{"GROUP_OWNER"},
				}}, nil)
				s.EXPECT().Unassign(ctx, "projectID", "teamID_1").Return(nil)
			},
			assertions: func(t *testing.T, k8sClient client.Client, teamsService *translation.TeamsServiceMock, project *akov2.AtlasProject, team *akov2.AtlasTeam, result workflow.DeprecatedResult) {
				assert.True(t, result.IsOk())
				teamsService.AssertNotCalled(t, "Assign", mock.Anything, mock.Anything, mock.Anything)

				gotProject := &akov2.AtlasProject{}
				assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(project), gotProject))
				assert.Empty(t, gotProject.Status.Teams)

				gotTeam := &akov2.AtlasTeam{}
				assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(team), gotTeam))
				assert.Empty(t, gotTeam.Status.Projects)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			project, team := tc.buildObjects(t)

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))

			objects := []client.Object{project}
			statusObjects := []client.Object{project}
			if team != nil {
				objects = append(objects, team)
				statusObjects = append(statusObjects, team)
			}

			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(objects...).
				WithStatusSubresource(statusObjects...).
				Build()

			teamsService := translation.NewTeamsServiceMock(t)
			if tc.setupService != nil {
				tc.setupService(teamsService)
			}

			logger := zaptest.NewLogger(t).Sugar()
			workflowCtx := &workflow.Context{Context: ctx, Log: logger}
			r := &AtlasProjectReconciler{
				Client:        k8sClient,
				EventRecorder: record.NewFakeRecorder(10),
				Log:           logger,
			}

			result := r.ensureAssignedTeams(workflowCtx, teamsService, project)
			tc.assertions(t, k8sClient, teamsService, project, team, result)
		})
	}
}

func mustMarshalProjectSpec(t *testing.T, spec akov2.AtlasProjectSpec) string {
	t.Helper()
	raw, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal project spec: %v", err)
	}

	return string(raw)
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
