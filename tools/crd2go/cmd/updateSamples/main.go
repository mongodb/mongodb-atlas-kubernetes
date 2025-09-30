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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	debug = false

	OpenAPI2CRDDir = "../openapi2crd"

	samplesDir = "../crd2go/pkg/crd2go/samples/"

	crdsFile = samplesDir + "crds.yaml"
)

func main() {
	if err := updateSamples(); err != nil {
		log.Fatal(err)
	}
}

func updateSamples() error {
	if err := generateCRDs(); err != nil {
		return fmt.Errorf("CRD generation failed: %w", err)
	}
	if err := generateSamples(); err != nil {
		return fmt.Errorf("Samples generation failed: %w", err)
	}
	return nil
}

func generateCRDs() error {
	return runAt(OpenAPI2CRDDir,
		"go", "run", "main.go", "-c", "config.yaml", "-o", crdsFile)
}

func generateSamples() error {
	return runAt(".",
		"go", "run", "cmd/crd2go/main.go", "-config", "crd2go.yaml")
}

func runAt(dir, command string, args ...string) error {
	debugIt("%s > %s %s", dir, command, strings.Join(args, " "))
	cmd := exec.CommandContext(context.Background(), command, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, err := os.Stderr.Write(out); err != nil {
			log.Printf("err write failed: %v", err)
		}
		return fmt.Errorf("run at directory %s failed: %w", dir, err)
	}
	debugIt("output:\n%s", string(out))
	return nil
}

func debugIt(msg string, args ...any) {
	if !debug {
		return
	}
	log.Printf(msg, args...)
}
