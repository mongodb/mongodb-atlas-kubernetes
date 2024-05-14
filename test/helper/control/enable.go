package control

import (
	"os"
	"strings"
	"testing"
)

func Enabled(envvar string) bool {
	value := strings.ToLower(os.Getenv(envvar))
	return value == "1"
}

func SkipTestUnless(t *testing.T, envvar string) {
	if !Enabled(envvar) {
		t.Skipf("Skipping tests, %s is not set", envvar)
	}
}
