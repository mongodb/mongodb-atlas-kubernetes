package cli

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

func Execute(command string, args ...string) *gexec.Session {
	// GinkgoWriter.Write([]byte("\n " + command + " " + strings.Join(args, " "))) // TODO for the local run only
	cmd := exec.Command(command, args...)
	session, _ := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	return session
}