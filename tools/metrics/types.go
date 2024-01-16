package main

import (
	"fmt"
	"path"
	"strings"
	"time"
)

const (
	ghURL = "https://github.com"

	akoAuthor = "mongodb"

	// #nosec G101 false positive detected by gosec linter, this is not a secret
	ako = "mongodb-atlas-kubernetes"

	akoWorkflowFilename = "test.yml"

	runsPathFmt = "%s/%s/actions/runs/%d"

	jobsPathFmt = "%s/%s/actions/runs/%d/job/%d"
)

const (
	// PerPage is the results per page on paged queries
	PerPage = 100
)

const (
	Weekly = 7 * 24 * time.Hour

	DayFormat = "2006/01/02"
)

type TestType int

const (
	Unit TestType = iota
	Integration
	E2E
)

func (tt TestType) String() string {
	switch tt {
	case Unit:
		return "Unit"
	case Integration:
		return "Integration"
	case E2E:
		return "e2e"
	default:
		return fmt.Sprintf("??? (unsupported test type %d)", tt)
	}
}

type runID int64

func (rid runID) String() string {
	return path.Join(ghURL, fmt.Sprintf(runsPathFmt, akoAuthor, ako, rid))
}

type jobID struct {
	Name  string
	RunID runID
	JobID int64
}

func (jid jobID) URL() string {
	return path.Join(ghURL, fmt.Sprintf(jobsPathFmt, akoAuthor, ako, jid.RunID, jid.JobID))
}

func (jid jobID) String() string {
	return fmt.Sprintf("%q %s", jid.Name, jid.URL())
}

type testIdentifier struct {
	Name     string
	testType TestType
}

func (tid testIdentifier) String() string {
	return fmt.Sprintf("%q %s", tid.Name, tid.testType)
}

type interval struct {
	start time.Time
	end   time.Time
}

func (it interval) String() string {
	return fmt.Sprintf("%s -> %s", it.start, it.end)
}

func slotForTimestamp(period time.Duration, notAfter, timestamp time.Time) int {
	slot := int(notAfter.Sub(timestamp) / period)
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
		Name:     testName,
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
