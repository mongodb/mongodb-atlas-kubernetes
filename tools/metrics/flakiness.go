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

	"github.com/google/go-github/v57/github"
)

type FlakinessQuerier interface {
}

type testFlakiness struct {
	testIdentifier
	tests []jobID
}

type slotFlakiness struct {
	interval
	successfulCloudTestRuns int
	flakyTests              flakyRank
}

func (sr slotFlakiness) count() int {
	return len(sr.flakyTests.rank)
}

type flakyRank struct {
	rank []*testFlakiness
}

func (fr *flakyRank) String() string {
	var sb strings.Builder
	for _, tf := range fr.rank {
		fmt.Fprintf(&sb, "id: %v %d tests: %v\n", tf.testIdentifier, len(tf.tests), tf.tests)
	}
	return sb.String()
}

func (fr *flakyRank) add(id testIdentifier, jid jobID) {
	insertedAt := -1
	for i, tf := range fr.rank {
		if tf.Name == jid.Name {
			tf.tests = append(tf.tests, jid)
			insertedAt = i
		}
	}
	if insertedAt > 0 {
		j := insertedAt
		for ; j > 0 && len(fr.rank[j-1].tests) < len(fr.rank[j].tests); j-- {
			tf := fr.rank[j]
			fr.rank[j] = fr.rank[j-1]
			fr.rank[j-1] = tf
		}
	} else if insertedAt < 0 {
		fr.rank = append(fr.rank, &testFlakiness{id, []jobID{jid}})
	}
}

type slotFlakinessResult []*slotFlakiness

func (sfr slotFlakinessResult) count() int {
	if len(sfr) == 0 {
		return 0
	}
	total := 0
	for _, sr := range sfr {
		total += sr.count()
	}
	return total
}

func QueryFlakiness(qc QueryClient, notAfter time.Time, period time.Duration, slots int) (slotFlakinessResult, error) {
	sfr := make(slotFlakinessResult, 0, slots)
	for slot := 0; slot < slots; slot++ {
		sfr = append(sfr, &slotFlakiness{
			interval:   slotInterval(notAfter, period, slot),
			flakyTests: flakyRank{},
		})
	}
	page := 0
	notBefore := notAfter.Add(period * time.Duration(-slots))
	for {
		page += 1
		wfRuns, err := qc.TestWorkflowRuns("", "pull_request", page)
		if err != nil {
			return nil, fmt.Errorf("failed to query runs for %q: %w", akoWorkflowFilename, err)
		}
		if wfRuns.TotalCount != nil && *wfRuns.TotalCount == 0 {
			return sfr, nil
		}
		for _, run := range wfRuns.WorkflowRuns {
			if run.CreatedAt.Time.After(notAfter) {
				continue // skip anything after the end date
			}
			if run.CreatedAt.Time.Before(notBefore) {
				return sfr, nil // data is returned in chronological descendent order
			}
			if !strings.HasPrefix(*run.Name, "Test") {
				continue // skip non tests
			}
			slot := slotForTimestamp(period, notAfter, run.CreatedAt.Time)
			rid := *run.ID
			jobs, err := queryAllJobs(qc, rid)
			if err != nil {
				return nil, err
			}
			if run.Conclusion != nil && *run.Conclusion != "success" {
				// if it failed completely, it is not flaky
				continue
			}
			if isCloudTest(jobs) {
				sfr[slot].successfulCloudTestRuns += 1
			}
			failed, err := queryJobFlakiness(rid, jobs)
			if err != nil {
				return nil, err
			}
			for _, failure := range failed {
				registerFlakiness(sfr[slot], failure)
			}
		}
	}
}

func queryAllJobs(qc QueryClient, rid int64) (*github.Jobs, error) {
	jobs, err := qc.TestWorkflowRunJobs(rid, "all", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to query job run %d: %w", rid, err)
	}
	if len(jobs.Jobs) > PerPage {
		return nil, fmt.Errorf("too many jobs in run (%d > %d)", len(jobs.Jobs), PerPage)
	}
	return jobs, nil
}

func isCloudTest(jobs *github.Jobs) bool {
	for _, job := range jobs.Jobs {
		if *job.Name == "cloud-tests" && job.Conclusion != nil && *job.Conclusion != "skipped" {
			return true
		}
	}
	return false
}

func queryJobFlakiness(rid int64, jobs *github.Jobs) ([]jobID, error) {
	failed := []jobID{}
	for _, job := range jobs.Jobs {
		if job.Conclusion != nil && *job.Conclusion == "failure" {
			failed = append(failed, jobID{Name: *job.Name, RunID: runID(rid), JobID: *job.ID})
		}
	}
	return failed, nil
}

func registerFlakiness(reg *slotFlakiness, jid jobID) {
	reg.flakyTests.add(identify(jid.Name), jid)
}
