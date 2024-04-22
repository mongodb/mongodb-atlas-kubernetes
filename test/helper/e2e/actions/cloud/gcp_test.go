package cloud

import (
	"context"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateUsedVirtualAddress(t *testing.T) {
	if _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !ok {
		t.Skipf("Can't run test without GCE access credentials")
	}
	ctx := context.Background()
	gt := ginkgo.GinkgoT()
	ga, err := NewGCPAction(gt, GoogleProjectID)
	require.NoError(t, err)

	err = ga.createVirtualAddress(ctx, "10.3.0.55", "name1", Subnet2Name, GCPRegion)
	require.NoError(t, err)
	defer ga.deleteVirtualAddress(ctx, "name1", GCPRegion)
	expectedErr := ga.createVirtualAddress(ctx, "10.3.0.55", "name2", Subnet2Name, GCPRegion)
	require.ErrorContains(t, expectedErr, "IP_IN_USE_BY_ANOTHER_RESOURCE")
}
