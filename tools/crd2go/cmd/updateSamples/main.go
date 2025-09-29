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
//

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/checkerr"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/fileinput"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/run"
)

const (
	// visit https://github.com/s-urbaniak/atlas2crd/blob/main/crds.yaml
	CRDsURL = "https://raw.githubusercontent.com/s-urbaniak/atlas2crd/refs/heads/main/crds.yaml?token=%s"

	samplesDir = "./pkg/crd2go/samples/"

	crdsFile = samplesDir + "crds.yaml"

	targetDir = samplesDir + "v1"
)

func main() {
	if err := updateSamples(); err != nil {
		log.Fatal(err)
	}
}

func updateSamples() error {
	token := mustGetenv("GITHUB_URL_TOKEN")
	url := fmt.Sprintf(CRDsURL, token)
	log.Printf("Downloading %s on to %s", CRDsURL, crdsFile)
	n, err := downloadTo(url, crdsFile)
	if err != nil {
		return fmt.Errorf("failed to download CRD YAML: %w", err)
	}
	log.Printf("Downloaded %d bytes on to %s", n, crdsFile)

	log.Printf("Generating Go structs from CRDs to %s...", targetDir)
	crd2go := mustGetenv("CRD2GO_BIN")
	if err := run.Run(crd2go, "-input", crdsFile, "-output", targetDir); err != nil {
		return fmt.Errorf("failed to generate CRDs Go structs: %w", err)
	}
	return nil
}

func downloadTo(url, filename string) (int64, error) {
	// #nosec G107 URL is safe as we are just adding a token, it cannot be re-pathed
	//nolint:noctx
	rsp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to download from %s: %w", url, err)
	}
	if rsp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to request %s with status: %q", url, rsp.Status)
	}
	f, err := os.Create(fileinput.MustBeSafe(filename))
	if err != nil {
		return 0, fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer checkerr.CheckErr("closing download file", f.Close)
	n, err := io.Copy(f, rsp.Body)
	if err != nil {
		return n, fmt.Errorf("failed to write downloaded data to file %s: %w", filename, err)
	}
	return n, nil
}

func mustGetenv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("%s env var must be set", name))
	}
	return value
}
