// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"testing"
)

func TestRegisteredGenerators(t *testing.T) {
	// Verify that all expected generators are registered via init()
	expectedGenerators := []string{
		AtlasControllersGeneratorName,
		AtlasExportersGeneratorName,
		IndexersGeneratorName,
	}

	registeredNames := List()

	if len(registeredNames) != len(expectedGenerators) {
		t.Errorf("expected %d registered generators, got %d: %v", len(expectedGenerators), len(registeredNames), registeredNames)
	}

	for _, expected := range expectedGenerators {
		g := Get(expected)
		if g == nil {
			t.Errorf("expected generator %q to be registered", expected)
			continue
		}
		if g.Name() != expected {
			t.Errorf("generator name mismatch: expected %q, got %q", expected, g.Name())
		}
		if g.Description() == "" {
			t.Errorf("generator %q has empty description", expected)
		}
	}
}

func TestGetByNames(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		wantCount   int
		wantErr     bool
		errContains string
	}{
		{
			name:      "single valid generator",
			input:     []string{IndexersGeneratorName},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "multiple valid generators",
			input:     []string{AtlasControllersGeneratorName, IndexersGeneratorName},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "all generators",
			input:     []string{AtlasControllersGeneratorName, AtlasExportersGeneratorName, IndexersGeneratorName},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:        "unknown generator",
			input:       []string{"unknown-generator"},
			wantErr:     true,
			errContains: "unknown generator",
		},
		{
			name:        "mix of valid and invalid",
			input:       []string{IndexersGeneratorName, "invalid"},
			wantErr:     true,
			errContains: "unknown generator",
		},
		{
			name:      "empty list",
			input:     []string{},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gens, err := GetByNames(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(gens) != tt.wantCount {
				t.Errorf("expected %d generators, got %d", tt.wantCount, len(gens))
			}
		})
	}
}

func TestList(t *testing.T) {
	names := List()

	// Verify the list is sorted
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Errorf("List() is not sorted: %v comes before %v", names[i-1], names[i])
		}
	}

	// Verify all registered generators are in the list
	for _, name := range names {
		if Get(name) == nil {
			t.Errorf("List() returned %q but Get(%q) returns nil", name, name)
		}
	}
}

func TestAll(t *testing.T) {
	gens := All()
	names := List()

	if len(gens) != len(names) {
		t.Errorf("All() returned %d generators but List() returned %d names", len(gens), len(names))
	}
}

func TestGeneratorDescriptions(t *testing.T) {
	gens := All()

	for _, gen := range gens {
		t.Run(gen.Name(), func(t *testing.T) {
			desc := gen.Description()
			if desc == "" {
				t.Error("Description() returned empty string")
			}
			if len(desc) < 10 {
				t.Errorf("Description() is too short: %q", desc)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
