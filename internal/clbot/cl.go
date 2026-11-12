// Copyright 2026 MongoDB Inc
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

// Package clbot contains the changelog logic used by the clbot command.
package clbot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Version is the CLI version. Release builds may replace this value at build time.
const Version = "dev"

// Change is a change supplied by a datasource.
type Change struct {
	ID          string
	Type        ChangeType
	Title       string
	Description string
}

// ChangeType classifies a changelog entry.
type ChangeType string

const (
	ChangeTypeSkipCL   ChangeType = "skip-cl"
	ChangeTypePatch    ChangeType = "patch"
	ChangeTypeFeature  ChangeType = "feature"
	ChangeTypeBreaking ChangeType = "breaking"
	ChangeTypeUnset    ChangeType = "unset"
)

// DataSource supplies releases and changes.
type DataSource interface {
	LatestRelease(ctx context.Context) (string, error)
	ChangesSince(ctx context.Context, release string) ([]Change, error)
}

// InputResolver resolves datasource-specific input for validation. Git
// implementations interpret input as a commit; file implementations may
// interpret it as a file path.
type InputResolver interface {
	ChangeForInput(ctx context.Context, input string) (Change, error)
}

// GitDataSource reads releases and changes from a Git repository.
type GitDataSource struct {
	Repository string
}

// NewGitDataSource constructs the default Git datasource for repository.
func NewGitDataSource(repository string) *GitDataSource {
	return &GitDataSource{Repository: repository}
}

func (g *GitDataSource) git(ctx context.Context, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-C", g.Repository}, args...)
	output, err := exec.CommandContext(ctx, "git", commandArgs...).Output()
	if err != nil {
		var stderr string
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr = strings.TrimSpace(string(exitErr.Stderr))
		}
		if stderr != "" {
			return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, stderr)
		}
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return output, nil
}

