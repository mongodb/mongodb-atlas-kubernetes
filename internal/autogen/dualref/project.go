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

package dualref

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/atlas"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

var (
	ErrMissingKubeProject = errors.New("missing Kubernetes Atlas Project")
)

func ResolveProject(ctx context.Context, kubeClient client.Client, pro project.ProjectReferrerObject) (*project.Project, error) {
	atlasClientSet := atlas.FromContext(ctx)
	projectsService := project.NewProjectAPIService(atlasClientSet.SdkClient20231115008.ProjectsApi)
	ref := pro.ProjectDualRef()
	if ref.ProjectRef != nil {
		project, err := projectFromKube(ctx, kubeClient, pro, projectsService)
		if err != nil {
			return nil, fmt.Errorf("failed to get project via Kubernetes reference: %w", err)
		}
		return project, nil
	}

	prj, err := projectFromAtlas(ctx, projectsService, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get project via Atlas by ID: %w", err)
	}

	return prj, nil
}

func projectFromKube(ctx context.Context, kubeClient client.Client, pro project.ProjectReferrerObject, service project.ProjectService) (*project.Project, error) {
	kubeProject, err := fetchProject(ctx, kubeClient, pro)
	if err != nil {
		return nil, err
	}

	return service.GetProjectByName(ctx, kubeProject.Spec.Name)
}

func projectFromAtlas(ctx context.Context, projectsService project.ProjectService, pdr *akov2.ProjectDualReference) (*project.Project, error) {
	return projectsService.GetProject(ctx, pdr.ExternalProjectRef.ID)
}

func fetchProject(ctx context.Context, kubeClient client.Client, pro project.ProjectReferrerObject) (*akov2.AtlasProject, error) {
	ref := pro.ProjectDualRef()
	if ref.ProjectRef == nil {
		return nil, nil
	}

	project := akov2.AtlasProject{}
	ns := pro.GetNamespace()
	if ref.ProjectRef.Namespace != "" {
		ns = ref.ProjectRef.Namespace
	}

	key := client.ObjectKey{Name: ref.ProjectRef.Name, Namespace: ns}
	err := kubeClient.Get(ctx, key, &project)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.Join(ErrMissingKubeProject, err)
		}
		return nil, fmt.Errorf("error getting AtlasProject: %w", err)
	}

	return &project, nil
}
