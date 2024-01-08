package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
)

type RegressionQuerier interface {
	// TestWorkflowRuns are all the test run by AKO at page X (descendent order)
	TestWorkflowRuns(branch, event string, page int) (*github.WorkflowRuns, error)

	// TestWorkflowRunJobs are all the jobs at a given Workflow Run at page X (descendent order)
	TestWorkflowRunJobs(runID int64, filter string, page int) (*github.Jobs, error)
}

type testRegressions struct {
	testIdentifier
	regressions []runID
}

type slotRegressions struct {
	interval
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
	workflowFilename := ".github/workflows/test.yml"
	for {
		page += 1
		wfRuns, err := qc.TestWorkflowRuns("main", "push", page)
		if err != nil {
			return nil, fmt.Errorf("failed to query runs for %q: %w", workflowFilename, err)
		}
		for _, run := range wfRuns.WorkflowRuns {
			if run.CreatedAt.Before(notBefore) {
				return srr, nil // data is returned in chronological descendent order
			}
			if !strings.HasPrefix(*run.Name, "Test") || (run.Conclusion != nil && *run.Conclusion == "success") {
				continue
			}
			rid := *run.ID
			failed, err := queryJobRegressions(qc, rid)
			if err != nil {
				return nil, err
			}
			slot := slotForTimestamp(period, run.CreatedAt.Time)
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
		if *job.Conclusion == "failure" {
			failed = append(failed, *job.Name)
		}
	}
	return failed, nil
}

func slotForTimestamp(period time.Duration, timestamp time.Time) int {
	slot := int(time.Since(timestamp) / period)
	return slot
}

func slotInterval(notAfter time.Time, period time.Duration, slot int) interval {
	return interval{
		start: notAfter.Add(-(time.Duration(slot+1) * period)),
		end:   notAfter.Add(-(time.Duration(slot) * period)),
	}
}

func identify(testName string) testIdentifier {
	return testIdentifier{
		test:     testName,
		testType: testTypeFor(testName),
	}
}

func testTypeFor(name string) TestType {
	if strings.Contains(name, "unit-tests") {
		return Unit
	}
	if strings.Contains(name, "int-tests") {
		return Integration
	}
	return E2E
}

func registerRegression(reg *slotRegressions, id testIdentifier, rid runID) {
	tr, ok := reg.regressions[id.test]
	if !ok {
		reg.regressions[id.test] = &testRegressions{testIdentifier: id, regressions: []runID{rid}}
		return
	}
	tr.regressions = append(tr.regressions, rid)
}
