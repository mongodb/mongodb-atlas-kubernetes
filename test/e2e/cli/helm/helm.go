package kube

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// GenKubeVersion
func GetVersionOutput() {
	session := cli.Execute("helm", "version")
	ExpectWithOffset(1, session).Should(Say("version.BuildInfo"), "Please, install HELM")
}

func Uninstall(name string) {
	session := cli.Execute("helm", "uninstall", name)
	EventuallyWithOffset(1, session.Wait()).Should(Say("uninstalled"), "HELM. Can't unninstall "+name)
}

func Install(args ...string) {
	args = append([]string{"install"}, args...)
	session := cli.Execute("helm", args...)
	EventuallyWithOffset(1, session.Wait()).Should(Say("STATUS: deployed"), "HELM. Can't install release")
}

func InstallTestApplication(input model.UserInputs, user model.DBUser, port string) {
	Install(
		"test-app-"+user.Spec.Username,
		config.HELMTestAppPath,
		"--set-string", "secret="+input.ProjectName+"-"+input.Clusters[0].Spec.Name+"-"+user.Spec.Username,
		"--set-string", "nodePort="+port,
		"-n", input.Namespace,
	)
}
