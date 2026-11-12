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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGit struct {
	release    string
	changes    []Change
	input      Change
	inputSeen  string
	releaseErr error
	changesErr error
	inputErr   error
}

func (f *fakeGit) LatestRelease(context.Context) (string, error) { return f.release, f.releaseErr }
func (f *fakeGit) ChangesSince(_ context.Context, release string) ([]Change, error) {
	if release != f.release {
		return nil, errors.New("queried changes from unexpected release")
	}
	return f.changes, f.changesErr
}
func (f *fakeGit) ChangeForInput(_ context.Context, input string) (Change, error) {
	f.inputSeen = input
	return f.input, f.inputErr
}

type fakeFileDataSource struct{ fakeGit }

func TestDatasourceIsInjectableAndListUsesLatestRelease(t *testing.T) {
	ds := &fakeGit{release: "v1.4.2", changes: []Change{
		{Type: ChangeTypeFeature, Title: "[feature] Add feature", Description: "Details"},
		{Type: ChangeTypeSkipCL, Title: "[skip-cl] generated update", Description: "Do not include"},
		{Type: ChangeTypePatch, Title: "[patch] Fix issue", Description: ""},
	}}
	got, err := List(context.Background(), ds)
	require.NoError(t, err)
	want := []Change{ds.changes[0], ds.changes[2]}
	assert.Equal(t, want, got)
}

func TestValidateRequiresExactlyOneRecognizedMarkerOnTitleOrDescription(t *testing.T) {
	tests := []struct {
		name, title, description string
		typeValue                ChangeType
		valid                    bool
	}{
		{"title marker", "[feature] add thing", "details", ChangeTypeFeature, true},
		{"description marker", "add thing", "details [bugfix]", ChangeTypePatch, true},
		{"patch alias", "[patch] fix thing", "", ChangeTypePatch, true},
		{"skip marker", "[skip-cl] generated", "", ChangeTypeSkipCL, true},
		{"missing marker", "ordinary title", "ordinary description", ChangeTypeUnset, false},
		{"unknown marker", "[docs] update docs", "", ChangeTypeUnset, false},
		{"duplicate in title", "[feature] [patch] ambiguous", "", ChangeTypeUnset, false},
		{"duplicate across fields", "[feature] title", "[breaking] description", ChangeTypeUnset, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &fakeGit{input: Change{Type: tt.typeValue, Title: tt.title, Description: tt.description}}
			err := Validate(context.Background(), ds, "commit:abc123")
			assert.Equal(t, tt.valid, err == nil)
			assert.Equal(t, "commit:abc123", ds.inputSeen, "datasource input should match exact commit input")
		})
	}

	guidance := markerGuidance

	missing := &fakeGit{input: Change{Type: ChangeTypeUnset, Title: "ordinary title", Description: "ordinary description"}}
	require.EqualError(t, Validate(context.Background(), missing, "commit:abc123"),
		"Rejected changelog entry:\nordinary title\nordinary description\n\nChange entry is missing one marker of\n"+guidance)

	duplicate := &fakeGit{input: Change{Type: ChangeTypeUnset, Title: "[feature] [patch] ambiguous", Description: "details"}}
	require.EqualError(t, Validate(context.Background(), duplicate, "commit:abc123"),
		"Rejected changelog entry:\n[feature] [patch] ambiguous\ndetails\n\nChange entry has multiple markers; exactly one is allowed\n"+guidance)
}

func TestValidateDelegatesDatasourceSpecificCommitAndFileInputs(t *testing.T) {
	for _, input := range []string{"commit:deadbeef", "file:changes/feature.md"} {
		t.Run(input, func(t *testing.T) {
			ds := &fakeFileDataSource{fakeGit: fakeGit{input: Change{Type: ChangeTypeFeature, Title: "[feature] valid"}}}
			require.NoError(t, Validate(context.Background(), ds, input))
			assert.Equal(t, input, ds.inputSeen)
		})
	}
}

