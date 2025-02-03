package atlasnetworkpeering

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	peeringMock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

const (
	sampleProjectID = "fake-632842834"

	sampleContainerID = "container-id"
)

var errFakeFailure = errors.New("failure")

func TestHandleContainer(t *testing.T) {
	for _, tc := range []struct {
		title             string
		req               *reconcileRequest
		expectedContainer *networkpeering.ProviderContainer
		expectedError     error
	}{
		{
			title: "fail to create container",
			req: &reconcileRequest{
				workflowCtx: &workflow.Context{},
				service: func() networkpeering.NetworkPeeringService {
					netpeeringService := peeringMock.NewNetworkPeeringServiceMock(t)
					netpeeringService.EXPECT().GetContainer(mock.Anything, sampleProjectID, sampleContainerID).Return(nil, nil)
					netpeeringService.EXPECT().CreateContainer(mock.Anything, sampleProjectID, mock.Anything).Return(nil, errFakeFailure)
					return netpeeringService
				}(),
				projectID: sampleProjectID,
				networkPeering: &v1.AtlasNetworkPeering{
					Spec: v1.AtlasNetworkPeeringSpec{
						AtlasNetworkPeeringConfig: v1.AtlasNetworkPeeringConfig{
							Provider:    string(provider.ProviderAWS),
							ContainerID: sampleContainerID,
						},
						AtlasProviderContainerConfig: v1.AtlasProviderContainerConfig{},
					},
				},
			},
			expectedError: errFakeFailure,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			r := &AtlasNetworkPeeringReconciler{}
			container, err := r.handleContainer(tc.req)
			assert.Equal(t, tc.expectedContainer, container)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
