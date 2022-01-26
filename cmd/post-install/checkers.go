package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

type atlasResourceReadinessChecker interface {
	isReady(namespace string, logger *zap.SugaredLogger) (bool, error)
	name() string
}

type clusterReadinessChecker struct {
	client.Client
}

func (c *clusterReadinessChecker) name() string {
	return os.Getenv("CLUSTER_NAME")
}

func (c *clusterReadinessChecker) isReady(namespace string, logger *zap.SugaredLogger) (bool, error) {
	atlasCluster := mdbv1.AtlasCluster{}
	if err := c.Get(context.TODO(), kube.ObjectKey(namespace, c.name()), &atlasCluster); err != nil {
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
	return false, nil
}

type userReadinessChecker struct {
	client.Client
}

func (c *userReadinessChecker) name() string {
	return os.Getenv("USER_NAME")
}

func (c *userReadinessChecker) isReady(namespace string, logger *zap.SugaredLogger) (bool, error) {
	atlasDatabaseUser := mdbv1.AtlasDatabaseUser{}
	if err := c.Get(context.TODO(), kube.ObjectKey(namespace, c.name()), &atlasDatabaseUser); err != nil {
		return false, err
	}

	// the atlas user has reached the ready state.
	for _, cond := range atlasDatabaseUser.Status.Conditions {
		if cond.Type == status.DatabaseUserReadyType {
			if cond.Status == corev1.ConditionTrue {
				return true, nil
			}
			logger.Infof("Atlas Database User %s is not yet ready", atlasDatabaseUser.Name)
		}
	}
	return false, nil
}
