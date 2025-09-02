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
	"github.com/josvazg/crd2go/internal/crd/hooks"
	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/internal/render"
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
		if err := render.Default.RenderDoc(req, group, version); err != nil {
			return fmt.Errorf("failed to generate the doc.go file for group version '%s/%s': %w", group, version, err)
		}
		if err := render.Default.RenderSchema(req, group, version); err != nil {
			return fmt.Errorf("failed to generate the schema.go file for group version '%s/%s': %w", group, version, err)
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
		versionedCRD := crd.SelectVersion(&crdSchema.Spec, req.Version)
		if versionedCRD == nil {
			if req.Version == "" {
				return nil, fmt.Errorf("no versions to generate code from")
			}
			return nil, fmt.Errorf("no version %q to generate code from", req.Version)
		}
		goCRD, err := buildCRDType(req, versionedCRD)
		if err != nil {
			return nil, fmt.Errorf("could not build CRD type: %w", err)
		}

		renderReq := render.CRDRenderRequest{
			Request:  *req,
			Filename: crd.Kind2Filename(versionedCRD.Kind),
			Version:  versionedCRD.Version.Name,
			Kind:     versionedCRD.Kind,
			Type:     goCRD,
		}
		if err := render.Default.RenderCRD(&renderReq); err != nil {
			return nil, fmt.Errorf("failed to generate CRD code: %w", err)
		}
		generated = true
		gvr := fmt.Sprintf("%s/%s/%s",
			crdSchema.Spec.Group, versionedCRD.Version.Name, crdSchema.Spec.Names.Plural)
		gvr = strings.TrimPrefix(gvr, "/")
		generatedGVRs = append(generatedGVRs, gvr)
	}
}

func buildCRDType(req *gotype.Request, versionedCRD *crd.VersionedCRD) (*gotype.GoType, error) {
	req.TypeDict.Add(gotype.NewStruct(versionedCRD.Kind, nil)) // reserve the name of the root type not to be taken
	specSchema := versionedCRD.Version.Schema.OpenAPIV3Schema.Properties["spec"]
	spec, err := crd.FromOpenAPIType(req.TypeDict, hooks.Hooks, &crd.CRDType{
		Name:    versionedCRD.SpecTypename(),
		Parents: []string{versionedCRD.Kind},
		Schema:  &specSchema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate spec type: %w", err)
	}

	statusSchema := versionedCRD.Version.Schema.OpenAPIV3Schema.Properties["status"]
	status, err := crd.FromOpenAPIType(req.TypeDict, hooks.Hooks, &crd.CRDType{
		Name:    versionedCRD.StatusTypename(),
		Parents: []string{versionedCRD.Kind},
		Schema:  &statusSchema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate status type: %w", err)
	}
	return combineSpecAndStatus(versionedCRD.Kind, spec, status)
}

func combineSpecAndStatus(kind string, spec, status *gotype.GoType) (*gotype.GoType, error) {
	metav1Package := config.ImportInfo{
		Alias: "metav1", Path: "k8s.io/apimachinery/pkg/apis/meta/v1",
	}
	typeMeta := gotype.NewAutoImportType(&config.ImportedTypeConfig{
		ImportInfo: metav1Package, Name: "TypeMeta",
	})
	objectMeta := gotype.NewAutoImportType(&config.ImportedTypeConfig{
		ImportInfo: metav1Package, Name: "ObjectMeta",
	})
	goCRD := gotype.NewStruct(kind, []*gotype.GoField{
		gotype.NewEmbeddedField(typeMeta).WithOptions(gotype.Required(true), gotype.JSONTag(",inline")),
		gotype.NewEmbeddedField(objectMeta).WithOptions(gotype.Required(true), gotype.JSONTag("metadata,omitempty")),
		gotype.NewGoField("Spec", spec).WithOptions(gotype.Required(true), gotype.JSONTag("spec,omitempty")),
		gotype.NewGoField("Status", status).WithOptions(gotype.Required(true), gotype.JSONTag("status,omitempty")),
	})
	return goCRD, nil
}

func in[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
