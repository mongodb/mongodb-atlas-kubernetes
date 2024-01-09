package main

import (
	"fmt"
	"strings"
	"time"
)

type FlakinessQuerier interface {
}

type testFlakiness struct {
	testIdentifier
	tests []jobID
}

type slotFlakiness struct {
	interval
	flakyTests flakyRank
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
			if run.CreatedAt.Time.Before(notBefore) {
				return sfr, nil // data is returned in chronological descendent order
			}
			if !strings.HasPrefix(*run.Name, "Test") || (run.Conclusion != nil && *run.Conclusion != "success") {
				continue // if it failed completely, it is not flaky
			}
			rid := *run.ID
			failed, err := queryJobFlakiness(qc, rid)
			if err != nil {
				return nil, err
			}
			slot := slotForTimestamp(period, notAfter, run.CreatedAt.Time)
			for _, failure := range failed {
				registerFlakiness(sfr[slot], failure)
			}
		}
	}
}

func queryJobFlakiness(qc QueryClient, rid int64) ([]jobID, error) {
	jobs, err := qc.TestWorkflowRunJobs(rid, "all", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to query job run %d: %w", rid, err)
	}
	if len(jobs.Jobs) > PerPage {
		return nil, fmt.Errorf("too many jobs in run (%d > %d)", len(jobs.Jobs), PerPage)
	}
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
