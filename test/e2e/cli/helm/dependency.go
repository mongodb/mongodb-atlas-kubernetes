package helm

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// dependencyAsFile changes Chart dependency 'repository' >> 'file'
// for not released dependencies
func dependencyAsFileForCRD() {
	chart := filepath.Join(config.AtlasOperatorHelmChartPath, "Chart.yaml")
	data, _ := ioutil.ReadFile(filepath.Clean(chart))
	r, err := regexp.Compile("repository: \"https://mongodb.github.io/helm-charts\"") // nolint:gocritic // if a test runs locally, it could be already changed
	if err == nil {
		utils.SaveToFile(
			chart,
			r.ReplaceAll(data, []byte("repository:  \"file://"+config.AtlasOperatorCRDHelmChartPath+"\"")),
		)
	}
}
