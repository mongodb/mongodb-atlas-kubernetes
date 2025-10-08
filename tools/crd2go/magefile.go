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

//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// CI runs all linting and validation checks.
func CI() {
	mg.Deps(Build, UnitTests, Addlicense, Checklicense, CheckGCI, Lint, Govulncheck)
	fmt.Println("✅ CI PASSED all checks")
}

// Build checks all execitable build properly
func Build() error {
	return wrapRun("🛠️  Build", "go", "build", "./...")
}

// UnitTests runs the go tests
func UnitTests() error {
	return wrapRun("🧪 Unit tests", "go", "test", "-cover", "./...")
}

// Addlicense runs the addlicense check to ensure source files have license headers
func Addlicense() error {
	return wrapRun("🛠️  License header check",
		"go", "tool",
		"addlicense",
		"-check",
		"-l", "apache",
		"-c", "MongoDB Inc",
		"-ignore", "**/*.md",
		"-ignore", "**/*.yaml",
		"-ignore", "**/*.yml",
		"-ignore", "**/*.nix",
		"-ignore", ".devbox/**",
		"-ignore", "internal/samples/**",
		"-ignore", "magefile.go",
		".",
	)
}

// Checklicense runs the go-licenses tool to check license compliance
func Checklicense() error {
	return wrapRun("🔬 License compliance check",
		"go", "tool",
		"go-licenses", "check",
		"--include_tests",
		"--disallowed_types", "restricted,reciprocal",
		"./...",
	)
}

// GCI runs gci to check that Go import orders are correct
func GCI() error {
	return wrapRun("🧹 Format code",
		"go", "tool",
		"gci", "write",
		"--skip-generated",
		"-s", "standard",
		"-s", "default",
		"-s", "localmodule",
		".",
	)
}

// GitClean check git is clean of changes
func GitClean() error {
	if err := sh.Run("git", "diff-index", "--quiet", "HEAD", "--"); err != nil {
		fmt.Println("❗️ The following files have changes:")
		sh.RunV("git", "diff-index", "--name-only", "HEAD")
		return fmt.Errorf("please run 'mage gci' and commit any changes")
	}
	return nil
}

// CheckGCI check GCI formatting was committed as expected
func CheckGCI() {
	mg.SerialDeps(GCI, GitClean)
}

// Lint runs the golangci-lint tool
func Lint() error {
	if err := os.Setenv("CGO_ENABLED", "0"); err != nil {
		return nil
	}
	return wrapRun("🔬 Lint check",
		"go", "tool", "golangci-lint", "run",
		"./cmd/...", "./internal/...", "./k8s/...", "./pkg/...")
}

// Govulncheck checks for Go toolchain or library vulnerabilities
func Govulncheck() error {
	return wrapRun("🔬 Vulnerability Check",
		"go", "tool", "govulncheck", "./...")
}

func wrapRun(action, cmd string, args ...string) error {
	fmt.Printf("▶️ Started: %s\n", action)
	if err := sh.RunV(cmd, args...); err != nil {
		return err
	}
	fmt.Printf("✅ Succeeded: %s\n", action)
	return nil
}
