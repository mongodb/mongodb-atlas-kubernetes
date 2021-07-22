package oc

import (
	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

func Version() {
	session := cli.Execute("oc", "version")
	session.Wait()
}

func Login(code string) {
	session := cli.ExecuteWithoutWriter("oc", "login", "--token="+code, "--server=https://api.openshift.mongokubernetes.com:6443", "--insecure-skip-tls-verify")
	session.Wait("2m")
}

func Apply(path string) {
	session := cli.Execute("oc", "apply", "-f", path)
	session.Wait("2m")
}

func Delete(path string) {
	session := cli.Execute("oc", "delete", "-f", path)
	session.Wait("2m")
}