func TestValidateRejectsDatasourceUnsetTypeEvenWhenTitleHasMarker(t *testing.T) {
	ds := &fakeGit{input: Change{Type: ChangeTypeUnset, Title: "[feature] marker in title"}}
	err := Validate(context.Background(), ds, "commit:abc123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `Change entry has unexpected type "unset"`)
}

func TestValidateRejectsUnexpectedDatasourceType(t *testing.T) {
	ds := &fakeGit{input: Change{Type: ChangeType("documentation"), Title: "[feature] marker in title"}}
	err := Validate(context.Background(), ds, "commit:abc123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `Change entry has unexpected type "documentation"`)
}

func TestNextVersionUsesMarkerPrecedenceAndSemverResets(t *testing.T) {
	tests := []struct {
		name, current, want string
		changes             []Change
	}{
		{"breaking dominates feature and patch", "v1.4.9", "2.0.0", []Change{{Type: ChangeTypePatch, Title: "[patch] p"}, {Type: ChangeTypeFeature, Title: "[feature] f"}, {Type: ChangeTypeBreaking, Title: "[breaking] b"}}},
		{"feature dominates patch", "1.4.9", "1.5.0", []Change{{Type: ChangeTypePatch, Title: "[patch] p"}, {Type: ChangeTypeFeature, Title: "[feature] f"}}},
		{"patch alias", "1.4.9", "1.4.10", []Change{{Type: ChangeTypePatch, Title: "[bugfix] b"}}},
		{"patch marker", "1.4.9", "1.4.10", []Change{{Type: ChangeTypePatch, Title: "[patch] p"}}},
		{"default patch for empty changes", "1.4.9", "1.4.10", nil},
		{"default patch for unset changes", "1.4.9", "1.4.10", []Change{{Type: ChangeTypeUnset, Title: "ordinary entry"}}},
		{"skipped changes do not affect bump", "1.4.9", "1.4.10", []Change{{Type: ChangeTypeSkipCL, Title: "[skip-cl] ignored"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &fakeGit{release: tt.current, changes: tt.changes}
			got, err := NextVersion(context.Background(), ds)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDatasourceErrorsAreReturned(t *testing.T) {
	wantErr := errors.New("datasource unavailable")
	tests := []struct {
		name string
		ds   *fakeGit
		run  func(*fakeGit) error
	}{
		{"latest release in list", &fakeGit{releaseErr: wantErr}, func(ds *fakeGit) error { _, err := List(context.Background(), ds); return err }},
		{"changes in list", &fakeGit{release: "v1.0.0", changesErr: wantErr}, func(ds *fakeGit) error { _, err := List(context.Background(), ds); return err }},
		{"latest release in version", &fakeGit{releaseErr: wantErr}, func(ds *fakeGit) error { _, err := NextVersion(context.Background(), ds); return err }},
		{"changes in version", &fakeGit{release: "v1.0.0", changesErr: wantErr}, func(ds *fakeGit) error { _, err := NextVersion(context.Background(), ds); return err }},
		{"input lookup in validate", &fakeGit{inputErr: wantErr}, func(ds *fakeGit) error { return Validate(context.Background(), ds, "commit:sha") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, tt.run(tt.ds), wantErr)
		})
	}
}

func TestChangeTypeMapsRecognizedMarkers(t *testing.T) {
	for _, test := range []struct {
		marker string
		want   ChangeType
	}{
		{"[skip-cl]", ChangeTypeSkipCL},
		{"[bugfix]", ChangeTypePatch},
		{"[patch]", ChangeTypePatch},
		{"[feature]", ChangeTypeFeature},
		{"[breaking]", ChangeTypeBreaking},
	} {
		t.Run(test.marker, func(t *testing.T) {
			assert.Equal(t, test.want, changeType("title "+test.marker, "details"))
		})
	}
	assert.Equal(t, ChangeTypeUnset, changeType("ordinary", "[docs]"))
	assert.Equal(t, ChangeTypeUnset, changeType("[feature] [patch]", ""))
}

func TestListIncludesUnsetChanges(t *testing.T) {
	withID := Change{ID: "abc123", Type: ChangeTypeUnset, Title: "unclassified change"}
	withoutID := Change{Type: ChangeTypeUnset, Title: "another unclassified change"}
	skipped := Change{ID: "skip", Type: ChangeTypeSkipCL, Title: "[skip-cl] generated"}
	ds := &fakeGit{release: "v1.0.0", changes: []Change{withID, skipped, withoutID}}
	got, err := List(context.Background(), ds)
	require.Error(t, err)
	var unsetErr *UnsetChangesError
	require.ErrorAs(t, err, &unsetErr)
	assert.Equal(t, []Change{withID, withoutID}, got)
	assert.Equal(t, 2, unsetErr.Count)
	assert.EqualError(t, err, "found 2 unset changelog entries")
}

func TestUnsetChangesErrorSingular(t *testing.T) {
	assert.EqualError(t, (&UnsetChangesError{Count: 1}), "found 1 unset changelog entry")
}
