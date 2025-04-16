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
	"fmt"
	"strings"
	"time"
)

type testRegressions struct {
	testIdentifier
	regressions []runID
}

type slotRegressions struct {
	interval
	runs        int
	regressions map[string]*testRegressions
}

func (sr slotRegressions) count() int {
	return len(sr.regressions)
}

type slotRegressionsResult []*slotRegressions

func (srr slotRegressionsResult) count() int {
	if len(srr) == 0 {
		return 0
	}
	total := 0
	for _, sr := range srr {
		total += sr.count()
	}
	return total
}

func QueryRegressions(qc QueryClient, notAfter time.Time, period time.Duration, slots int) (slotRegressionsResult, error) {
	srr := make(slotRegressionsResult, 0, slots)
	for slot := 0; slot < slots; slot++ {
		srr = append(srr, &slotRegressions{
			interval:    slotInterval(notAfter, period, slot),
			regressions: map[string]*testRegressions{},
		})
	}
	page := 0
	notBefore := notAfter.Add(period * time.Duration(-slots))
	for {
		page += 1
		wfRuns, err := qc.TestWorkflowRuns("main", "push", page)
		if err != nil {
			return nil, fmt.Errorf("failed to query runs for %q: %w", akoWorkflowFilename, err)
		}
		if wfRuns.TotalCount != nil && *wfRuns.TotalCount == 0 {
			return srr, nil
		}
		for _, run := range wfRuns.WorkflowRuns {
			if run.CreatedAt.Time.After(notAfter) {
				continue // skip anything after the end date
			}
			if run.CreatedAt.Before(notBefore) {
				return srr, nil // data is returned in chronological descendent order
			}
			if !strings.HasPrefix(*run.Name, "Test") {
				continue
			}
			rid := *run.ID
			slot := slotForTimestamp(period, notAfter, run.CreatedAt.Time)
			srr[slot].runs += 1
			if run.Conclusion != nil && *run.Conclusion == "success" {
				continue
			}
			failed, err := queryJobRegressions(qc, rid)
			if err != nil {
				return nil, err
			}
			for _, failure := range failed {
				registerRegression(srr[slot], identify(failure), runID(rid))
			}
		}
	}
}

func queryJobRegressions(qc QueryClient, rid int64) ([]string, error) {
	jobs, err := qc.TestWorkflowRunJobs(rid, "latest", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to query job run %d: %w", rid, err)
	}
	if len(jobs.Jobs) > PerPage {
		return nil, fmt.Errorf("too many jobs in run (%d > %d)", len(jobs.Jobs), PerPage)
	}
	failed := []string{}
	for _, job := range jobs.Jobs {
		if job.Conclusion != nil && *job.Conclusion == "failure" {
			failed = append(failed, *job.Name)
		}
	}
	return failed, nil
}

func registerRegression(reg *slotRegressions, id testIdentifier, rid runID) {
	tr, ok := reg.regressions[id.Name]
	if !ok {
		reg.regressions[id.Name] = &testRegressions{testIdentifier: id, regressions: []runID{rid}}
		return
	}
	tr.regressions = append(tr.regressions, rid)
}
