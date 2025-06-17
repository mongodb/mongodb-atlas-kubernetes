// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"io"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
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
