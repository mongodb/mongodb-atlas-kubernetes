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

package clbot

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test builds a real local repository from the checked-in sandbox seed.
func TestGitDataSourceReadsLatestTagAndCommitText(t *testing.T) {
	seed := filepath.Join("testdata", "git-sandbox", "seed.txt")
	seedContents, err := os.ReadFile(seed) // #nosec G304 -- seed is a fixed checked-in fixture path.
	require.NoError(t, err)
	repo := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repo, "seed.txt"), seedContents, 0o600)) // #nosec G703 -- output joins a test-generated temporary repo with a constant filename.
	git(t, repo, "init", "-q")
	git(t, repo, "config", "user.name", "CL Integration Test")
	git(t, repo, "config", "user.email", "cl-test@example.invalid")
	git(t, repo, "add", "seed.txt")
	git(t, repo, "commit", "-q", "-m", "Initial release")
	git(t, repo, "tag", "v1.2.3")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[feature] Add local sandbox support", "-m", "Detailed description from the commit body.")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[bugfix] Correct sandbox metadata")
	git(t, repo, "tag", "v1.2.4")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[patch] Follow-up after release", "-m", "This commit is after the latest tag.")

	source := NewGitDataSource(repo)
	gotRelease, err := source.LatestRelease(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "v1.2.4", gotRelease)
	gotChanges, err := source.ChangesSince(context.Background(), gotRelease)
	require.NoError(t, err)
	patchSHA, err := exec.CommandContext(t.Context(), "git", "-C", repo, "rev-parse", "HEAD").Output() // #nosec G204 -- git arguments are fixed by this test.
	require.NoError(t, err)
	want := []Change{{ID: strings.TrimSpace(string(patchSHA)), Type: ChangeTypePatch, Title: "[patch] Follow-up after release", Description: "This commit is after the latest tag."}}
	require.Len(t, gotChanges, len(want))
	assert.Equal(t, want[0], gotChanges[0])

	gotHead, err := source.ChangeForInput(context.Background(), "HEAD")
	require.NoError(t, err)
	assert.Equal(t, want[0], gotHead)

	var stdout, stderr bytes.Buffer
	require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"latest"}, Stdin: bytes.NewReader(nil), Stdout: &stdout, Stderr: &stderr, DataSource: source}))
	assert.Equal(t, "1.2.4\n", stdout.String())
	assert.Empty(t, stderr.String())

	stdout.Reset()
	stderr.Reset()
	require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"version"}, Stdin: bytes.NewReader(nil), Stdout: &stdout, Stderr: &stderr, DataSource: source}))
	assert.Equal(t, "1.2.5\n", stdout.String())
	assert.Empty(t, stderr.String())

	stdout.Reset()
	stderr.Reset()
	require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"verify", "HEAD"}, Stdin: bytes.NewReader(nil), Stdout: &stdout, Stderr: &stderr, DataSource: source}))
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())

	featureCommit, err := exec.CommandContext(t.Context(), "git", "-C", repo, "rev-parse", "HEAD~2").Output() // #nosec G204 -- git arguments are fixed by this test.
	require.NoError(t, err)
	featureInput := "commit:" + strings.TrimSpace(string(featureCommit))[:8]
	gotFeature, err := source.ChangeForInput(context.Background(), featureInput)
	require.NoError(t, err)
	assert.Equal(t, Change{ID: strings.TrimSpace(string(featureCommit)), Type: ChangeTypeFeature, Title: "[feature] Add local sandbox support", Description: "Detailed description from the commit body."}, gotFeature)
	require.NoError(t, Validate(context.Background(), source, featureInput))
}

func TestGitDataSourceLatestReleaseIgnoresPrereleaseTags(t *testing.T) {
	repo := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repo, "seed.txt"), []byte("seed\n"), 0o600)) // #nosec G703 -- output joins a test-generated temporary repo with a constant filename.
	git(t, repo, "init", "-q")
	git(t, repo, "config", "user.name", "CL Integration Test")
	git(t, repo, "config", "user.email", "cl-test@example.invalid")
	git(t, repo, "add", "seed.txt")
	git(t, repo, "commit", "-q", "-m", "Initial release")
	git(t, repo, "tag", "v1.2.3")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[feature] Add thing")
	git(t, repo, "tag", "v1.2.4")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[patch] Follow-up after release")
	git(t, repo, "tag", "v1.2.5-pre-release")
	git(t, repo, "commit", "--allow-empty", "-q", "-m", "[patch] Later patch")

	source := NewGitDataSource(repo)
	gotRelease, err := source.LatestRelease(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "v1.2.4", gotRelease)

	gotVersion, err := NextVersion(context.Background(), source)
	require.NoError(t, err)
	assert.Equal(t, "1.2.5", gotVersion)
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "git", append([]string{"-C", dir}, args...)...) // #nosec G204 -- git subprocess arguments are authored by this test.
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %s: %s", strings.Join(args, " "), output)
}
