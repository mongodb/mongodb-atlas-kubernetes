package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	pollingInterval  = time.Second * 10
	pollingDuration  = time.Minute * 10
	defaultNamespace = "default"
)

func setupLogger() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()
	if err != nil {
		zap.S().Errorf("Error building logger config: %s", err)
		os.Exit(1)
	}

	return log.Sugar()
}

// createK8sClient creates an in cluster client which can be used to fetch the current state of the AtlasCluster
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

	k8sClient.Scheme().AddKnownTypes(schema.GroupVersion{Group: "atlas.mongodb.com", Version: "v1"}, &mdbv1.AtlasCluster{}, &mdbv1.AtlasClusterList{})
	return k8sClient, nil
}

// getNamespace returns the current namespace.
func getNamespace() (string, error) {
	data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
		return ns, nil
	}
	return defaultNamespace, nil
}

// isClusterReady returns a boolean indicating if the cluster has reached the ready state and is
// ready to be used.
func isClusterReady(logger *zap.SugaredLogger) (bool, error) {
	k8sClient, err := createK8sClient()
	if err != nil {
		return false, err
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	clusterName := os.Getenv("CLUSTER_NAME")

	namespace, err := getNamespace()
	if err != nil {
		return false, err
	}

	totalTime := time.Duration(0)
	for range ticker.C {
		if totalTime > pollingDuration {
			break
		}
		totalTime += pollingInterval

		atlasCluster := mdbv1.AtlasCluster{}
		if err := k8sClient.Get(context.TODO(), kube.ObjectKey(namespace, clusterName), &atlasCluster); err != nil {
			return false, err
		}

		// the atlas project has reached the ClusterReady state.
		for _, cond := range atlasCluster.Status.Conditions {
			if cond.Type == status.ClusterReadyType {
				if cond.Status == corev1.ConditionTrue {
					return true, nil
				}
				logger.Infof("Atlas Cluster %s is not yet ready", atlasCluster.Name)
			}
		}
	}
	return false, nil
}

func main() {
	logger := setupLogger()

	clusterIsReady, err := isClusterReady(logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	exitCode := 1
	if clusterIsReady {
		exitCode = 0
	}
	os.Exit(exitCode)
}
