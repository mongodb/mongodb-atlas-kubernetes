package podman

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

func Version() {
	session := cli.Execute("podman", "version")
	session.Wait()
}

func Login(registry, user, pass string) {
	session := cli.Execute("podman", "login", "-u", user, "-p", pass, registry)
	Eventually(session, "2m", "10s").Should(gexec.Exit(0))
}

func PushIndexCatalog(catalogURL string) {
	session := cli.Execute("podman", "push", catalogURL)
	Eventually(session, "5m", "10s").Should(gexec.Exit(0))
}
