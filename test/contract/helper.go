package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
)

func NewAPIClient() (*admin.APIClient, error) {
	client, err := admin.NewClient(
		admin.UseBaseURL(Domain()),
		admin.UseDigestAuth(PublicAPIKey(), PrivateAPIKey()),
		admin.UseUserAgent(contractTestUserAgent()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %w", err)
	}
	return client, nil
}

func DefaultProviderName() string {
	return "AWS"
}

func DefaultRegion() string {
	return "US_EAST_2"
}

func contractTestUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBContractTestsAKO", version.Version, runtime.GOOS, runtime.GOARCH)
}

func OrgID() string {
	return mustGetEnv("MCLI_ORG_ID")
}

func Domain() string {
	return mustGetEnv("MCLI_OPS_MANAGER_URL")
}

func PublicAPIKey() string {
	return mustGetEnv("MCLI_PUBLIC_API_KEY")
}

func PrivateAPIKey() string {
	return mustGetEnv("MCLI_PRIVATE_API_KEY")
}

func BoolEnv(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value == strings.ToLower("1")
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic("expected MCLI_ORG_ID was not set")
	}
	return value
}

func NewRandomName(prefix string) string {
	randomSuffix := uuid.New().String()[0:6]
	return fmt.Sprintf("%s-%s", prefix, randomSuffix)
}

func Jsonize(obj any) string {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

func SkipContractTesting() bool {
	if !control.Enabled("AKO_CONTRACT_TEST") {
		log.Print("Skipping contract tests, AKO_CONTRACT_TEST is not set")
		return true
	}
	return false
}

type setupFunc func(context.Context) (*TestResources, error)

func RunTests(m *testing.M, resourcesVar **TestResources, setupFn setupFunc) int {
	if SkipContractTesting() {
		return 0
	}
	ctx := context.Background()
	resources, err := setupFn(ctx)
	if err != nil {
		log.Print(err.Error())
		return 1
	}
	*resourcesVar = resources
	defer resources.Recycle(ctx)
	return m.Run()
}
