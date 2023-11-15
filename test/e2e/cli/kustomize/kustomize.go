package kustomize

import (
	cli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/cli"
)

func Version() {
	session := cli.Execute("kustomize", "version")
	session.Wait()
}

func Build(source string) []byte {
	session := cli.Execute("kustomize", "build", source)
	session.Wait("2m")
	return session.Out.Contents()
}
