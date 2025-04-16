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

package main

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

const (
	pollingInterval = time.Second * 10
	pollingDuration = time.Minute * 30
)

func setupLogger() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()
	if err != nil {
		zap.S().Errorf("Error building logger config: %s", err)
		os.Exit(1)
	}

	return log.Sugar()
}

// createK8sClient creates a client which can be used to fetch the current state of the AtlasDeployment
// resource.
func createK8sClient() (client.Client, error) {
	restCfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(restCfg, client.Options{})

	if err != nil {
		return nil, err
	}

	k8sClient.Scheme().AddKnownTypes(schema.GroupVersion{Group: "atlas.mongodb.com", Version: "v1"}, &akov2.AtlasDeployment{}, &akov2.AtlasDeploymentList{})
	return k8sClient, nil
}

// isDeploymentReady returns a boolean indicating if the deployment has reached the ready state and is
// ready to be used.
func isDeploymentReady(ctx context.Context, logger *zap.SugaredLogger) (bool, error) {
	k8sClient, err := createK8sClient()
	if err != nil {
		return false, err
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	deploymentName := os.Getenv("DEPLOYMENT_NAME")
	namespace := os.Getenv("NAMESPACE")

	totalTime := time.Duration(0)
	for range ticker.C {
		if totalTime > pollingDuration {
			break
		}
		totalTime += pollingInterval

		atlasDeployment := akov2.AtlasDeployment{}
		if err := k8sClient.Get(ctx, kube.ObjectKey(namespace, deploymentName), &atlasDeployment); err != nil {
			return false, err
		}

		// the atlas project has reached the DeploymentReady state.
		for _, cond := range atlasDeployment.Status.Conditions {
			if cond.Type == api.DeploymentReadyType {
				if cond.Status == corev1.ConditionTrue {
					return true, nil
				}
				logger.Infof("Atlas Deployment %s is not yet ready", atlasDeployment.Name)
			}
		}
	}
	return false, nil
}

func main() {
	ctx := signals.SetupSignalHandler()
	logger := setupLogger()

	deploymentIsReady, err := isDeploymentReady(ctx, logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	exitCode := 1
	if deploymentIsReady {
		exitCode = 0
	}
	os.Exit(exitCode)
}
