package run

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run(command string, args ...string) error {
	log.Printf("Running:\n  %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
