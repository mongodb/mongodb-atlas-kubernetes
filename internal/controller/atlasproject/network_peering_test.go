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
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestNetworkPeeringsNonGreedyBehaviour(t *testing.T) {
	for _, tc := range []struct {
		title                   string
		lastAppliedNetworkPeers []string
		specNetworkPeers        []string
		atlasNetworkPeers       []string
		wantRemoved             []string
	}{
		{
			title:                   "no last applied no removal in Atlas",
			lastAppliedNetworkPeers: []string{},
			specNetworkPeers:        []string{},
			atlasNetworkPeers:       []string{"np1", "np2"},
			wantRemoved:             []string{},
		},
		{
			title:                   "removed from last applied removes from Atlas",
			lastAppliedNetworkPeers: []string{"np1", "np2"},
			specNetworkPeers:        []string{"np1"},
			atlasNetworkPeers:       []string{"np1", "np2"},
			wantRemoved:             []string{"np2"},
		},
		{
			title:                   "removed all from last applied removes all from Atlas",
			lastAppliedNetworkPeers: []string{"np1", "np2"},
			specNetworkPeers:        []string{},
			atlasNetworkPeers:       []string{"np1", "np2"},
			wantRemoved:             []string{"np1", "np2"},
		},
		{
			title:                   "not in last applied not removed from Atlas",
			lastAppliedNetworkPeers: []string{"np1"},
			specNetworkPeers:        []string{"np1"},
			atlasNetworkPeers:       []string{"np1", "np2"},
			wantRemoved:             []string{},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			prj := newNetworkPeeringTestProject(tc.specNetworkPeers)
			lastPrj := newNetworkPeeringTestProject(tc.lastAppliedNetworkPeers)
			prj.Annotations[customresource.AnnotationLastAppliedConfiguration] = jsonize(t, lastPrj.Spec)

			peeringAPI := mockadmin.NewNetworkPeeringApi(t)
			peeringAPI.EXPECT().ListGroupPeersWithParams(mock.Anything, mock.Anything).
				Return(admin.ListGroupPeersApiRequest{ApiService: peeringAPI}).Once()
			peeringAPI.EXPECT().ListGroupPeersExecute(
				mock.Anything).Return(
				synthesizeAtlasNetworkPeerings(tc.atlasNetworkPeers), nil, nil,
			).Once()
			peeringAPI.EXPECT().ListGroupPeersWithParams(mock.Anything, mock.Anything).
				Return(admin.ListGroupPeersApiRequest{ApiService: peeringAPI}).Twice()
			peeringAPI.EXPECT().ListGroupPeersExecute(
				mock.Anything).Return(
				nil, nil, nil,
			).Twice()

			removals := len(tc.wantRemoved)
			if removals > 0 {
				peeringAPI.EXPECT().DeleteGroupPeer(
					mock.Anything, mock.Anything, mock.Anything,
				).Return(admin.DeleteGroupPeerApiRequest{ApiService: peeringAPI}).Times(removals)
				peeringAPI.EXPECT().DeleteGroupPeerExecute(
					mock.Anything).Return(
					nil, nil, nil,
				).Times(removals)
			}

			peeringAPI.EXPECT().ListGroupContainerAll(mock.Anything, mock.Anything).
				Return(admin.ListGroupContainerAllApiRequest{ApiService: peeringAPI}).Maybe()
			peeringAPI.EXPECT().ListGroupContainerAllExecute(
				mock.Anything).Return(
				nil, nil, nil,
			).Maybe()

			workflowCtx := workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						NetworkPeeringApi: peeringAPI,
					},
				},
			}

			result := ensureNetworkPeers(&workflowCtx, prj)
			require.Equal(t, workflow.OK(), result)
		})
	}
}

func newNetworkPeeringTestProject(networkPeers []string) *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name:         "test-project",
			NetworkPeers: synthesizeNetworkPeerings(networkPeers),
		},
	}
}

func synthesizeNetworkPeerings(peeringIDs []string) []akov2.NetworkPeer {
	peers := make([]akov2.NetworkPeer, 0, len(peeringIDs))
	for _, id := range peeringIDs {
		peers = append(peers, akov2.NetworkPeer{
			ProviderName: "AWS",
			VpcID:        id,
			ContainerID:  fmt.Sprintf("container-%s", id),
		})
	}
	return peers
}

func synthesizeAtlasNetworkPeerings(peeringIDs []string) *admin.PaginatedContainerPeer {
	atlasPeers := make([]admin.BaseNetworkPeeringConnectionSettings, 0, len(peeringIDs))
	for _, id := range peeringIDs {
		atlasPeers = append(atlasPeers, admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:  fmt.Sprintf("container-%s", id),
			Id:           pointer.MakePtr(fmt.Sprintf("np-%s", id)),
			ProviderName: pointer.MakePtr("AWS"),
			VpcId:        pointer.MakePtr(id),
			StatusName:   pointer.MakePtr(StatusReady),
		})
	}
	return &admin.PaginatedContainerPeer{
		Results: &atlasPeers,
	}
}
