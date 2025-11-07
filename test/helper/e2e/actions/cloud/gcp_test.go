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
	ga, err := NewGCPAction(ctx, gt, GoogleProjectID)
	require.NoError(t, err)

	err = ga.createVirtualAddress(ctx, "10.0.0.155", "name1", Subnet2Name, GCPRegion)
	require.NoError(t, err)
	defer ga.deleteVirtualAddress(ctx, "name1", GCPRegion)
	expectedErr := ga.createVirtualAddress(ctx, "10.0.0.155", "name2", Subnet2Name, GCPRegion)
	require.ErrorContains(t, expectedErr, "IP_IN_USE_BY_ANOTHER_RESOURCE")
}
