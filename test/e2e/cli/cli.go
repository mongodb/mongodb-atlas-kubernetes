package cli

import (
	"io"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/gomega/gbytes"
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

func SessionShouldExit(session *gexec.Session) {
	EventuallyWithOffset(
		2,
		func() int { return session.ExitCode() },
		"3m",
		"5s",
	).ShouldNot(Equal(-1))
}

func GetSessionExitMsg(session *gexec.Session) *Buffer {
	SessionShouldExit(session)
	if session.ExitCode() != 0 {
		return session.Err
	}
	return session.Out
}
