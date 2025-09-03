package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/josvazg/crd2go/internal/checkerr"
)

const (
	// visit https://github.com/s-urbaniak/atlas2crd/blob/main/crds.yaml
	CRDsURL = "https://raw.githubusercontent.com/s-urbaniak/atlas2crd/refs/heads/main/crds.yaml?token=%s"

	samplesDir = "pkg/crd2go/samples/"

	crdsFile = samplesDir + "crds.yaml"

	targetDir = samplesDir + "v1"

	controllerGenPaths = "./" + targetDir
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
	if err := run(crd2go, "-input", crdsFile, "-output", targetDir); err != nil {
		return fmt.Errorf("failed to generate CRDs Go structs: %w", err)
	}

	log.Print("Generating Go deep copy code...")
	if err := run("controller-gen", "object", fmt.Sprintf("paths=%q", controllerGenPaths)); err != nil {
		return fmt.Errorf("failed to generate Go deep copy code: %w", err)
	}
	return nil
}

func downloadTo(url, filename string) (int64, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to download from %s: %w", url, err)
	}
	f, err := os.Create(filename)
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

func run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
