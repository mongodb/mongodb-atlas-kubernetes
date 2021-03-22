package cli

import (
	"io"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

func Execute(command string, args ...string) *gexec.Session {
	GinkgoWriter.Write([]byte("\n " + command + " " + strings.Join(args, " ") + "\n"))
	return ExecuteCommand(GinkgoWriter, command, args...)
}

func ExecuteWithoutWriter(command string, args ...string) *gexec.Session {
	return ExecuteCommand(nil, command, args...)
}

func ExecuteCommand(reporter io.Writer, command string, args ...string) *gexec.Session {
	cmd := exec.Command(command, args...)
	session, _ := gexec.Start(cmd, reporter, reporter)
	return session
}
