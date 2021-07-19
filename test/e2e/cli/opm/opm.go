package opm

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"strings"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli"
)

// Version print version
func Version() {
	session := cli.Execute("opm", "version")
	session.Wait()
}

func AddIndex(bundleImage string) string {
	catalogName := genIndexCatalogName(bundleImage)
	session := cli.Execute("opm", "index", "add", "--bundles", bundleImage, "--tag", catalogName)
	Eventually(session).Should(gexec.Exit(0))
	return catalogName
}

// genIndexCatalogName gen Index Catalog image name from bundle image URL
func genIndexCatalogName(bundleImage string) string {
	url := strings.Split(bundleImage, ":")
	Expect(len(url)).Should(Equal(2), "Can't split DOCKER IMAGE")
	split := strings.Split(url[0], "/")
	name := split[len(split)-1]
	return strings.Replace(bundleImage, name, name+"-index-catalog", 1)
}
