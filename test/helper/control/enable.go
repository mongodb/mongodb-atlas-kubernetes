package control

import (
	"os"
	"strings"
)

func Enabled(envvar string) bool {
	value := strings.ToLower(os.Getenv(envvar))
	return value == "1"
}
