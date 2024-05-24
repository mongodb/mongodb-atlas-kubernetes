package auditing

import (
	"context"
	_ "embed"
	"log"
	"testing"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/resource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/launcher" // TODO: (DELETE ME) Launcher should be in place soon, before the template is used

	"github.com/stretchr/testify/assert"
)

//go:embed test.yml
var testYml string

const (
	testVersion = "2.1.0"
)

// TestMain prepares the tests for Atlas API calls
// Tests are skipped unless env var AKO_CONTRACT_TEST is set (to protect unit tests)
// Cleanup is skipped with SKIP_CLEANUP if we want to iterate faster locally
func TestMain(m *testing.M) {
	if !control.Enabled("AKO_CONTRACT_TEST") {
		log.Printf("Skipping contract test as AKO_CONTRACT_TEST is unset")
		return
	}
	l := launcher.NewFromEnv(testVersion)
	l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", 30*time.Second))
	if !control.Enabled("SKIP_CLEANUP") { // allow to reuse Atlas resources for local tests
		defer l.Cleanup()
	}
	m.Run()
}

func TestGetResource(t *testing.T) {
	ctx := context.Background()
	rs := resource.NewProductionResources( /*something*/ )

	assert.PanicsWithValue(t, "unimplemented", func() {
		rs.GetResource(ctx, "some-id")
	})
}
