package oc

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/cli"
)

func Version() {
	session := cli.Execute("oc", "version")
	session.Wait()
}

func Login(code, serverAPI string) {
	session := cli.ExecuteWithoutWriter("oc", "login", "--token="+code, "--server="+serverAPI, "--insecure-skip-tls-verify")
	EventuallyWithOffset(1, session).Should(Say("Logged into"), "Can not login to "+serverAPI)
}

func Apply(path string) {
	session := cli.Execute("oc", "apply", "-f", path)
	session.Wait("2m")
}

func Delete(path string) {
	session := cli.Execute("oc", "delete", "-f", path)
	session.Wait("2m")
}
