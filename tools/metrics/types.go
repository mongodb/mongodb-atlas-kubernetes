package main

import (
	"fmt"
	"path"
	"time"
)

const (
	ghURL = "https://github.com"

	akoAuthor = "mongodb"

	// #nosec G101 false positive detected by gosec linter, this is not a secret
	ako = "mongodb-atlas-kubernetes"

	akoWorkflowFilename = "test.yml"

	runsPathFmt = "%s/actions/runs/%d"
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
	return path.Join(ghURL, fmt.Sprintf(runsPathFmt, ako, rid))
}

type testIdentifier struct {
	test     string
	testType TestType
}

type interval struct {
	start time.Time
	end   time.Time
}
