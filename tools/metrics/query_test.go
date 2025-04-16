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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/require"
)

const (
	// unit true means unit tests so use the playback test client
	// when false the recorder is used instead and it behaves as an end to end test
	unit = true

	recorderDir = "samples"
)

var lastRecordingTime = time.Date(2024, 01, 03, 19, 01, 0, 0, time.UTC)

func TestQuery(t *testing.T) {
	qc := newTestClient()
	wfRuns, err := qc.TestWorkflowRuns("main", "push", 1)
	require.NoError(t, err)
	require.NotNil(t, wfRuns)
}

func newTestClient() QueryClient {
	if unit {
		return newPlaybackGHQueryClient(recorderDir)
	}
	return newRecorderGHQueryClient(recorderDir)
}

type recorderGHQueryClient struct {
	ghQueryClient
	path string
}

func newRecorderGHQueryClient(path string) *recorderGHQueryClient {
	return &recorderGHQueryClient{ghQueryClient: *NewDefaultQueryClient(), path: path}
}

func (rec *recorderGHQueryClient) TestWorkflowRuns(branch, event string, page int) (*github.WorkflowRuns, error) {
	wfRuns, err := rec.ghQueryClient.TestWorkflowRuns(branch, event, page)
	if err != nil {
		return wfRuns, err
	}
	if err := record(testWorkFlowRunsFilename(rec.path, branch, event, page), wfRuns); err != nil {
		return wfRuns, fmt.Errorf("failed to record workflow runs: %w", err)
	}
	return wfRuns, nil
}

func (rec *recorderGHQueryClient) TestWorkflowRunJobs(runID int64, filter string, page int) (*github.Jobs, error) {
	jobs, err := rec.ghQueryClient.TestWorkflowRunJobs(runID, filter, page)
	if err != nil {
		return jobs, err
	}
	if err := record(testWorkflowRunJobsFilename(rec.path, runID, filter, page), jobs); err != nil {
		return jobs, fmt.Errorf("failed to record workflow run jobs: %w", err)
	}
	return jobs, nil
}

func record[T any](filename string, obj *T) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object of type %T: %w", obj, err)
	}
	return os.WriteFile(filename, data, 0600)
}

type playbackGHQueryClient struct {
	path string
}

func newPlaybackGHQueryClient(path string) *playbackGHQueryClient {
	return &playbackGHQueryClient{path: path}
}

func (pbc *playbackGHQueryClient) TestWorkflowRuns(branch, event string, page int) (*github.WorkflowRuns, error) {
	return playback(&github.WorkflowRuns{}, testWorkFlowRunsFilename(pbc.path, branch, event, page))
}

func (pbc *playbackGHQueryClient) TestWorkflowRunJobs(runID int64, filter string, page int) (*github.Jobs, error) {
	return playback(&github.Jobs{}, testWorkflowRunJobsFilename(pbc.path, runID, filter, page))
}

func playback[T any](obj *T, filename string) (*T, error) {
	cleanFilename := filepath.Clean(filename)
	if !strings.HasPrefix(cleanFilename, recorderDir) {
		panic(fmt.Errorf("Unsafe input %q does not start with %q", cleanFilename, recorderDir))
	}
	data, err := os.ReadFile(cleanFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read playback file %q: %w", filename, err)
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to parse playback file into a %T: %w", obj, err)
	}
	return obj, err
}

func testWorkFlowRunsFilename(path string, branch, event string, page int) string {
	return filepath.Join(path, fmt.Sprintf("testWorkflowRuns-branch-%s-event-%s-page-%d.json", branch, event, page))
}

func testWorkflowRunJobsFilename(path string, runID int64, filter string, page int) string {
	return filepath.Join(path, fmt.Sprintf("testWorkflowJobs-runId-%d-filter-%s-page-%d.json", runID, filter, page))
}
