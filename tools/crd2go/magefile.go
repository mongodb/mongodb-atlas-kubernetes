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
	mg.SerialDeps(Build, UnitTests, Addlicense, Checklicense, GCI, Lint)
	fmt.Println("✅ CI PASSED all checks")
}

// Build checks all execitable build properly
func Build() error {
	return wrapRun("🛠️  Building...", "go", "build", "./...")
}

// UnitTests runs the go tests
func UnitTests() error {
	return wrapRun("🧪 Running unit tests:\n", "go", "test", "-cover", "./...")
}

// Addlicense runs the addlicense check to ensure source files have license headers
func Addlicense() error {
	return wrapRun("🛠️  Running license header check...",
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
	return wrapRun("🔬 Running license compliance checks:\n",
		"go", "tool",
		"go-licenses", "check",
		"--include_tests",
		"--disallowed_types", "restricted,reciprocal",
		"./...",
	)
}

// GCI runs gci to check that Go import orders are correct
func GCI() error {
	fmt.Println("🧹 Formatting Go imports...")
	if err := sh.RunV(
		"go", "tool",
		"gci", "write",
		"--skip-generated",
		"-s", "standard",
		"-s", "default",
		"-s", "localmodule",
		".",
	); err != nil {
		return fmt.Errorf("gci write command failed: %w", err)
	}

	fmt.Println("🔍 Checking for changes...")
	if err := sh.Run("git", "diff-index", "--quiet", "HEAD", "--"); err != nil {
		fmt.Println("❗️ Go files were not correctly formatted. The following files have changes:")
		sh.RunV("git", "diff-index", "--name-only", "HEAD")
		return fmt.Errorf("please run 'mage gci' and commit the changes")
	}

	fmt.Println("✅ Go imports are correctly formatted.")
	return nil
}

// Lint runs the golangci-lint tool
func Lint() error {
	if err := os.Setenv("CGO_ENABLED", "0"); err != nil {
		return nil
	}
	return wrapRun("▶️ Run linting...",
		"go", "tool", "golangci-lint", "run",
		"./cmd/...", "./internal/...", "./k8s/...", "./pkg/...")
}

func wrapRun(msg, cmd string, args ...string) error {
	fmt.Print(msg)
	if err := sh.RunV(cmd, args...); err != nil {
		fmt.Println()
		return err
	}
	fmt.Println("✅ Success")
	return nil
}
