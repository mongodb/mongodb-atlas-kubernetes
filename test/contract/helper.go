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

func newRandomName(prefix string) string {
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

type ActionFunc func(ctx context.Context)

func TestMain(m *testing.M, setup, clear ActionFunc) {
	if !control.Enabled("AKO_CONTRACT_TEST") {
		log.Print("Skipping contract tests, AKO_CONTRACT_TEST is not set")
		return
	}
	ctx := context.Background()
	setup(ctx)
	code := m.Run()
	clear(ctx)
	os.Exit(code)
}
