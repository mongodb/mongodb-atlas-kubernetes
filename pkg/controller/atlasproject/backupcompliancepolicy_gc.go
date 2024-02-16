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
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
)

func (r *AtlasProjectReconciler) garbageCollectBackupResource(ctx context.Context, project *mdbv1.AtlasProject) error {
	policies := &mdbv1.AtlasBackupCompliancePolicyList{}

	err := r.Client.List(ctx, policies)
	if err != nil {
		return fmt.Errorf("failed to retrieve list of backup compliance policies: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, p := range policies.Items {
		policy := p
		g.Go(func() error {
			// get policy's associated projects
			annotation := policy.Annotations[ProjectAnnotation]
			annotations := strings.Split(annotation, ",")
			if !slices.Contains(annotations, project.ID()) {
				// project not covered by policy
				return nil
			}
			// policy thinks it covers project
			if project.Spec.BackupCompliancePolicyRef.Name == policy.Name ||
				project.Spec.BackupCompliancePolicyRef.Namespace == policy.Namespace {
				// project is still using the BCP
				r.Log.Debugw("adding deletion finalizer", "name", customresource.FinalizerLabel)
				customresource.SetFinalizer(policy, customresource.FinalizerLabel)
				return nil
			}

			// TODO: remove project ID annotation

			if policy.GetDeletionTimestamp().IsZero() {
				if projects, ok := policy.Annotations[ProjectAnnotation]; ok {
					if len(strings.Split(projects, ",")) == 0 {
						r.Log.Debugw("removing deletion finalizer", "name", customresource.FinalizerLabel)
						customresource.UnsetFinalizer(policy, customresource.FinalizerLabel)
					}
				}
			}

			if !policy.GetDeletionTimestamp().IsZero() && customresource.HaveFinalizer(policy, customresource.FinalizerLabel) {
				r.Log.Warnf("backup compliance policy %s is assigned to at least one Project. Remove it from all Projects before deletion", policy.Name)
			}
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}
