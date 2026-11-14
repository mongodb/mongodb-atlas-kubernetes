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
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cliFakeDataSource struct {
	release    string
	changes    []Change
	input      Change
	lastInput  string
	releaseErr error
	changesErr error
}

func (f *cliFakeDataSource) LatestRelease(context.Context) (string, error) {
	return f.release, f.releaseErr
}
func (f *cliFakeDataSource) ChangesSince(context.Context, string) ([]Change, error) {
	return f.changes, f.changesErr
}
func (f *cliFakeDataSource) ChangeForInput(_ context.Context, input string) (Change, error) {
	f.lastInput = input
	return f.input, nil
}

func TestRunNoArgumentsPrintsVersionAndUsage(t *testing.T) {
	var out, stderr bytes.Buffer
	err := Run(context.Background(), RuntimeInputs{Args: []string{}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: &cliFakeDataSource{}})
	require.NoError(t, err)
	stdout := out.String()
	assert.Contains(t, stdout, Version)
	for _, want := range []string{"Usage", "version", "released", "latest", "list", "validate <input|->", "reads stdin", "verify <input|->"} {
		assert.Contains(t, stdout, want)
	}
	assert.Empty(t, stderr.String())
}

func TestRunVersionCommand(t *testing.T) {
	var out, stderr bytes.Buffer
	source := &cliFakeDataSource{
		release: "v2.18.1",
		changes: []Change{{Type: ChangeTypePatch, Title: "[patch] Fix a bug"}},
	}
	err := Run(context.Background(), RuntimeInputs{Args: []string{"version"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: source})
	require.NoError(t, err)
	assert.Equal(t, "2.18.2\n", out.String())
}

func TestRunReleasedAndLatestAreAliases(t *testing.T) {
	for _, command := range []string{"released", "latest"} {
		t.Run(command, func(t *testing.T) {
			var out, stderr bytes.Buffer
			err := Run(context.Background(), RuntimeInputs{Args: []string{command}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: &cliFakeDataSource{release: "v3.2.1"}})
			require.NoError(t, err)
			assert.Equal(t, "3.2.1\n", out.String())
		})
	}
}

func TestRunListUsesChangesAfterLatestAndOmitsSkip(t *testing.T) {
	ds := &cliFakeDataSource{release: "v1.0.0", changes: []Change{
		{ID: "abc123", Type: ChangeTypeFeature, Title: "[feature] Keep this", Description: "Feature detail"},
		{Type: ChangeTypeSkipCL, Title: "[skip-cl] Omit this", Description: "Generated"},
		{ID: "def456", Type: ChangeTypePatch, Title: "[patch] Keep this without a description"},
		{ID: "ghi789", Type: ChangeTypeUnset, Title: "Unclassified but allowed"},
	}}
	var out, stderr bytes.Buffer
	err := Run(context.Background(), RuntimeInputs{Args: []string{"list"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds})
	var unsetErr *UnsetChangesError
	require.ErrorAs(t, err, &unsetErr)
	assert.EqualError(t, err, "found 1 unset changelog entry")
	assert.Equal(t, "abc123 feature [feature] Keep this\nFeature detail\n\ndef456 patch [patch] Keep this without a description\n\nghi789 unset Unclassified but allowed\n\n", out.String())
}

func TestRunListSourceErrorsReturnBeforeWriting(t *testing.T) {
	want := errors.New("source failed")
	for _, ds := range []*cliFakeDataSource{
		{releaseErr: want},
		{release: "v1.0.0", changes: []Change{{Title: "Should not print"}}, changesErr: want},
	} {
		var out, stderr bytes.Buffer
		err := Run(context.Background(), RuntimeInputs{Args: []string{"list"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds})
		assert.ErrorIs(t, err, want)
		assert.Empty(t, out.String())
	}
}

func TestRunValidateAndStdlibFlagSemantics(t *testing.T) {
	ds := &cliFakeDataSource{input: Change{Type: ChangeTypePatch, Title: "[patch] valid"}}
	var out, stderr bytes.Buffer
	require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"validate", "commit:abc"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}))
	assert.Empty(t, out.String())
	require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"verify", "HEAD~1..HEAD"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}))
	assert.Empty(t, out.String())
	assert.Equal(t, "HEAD~1..HEAD", ds.lastInput)
	assert.Error(t, Run(context.Background(), RuntimeInputs{Args: []string{"--definitely-unknown"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}), "Run() accepted an unknown standard-library flag")
	assert.Error(t, Run(context.Background(), RuntimeInputs{Args: []string{"validate"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}), "Run() accepted a command missing its required input")
	assert.EqualError(t, Run(context.Background(), RuntimeInputs{Args: []string{"verify"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}), "verify requires exactly one input (got 0)")
}

func TestRunReturnsDatasourceErrors(t *testing.T) {
	want := errors.New("source failed")
	ds := errorCLIDataSource{err: want}
	var out, stderr bytes.Buffer
	assert.ErrorIs(t, Run(context.Background(), RuntimeInputs{Args: []string{"latest"}, Stdin: bytes.NewReader(nil), Stdout: &out, Stderr: &stderr, DataSource: ds}), want)
}

func TestRunValidateReadsTitleAndDescriptionFromStdin(t *testing.T) {
	for _, input := range []string{
		"A title\n[feature] details\n",
		"[patch] A title\nDetails without another marker\n",
	} {
		t.Run(strings.ReplaceAll(input, "\n", "/"), func(t *testing.T) {
			var out, stderr bytes.Buffer
			ds := &cliFakeDataSource{}
			require.NoError(t, Run(context.Background(), RuntimeInputs{Args: []string{"validate", "-"}, Stdin: strings.NewReader(input), Stdout: &out, Stderr: &stderr, DataSource: ds}))
			assert.Empty(t, out.String())
			assert.Empty(t, stderr.String())
			assert.Empty(t, ds.lastInput, "stdin validation must not resolve a datasource input")
		})
	}
}

func TestRunValidateStdinRequiresExactlyOneMarkerAndShowsGuidance(t *testing.T) {
	for _, test := range []struct {
		name       string
		input      string
		wantReason string
	}{
		{name: "missing", input: "A title\nNo marker here", wantReason: "Change entry is missing one marker of"},
		{name: "duplicate across title and body", input: "[feature] A title\n[patch] details", wantReason: "Change entry has multiple markers; exactly one is allowed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := Run(context.Background(), RuntimeInputs{Args: []string{"validate", "-"}, Stdin: strings.NewReader(test.input), Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, DataSource: &cliFakeDataSource{}})
			require.Error(t, err)
			want := "Rejected changelog entry:\n" + strings.TrimSpace(test.input) + "\n\n" + test.wantReason + "\n" + markerGuidance
			assert.EqualError(t, err, want)
			assert.NotContains(t, err.Error(), "validate input")
			assert.NotContains(t, err.Error(), "validate stdin")
		})
	}
}

func TestRunValidateStdinReportsReaderErrors(t *testing.T) {
	want := errors.New("read failed")
	assert.EqualError(t, Run(context.Background(), RuntimeInputs{Args: []string{"validate", "-"}, Stdin: errorReader{err: want}, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, DataSource: &cliFakeDataSource{}}), "validate stdin: read input: read failed")
	assert.EqualError(t, Run(context.Background(), RuntimeInputs{Args: []string{"validate", "-"}, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}, DataSource: &cliFakeDataSource{}}), "validate stdin: input reader is nil")
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }

type errorCLIDataSource struct{ err error }

func (e errorCLIDataSource) LatestRelease(context.Context) (string, error) { return "", e.err }
func (e errorCLIDataSource) ChangesSince(context.Context, string) ([]Change, error) {
	return nil, e.err
}
