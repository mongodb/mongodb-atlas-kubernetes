package crd2go

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/josvazg/crd2go/internal/crd"
	"github.com/josvazg/crd2go/internal/gen"
	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/pkg/config"
)

func LoadConfig(r io.Reader) (*config.Config, error) {
	yml, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %w", err)
	}
	cfg := config.Config{}
	if err = yaml.Unmarshal(yml, &cfg); err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %v", err)
	}
	return &cfg, nil
}

// CodeWriterAtPath creates a file writer for the given CRD at the specified directory
func CodeWriterAtPath(dir string) config.CodeWriterFunc {
	return func(filename string, overwrite bool) (io.WriteCloser, error) {
		srcFile := filepath.Join(dir, filename)
		flags := os.O_CREATE | os.O_EXCL | os.O_WRONLY
		if overwrite {
			flags = os.O_CREATE | os.O_TRUNC | os.O_RDWR
		}
		w, err := os.OpenFile(srcFile, flags, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", srcFile, err)
		}
		return w, nil
	}
}

// GenerateToDir generates Go code from a CRD YAML file into a directory
func GenerateToDir(cfg *config.Config) error {
	in, err := os.Open(cfg.Input)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", cfg.Input, err)
	}
	req := gotype.Request{
		CoreConfig:   cfg.CoreConfig,
		CodeWriterFn: CodeWriterAtPath(cfg.Output),
		TypeDict:     gotype.NewTypeDict(cfg.Renames, gotype.KnownTypes()...),
	}
	return Generate(&req, in)
}

// Generate will write files using the CodeWriterFunc
func Generate(req *gotype.Request, r io.Reader) error {
	groupsVersions := map[string]struct{}{}
	group := ""
	version := ""
	generatedGVRs, err := GenerateStream(req, r)
	if err != nil {
		return fmt.Errorf("failed to generate CRDs: %w", err)
	}
	for _, gvr := range generatedGVRs {
		parts := strings.Split(gvr, "/")
		if len(parts) > 2 {
			group = parts[0]
			version = parts[1]
			gv := fmt.Sprintf("%s/%s", group, version)
			groupsVersions[gv] = struct{}{}
		}
	}
	if len(groupsVersions) == 1 {
		if err := gen.GenerateGroupVersionFiles(req, group, version); err != nil {
			return fmt.Errorf("failed to generate files for group version '%s/%s': %w", group, version, err)
		}
	}
	return nil
}

// GenerateStream generates Go code from a stream of CRDs within a YAML reader.
// It uses the provided CodeWriterFunc to write the generated code to the specified output.
// The version parameter specifies the version of the CRD to generate code for.
// The preloadedTypes parameter allows for preloading specific types to avoid name collisions.
func GenerateStream(req *gotype.Request, r io.Reader) ([]string, error) {
	preloaded := []*gotype.GoType{}
	for _, name := range req.Reserved {
		preloaded = append(preloaded, gotype.NewOpaqueType(name))
	}
	for _, importType := range req.Imports {
		preloaded = append(preloaded, gotype.NewAutoImportType(&importType))
	}
	overwrite := true
	generatedGVRs := []string{}
	generated := false
	scanner := bufio.NewScanner(r)
	req.TypeDict.AddAll(preloaded...)
	for {
		crdSchema, err := crd.ParseCRD(scanner)
		if errors.Is(err, io.EOF) {
			if generated {
				return generatedGVRs, nil
			}
			return nil, fmt.Errorf("failed to parse CRD: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		if in(req.SkipList, crdSchema.Spec.Names.Kind) {
			continue
		}
		w, err := req.CodeWriterFn(crd.CRD2Filename(crdSchema), overwrite)
		if err != nil {
			return nil, fmt.Errorf("failed to get writer for CRD %s: %w", crdSchema.Name, err)
		}
		defer w.Close()
		versionedCRD := crd.SelectVersion(&crdSchema.Spec, req.Version)
		if versionedCRD == nil {
			if req.Version == "" {
				return nil, fmt.Errorf("no versions to generate code from")
			}
			return nil, fmt.Errorf("no version %q to generate code from", req.Version)
		}
		stmt, err := gen.GenerateCRD(req.TypeDict, versionedCRD)
		if err != nil {
			return nil, fmt.Errorf("failed to generate CRD code: %w", err)
		}
		if _, err := w.Write(([]byte)(stmt.GoString())); err != nil {
			return nil, fmt.Errorf("failed to write Go code: %w", err)
		}
		generated = true
		gvr := fmt.Sprintf("%s/%s/%s",
			crdSchema.Spec.Group, versionedCRD.Version.Name, crdSchema.Spec.Names.Plural)
		gvr = strings.TrimPrefix(gvr, "/")
		generatedGVRs = append(generatedGVRs, gvr)
	}
}

func in[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
