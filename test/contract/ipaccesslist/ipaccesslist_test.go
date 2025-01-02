package ipaccesslist

import (
	"context"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
)

func TestList(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "get default auditing", func(ch contract.ContractHelper) {
		projectName := "default-auditing-project"

		prj := contract.DefaultAtlasProject(projectName).(*akov2.AtlasProject)

		for i := 0; i < 110; i++ {
			prj.Spec.ProjectIPAccessList = append(prj.Spec.ProjectIPAccessList, project.IPAccessList{IPAddress: generateRandomIP()})
		}

		require.NoError(t, ch.AddResources(ctx, 10*time.Minute, prj))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)

		ipService := ipaccesslist.NewIPAccessList(ch.AtlasClient().ProjectIPAccessListApi)
		entries, err := ipService.List(ctx, testProjectID)
		require.NoError(t, err)
		assert.Equal(t, 110, len(entries))
	})
}

func generateRandomIP() string {
	ip := make(net.IP, net.IPv4len)
	for i := 0; i < net.IPv4len; i++ {
		ip[i] = byte(rand.Intn(256))
	}
	return ip.String()
}
