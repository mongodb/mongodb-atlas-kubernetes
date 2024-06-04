package launcher_test

import (
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/launcher"
)

//go:embed test.yml
var testYml string

const (
	testVersion = "2.1.0"
)

func TestLaunch(t *testing.T) {
	control.SkipTestUnless(t, "AKO_LAUNCHER_TEST")

	l := launcher.NewFromEnv(testVersion)
	timeout := time.Minute
	assert.NoError(t, l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", timeout)))
	// retest should also work and faster (no need to re-create resources)
	assert.NoError(t, l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", timeout)))
	require.NoError(t, l.Cleanup())
}
