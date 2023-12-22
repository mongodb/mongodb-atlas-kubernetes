/*
Copyright 2023 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasproject

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func (r *AtlasProjectReconciler) garbageCollectBackupResource(ctx context.Context, clusterName string) error {
	policies := &mdbv1.AtlasBackupCompliancePolicyList{}

	err := r.Client.List(ctx, policies)
	if err != nil {
		return fmt.Errorf("failed to retrieve list of backup schedules: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, policy := range policies.Items {
		g.Go(func() error {
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}
