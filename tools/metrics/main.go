package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Weeks = 7
)

type ReportInfo struct {
	Type      string      `json:"type"`
	SlotCount int         `json:"slot_count"`
	Slots     []*SlotInfo `json:"slots"`
	Total     int         `json:"total"`
	QueryTme  string      `json:"query_time"`
}

type SlotInfo struct {
	SlotName   string       `json:"slot_name"`
	Start      string       `json:"start"`
	End        string       `json:"end"`
	Entries    []*EntryInfo `json:"entries"`
	EntryCount int          `json:"entry_count"`
}

type EntryInfo struct {
	TestName  string   `json:"test_name"`
	TestType  string   `json:"test_type"`
	Tests     []string `json:"entries"`
	TestCount int      `json:"test_count"`
}

func report(query string) (string, error) {
	switch query {
	case "regressions":
		return jsonize(regressions())
	case "flakiness":
		return jsonize(flakiness())
	default:
		return "", fmt.Errorf("query type %q unsupported, can only be 'regressions' or 'flakiness'", query)
	}
}

func regressions() (*ReportInfo, error) {
	start := time.Now()
	results, err := QueryRegressions(NewDefaultQueryClient(), time.Now(), Weekly, Weeks)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	return &ReportInfo{
		Type:      "regressions",
		SlotCount: len(results),
		Slots:     regressionsSlots(results),
		Total:     results.count(),
		QueryTme:  fmt.Sprintf("%v", elapsed),
	}, nil
}

func regressionsSlots(results slotRegressionsResult) []*SlotInfo {
	slots := []*SlotInfo{}
	for slot, sr := range results {
		slots = append(slots, &SlotInfo{
			SlotName:   slotName(slot),
			Start:      sr.interval.start.Format(DayFormat),
			End:        sr.interval.end.Format(DayFormat),
			Entries:    regressionEntries(sr.regressions),
			EntryCount: len(sr.regressions),
		})
	}
	return slots
}

func regressionEntries(regressions map[string]*testRegressions) []*EntryInfo {
	entries := []*EntryInfo{}
	for _, entry := range regressions {
		entries = append(entries, &EntryInfo{
			TestName:  entry.Name,
			TestType:  entry.testType.String(),
			Tests:     runURLs(entry.regressions),
			TestCount: len(entry.regressions),
		})
	}
	return entries
}

func runURLs(runIDs []runID) []string {
	urls := []string{}
	for _, rid := range runIDs {
		urls = append(urls, rid.String())
	}
	return urls
}

func flakiness() (*ReportInfo, error) {
	start := time.Now()
	results, err := QueryFlakiness(NewDefaultQueryClient(), time.Now(), Weekly, Weeks)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	return &ReportInfo{
		Type:      "flakiness",
		SlotCount: len(results),
		Slots:     flakinessSlots(results),
		Total:     results.count(),
		QueryTme:  fmt.Sprintf("%v", elapsed),
	}, nil
}

func flakinessSlots(results slotFlakinessResult) []*SlotInfo {
	slots := []*SlotInfo{}
	for slot, sr := range results {
		slots = append(slots, &SlotInfo{
			SlotName:   slotName(slot),
			Start:      sr.interval.start.Format(DayFormat),
			End:        sr.interval.end.Format(DayFormat),
			Entries:    flakyEntries(sr.flakyTests),
			EntryCount: len(sr.flakyTests.rank),
		})
	}
	return slots
}

func flakyEntries(flakyTests flakyRank) []*EntryInfo {
	entries := []*EntryInfo{}
	for _, entry := range flakyTests.rank {
		entries = append(entries, &EntryInfo{
			TestName:  entry.Name,
			TestType:  entry.testType.String(),
			Tests:     jobURLs(entry.tests),
			TestCount: len(entry.tests),
		})
	}
	return entries
}

func jobURLs(jobIDs []jobID) []string {
	urls := []string{}
	for _, jid := range jobIDs {
		urls = append(urls, jid.URL())
	}
	return urls
}

func slotName(slot int) string {
	switch slot {
	case 0:
		return "this week"
	case 1:
		return "last week"
	default:
		return fmt.Sprintf("%d weeks ago", slot)
	}
}

func jsonize(report *ReportInfo, err error) (string, error) {
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(report)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s {regressions|flakiness}\n", os.Args[0])
		os.Exit(1)
	}
	query := strings.ToLower(os.Args[1])
	if report, err := report(query); err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprint(os.Stdout, report)
	}
}
