package contract_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
)

const (
	disabledPrefix = "DISABLED_ENV_VAR_"

	modifiedEnvVars = "MODIFIED_ENV_VARS"
)

func TestContractTestSkip(t *testing.T) {
	testWithEnv(func() {
		contract.RunContractTest(t, "Skip contract test", func(_ *contract.ContractTest) {
			panic("should not have got here!")
		})
	}, "-AKO_CONTRACT_TEST")
}

func TestContractTestClientSetFails(t *testing.T) {
	testWithEnv(func() {
		assert.Panics(t, func() {
			contract.RunContractTest(t, "bad client settings panics", func(_ *contract.ContractTest) {})
		})
	},
		"AKO_CONTRACT_TEST=1",
		"-MCLI_OPS_MANAGER_URL",
		"-MCLI_PUBLIC_API_KEY",
		"-MCLI_PRIVATE_API_KEY")
}

func TestContractsWithResources(t *testing.T) {
	contract.RunContractTest(t, "run contract test list projects", func(ct *contract.ContractTest) {
		ct.AddResources(time.Minute, contract.DefaultAtlasProject("contract-tests-list-projects"))
		_, _, err := ct.AtlasClient.ProjectsApi.ListProjects(ct.Ctx).Execute()
		assert.NoError(t, err)
	})
	contract.RunContractTest(t, "run contract test list orgs", func(ct *contract.ContractTest) {
		ct.AddResources(time.Minute, contract.DefaultAtlasProject("contract-tests-list-orgs"))
		_, _, err := ct.AtlasClient.OrganizationsApi.ListOrganizations(ct.Ctx).Execute()
		assert.NoError(t, err)
	})
}

func testWithEnv(testFn func(), envEntries ...string) {
	for _, entry := range envEntries {
		if entry[0] == '-' {
			disableEnv(entry[1:])
			continue
		}
		parts := strings.Split(entry, "=")
		if len(parts) != 2 {
			panic(fmt.Sprintf("expected 'key=value' but got %q", entry))
		}
		setTestEnv(parts[0], parts[1])
	}
	defer restoreEnvs()
	testFn()
}

func disableEnv(varName string) {
	value := os.Getenv(varName)
	os.Setenv(disabledPrefix+varName, value)
	os.Unsetenv(varName)
	registerToRestore(varName)
}

func restoreEnvs() {
	envVars := os.Getenv(modifiedEnvVars)
	for _, varName := range strings.Split(envVars, ",") {
		value := os.Getenv(disabledPrefix + varName)
		os.Setenv(varName, value)
		os.Unsetenv(disabledPrefix + varName)
	}
	os.Unsetenv(modifiedEnvVars)
}

func setTestEnv(varName, value string) {
	previousValue := os.Getenv(varName)
	os.Setenv(disabledPrefix+varName, previousValue)
	os.Setenv(varName, value)
	registerToRestore(varName)
}

func registerToRestore(varName string) {
	envVars := os.Getenv(modifiedEnvVars)
	if envVars == "" {
		envVars = varName
	} else {
		envVars = strings.Join([]string{envVars, varName}, ",")
	}
	os.Setenv(modifiedEnvVars, envVars)
}
