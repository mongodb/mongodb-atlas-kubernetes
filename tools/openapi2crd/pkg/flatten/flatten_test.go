package flatten

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "regenerate golden files")

// TestFlattenGolden runs Flatten on every testdata/<name>/input.yaml and
// compares the output against testdata/<name>/golden.yaml.
// Run with -update to regenerate the golden files.
func TestFlattenGolden(t *testing.T) {
	inputs, err := filepath.Glob("testdata/*/input.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(inputs) == 0 {
		t.Fatal("no test fixtures found under testdata/")
	}

	for _, inputPath := range inputs {
		dir := filepath.Dir(inputPath)
		name := filepath.Base(dir)
		goldenPath := filepath.Join(dir, "golden.yaml")

		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("reading input: %v", err)
			}

			got, err := Flatten(data)
			if err != nil {
				t.Fatalf("Flatten: %v", err)
			}

			if *update {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("writing golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("reading golden (run with -update to create): %v", err)
			}

			if string(got) != string(want) {
				t.Errorf("output mismatch for %s:\n\ngot:\n%s\nwant:\n%s", name, got, want)
			}
		})
	}
}
