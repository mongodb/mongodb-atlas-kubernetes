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

package helm

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

// dependencyAsFile changes Chart dependency from 'repository' to 'file'
// used for not released dependencies
func dependencyAsFileForCRD() {
	chart := filepath.Join(config.AtlasOperatorHelmChartPath, "Chart.yaml")
	data, _ := os.ReadFile(filepath.Clean(chart))
	r, err := regexp.Compile("repository: \"https://mongodb.github.io/helm-charts\"") //nolint:gocritic // if a test runs locally, it could be already changed
	if err == nil {
		utils.SaveToFile(
			chart,
			r.ReplaceAll(data, []byte("repository:  \"file://"+config.AtlasOperatorCRDHelmChartPath+"\"")),
		)
	}
}