func (g *GitDataSource) LatestRelease(ctx context.Context) (string, error) {
	output, err := g.git(ctx, "describe", "--tags", "--abbrev=0", "HEAD")
	if err != nil {
		return "", fmt.Errorf("find latest reachable release tag: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *GitDataSource) ChangesSince(ctx context.Context, release string) ([]Change, error) {
	if strings.TrimSpace(release) == "" {
		return nil, fmt.Errorf("list changes since release: release is empty")
	}
	output, err := g.git(ctx, "log", "--format=%x1e%H%x00%s%x00%b", release+"..HEAD")
	if err != nil {
		return nil, fmt.Errorf("list changes since release %q: %w", release, err)
	}
	records := bytes.Split(output, []byte{0x1e})
	changes := make([]Change, 0, len(records)-1)
	for _, record := range records[1:] {
		fields := bytes.SplitN(record, []byte{0}, 3)
		if len(fields) != 3 {
			continue
		}
		id := strings.TrimSpace(string(fields[0]))
		title := strings.TrimSpace(string(fields[1]))
		body := strings.TrimSpace(string(fields[2]))
		if title == "" {
			continue
		}
		changes = append(changes, Change{ID: id, Type: changeType(title, body), Title: title, Description: body})
	}
	return changes, nil
}

func (g *GitDataSource) ChangeForInput(ctx context.Context, input string) (Change, error) {
	const prefix = "commit:"
	revision := input
	if after, ok := strings.CutPrefix(input, prefix); ok {
		revision = after
		if !regexp.MustCompile(`(?i)^[0-9a-f]{4,40}$`).MatchString(revision) {
			return Change{}, fmt.Errorf("resolve change input %q: malformed commit SHA", input)
		}
	}
	resolved, err := g.git(ctx, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		return Change{}, fmt.Errorf("resolve change input %q: %w", input, err)
	}
	output, err := g.git(ctx, "show", "-s", "--format=%s%x00%b%x00", strings.TrimSpace(string(resolved)))
	if err != nil {
		return Change{}, fmt.Errorf("resolve change input %q: %w", input, err)
	}
	fields := bytes.Split(output, []byte{0})
	if len(fields) < 2 || strings.TrimSpace(string(fields[0])) == "" {
		return Change{}, fmt.Errorf("resolve change input %q: commit has no title", input)
	}
	return Change{
		ID:          strings.TrimSpace(string(resolved)),
		Type:        changeType(strings.TrimSpace(string(fields[0])), strings.TrimSpace(string(fields[1]))),
		Title:       strings.TrimSpace(string(fields[0])),
		Description: strings.TrimSpace(string(fields[1])),
	}, nil
}

// UnsetChangesError reports the number of included changes that do not have a
// recognized changelog type. Changes are still returned so callers can display
// them alongside this validation error.
type UnsetChangesError struct {
	Count int
}

func (e *UnsetChangesError) Error() string {
	entry := "entries"
	if e.Count == 1 {
		entry = "entry"
	}
	return fmt.Sprintf("found %d unset changelog %s", e.Count, entry)
}

// List returns the non-skipped changes since the datasource's latest release.
func List(ctx context.Context, source DataSource) ([]Change, error) {
	release, err := source.LatestRelease(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest release: %w", err)
	}
	changes, err := source.ChangesSince(ctx, release)
	if err != nil {
		return nil, fmt.Errorf("get changes since latest release %q: %w", release, err)
	}
	included := make([]Change, 0, len(changes))
	unsetCount := 0
	for _, change := range changes {
		if change.Type == ChangeTypeSkipCL {
			continue
		}
		included = append(included, change)
		if change.Type == ChangeTypeUnset {
			unsetCount++
		}
	}
	if unsetCount > 0 {
		return included, &UnsetChangesError{Count: unsetCount}
	}
	return included, nil
}

// Validate checks that the supplied datasource-specific input has one marker.
func Validate(ctx context.Context, source DataSource, input string) error {
	resolver, ok := source.(InputResolver)
	if !ok {
		return fmt.Errorf("validate input %q: datasource does not support input resolution", input)
	}
	change, err := resolver.ChangeForInput(ctx, input)
	if err != nil {
		return fmt.Errorf("validate input %q: %w", input, err)
	}
	return validateChange(change)
}

const markerGuidance = "[skip-cl] use this for CI changes or other non-production changes\n" +
	"[bugfix] or [patch] use for patches or bugfixes\n" +
	"[feature] use for backward compatible added functionality\n" +
	"[breaking] use for breaking changes"

func validateChange(change Change) error {
	markers := regexp.MustCompile(`\[(skip-cl|bugfix|patch|feature|breaking)\]`)
	count := len(markers.FindAllString(change.Title+"\n"+change.Description, -1))
	content := strings.Join(nonEmpty(change.Title, change.Description), "\n")
	if count == 0 {
		return fmt.Errorf("Rejected changelog entry:\n%s\n\nChange entry is missing one marker of\n%s", content, markerGuidance)
	}
	if count > 1 {
		return fmt.Errorf("Rejected changelog entry:\n%s\n\nChange entry has multiple markers; exactly one is allowed\n%s", content, markerGuidance)
	}
	switch change.Type {
	case ChangeTypeSkipCL, ChangeTypePatch, ChangeTypeFeature, ChangeTypeBreaking:
		return nil
	default:
		return fmt.Errorf("Rejected changelog entry:\n%s\n\nChange entry has unexpected type %q", content, change.Type)
	}
}

var changeMarkers = regexp.MustCompile(`\[(skip-cl|bugfix|patch|feature|breaking)\]`)

func changeType(title, description string) ChangeType {
	markers := changeMarkers.FindAllStringSubmatch(title+"\n"+description, -1)
	if len(markers) != 1 {
		return ChangeTypeUnset
	}
	switch markers[0][1] {
	case "skip-cl":
		return ChangeTypeSkipCL
	case "bugfix", "patch":
		return ChangeTypePatch
	case "feature":
		return ChangeTypeFeature
	case "breaking":
		return ChangeTypeBreaking
	default:
		return ChangeTypeUnset
	}
}

func nonEmpty(parts ...string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// NextVersion calculates the next semver version from included changes.
func NextVersion(ctx context.Context, source DataSource) (string, error) {
	release, err := source.LatestRelease(ctx)
	if err != nil {
		return "", fmt.Errorf("get latest release: %w", err)
	}
	changes, err := source.ChangesSince(ctx, release)
	if err != nil {
		return "", fmt.Errorf("get changes since latest release %q: %w", release, err)
	}
	included := make([]Change, 0, len(changes))
	for _, change := range changes {
		if change.Type != ChangeTypeSkipCL {
			included = append(included, change)
		}
	}

	version := strings.TrimPrefix(release, "v")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid release version %q: expected X.Y.Z", release)
	}
	values := [3]int{}
	for i, part := range parts {
		value, parseErr := strconv.Atoi(part)
		if parseErr != nil || value < 0 || strconv.Itoa(value) != part {
			return "", fmt.Errorf("invalid release version %q: expected non-negative X.Y.Z", release)
		}
		values[i] = value
	}

	breaking, feature := false, false
	for _, change := range included {
		if change.Type == ChangeTypeBreaking {
			breaking = true
		} else if change.Type == ChangeTypeFeature {
			feature = true
		}
	}
	switch {
	case breaking:
		values[0]++
		values[1], values[2] = 0, 0
	case feature:
		values[1]++
		values[2] = 0
	default:
		values[2]++
	}
	return fmt.Sprintf("%d.%d.%d", values[0], values[1], values[2]), nil
}
