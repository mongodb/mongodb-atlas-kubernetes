package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
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

func validateEnvVars() error {
	if !hasExactlyOneEnv("CLUSTER_NAME", "USER_NAME") {
		return fmt.Errorf("must specify exactly one of [CLUSTER_NAME, USER_NAME]")
	}
	return nil
}

func hasExactlyOneEnv(envVarNames ...string) bool {
	foundOne := false
	for _, envVarName := range envVarNames {
		envVar := os.Getenv(envVarName)
		if envVar != "" {
			if foundOne {
				return false
			}
			foundOne = true
		}
	}
	return foundOne
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

	k8sClient.Scheme().AddKnownTypes(schema.GroupVersion{Group: "atlas.mongodb.com", Version: "v1"}, &mdbv1.AtlasCluster{}, &mdbv1.AtlasClusterList{}, &mdbv1.AtlasDatabaseUser{}, &mdbv1.AtlasDatabaseUserList{})
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

// isReady returns a boolean indicating whether or not the given checker's resource is
// ready.
func isReady(readinessChecker atlasResourceReadinessChecker, logger *zap.SugaredLogger) (bool, error) {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

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

		isReady, err := readinessChecker.isReady(namespace, logger)
		if err != nil {
			return false, err
		}
		if isReady {
			return true, nil
		}
	}
	return false, nil
}

func shouldCheckForUser() bool {
	return os.Getenv("USER_NAME") != ""
}

func main() {
	logger := setupLogger()
	if err := validateEnvVars(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	k8sClient, err := createK8sClient()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	var checker atlasResourceReadinessChecker = &clusterReadinessChecker{Client: k8sClient}
	if shouldCheckForUser() {
		checker = &userReadinessChecker{Client: k8sClient}
	}

	ready, err := isReady(checker, logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	exitCode := 1
	if !ready {
		exitCode = 0
	}
	os.Exit(exitCode)
}
