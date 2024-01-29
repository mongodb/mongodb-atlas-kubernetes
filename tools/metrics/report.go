package main

import (
	"encoding/json"
	"fmt"
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
	Entries    []*EntryInfo `json:"entries,omitempty"`
	EntryCount int          `json:"entry_count"`
}

type EntryInfo struct {
	TestName  string   `json:"test_name"`
	TestType  string   `json:"test_type"`
	Tests     []string `json:"entries,omitempty"`
	TestCount int      `json:"test_count"`
}

func report(query string) (string, error) {
	switch query {
	case "regressions":
		return format(regressions())
	case "flakiness":
		return format(flakiness())
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
		item := &SlotInfo{
			SlotName:   slotName(slot),
			Start:      sr.interval.start.Format(DayFormat),
			End:        sr.interval.end.Format(DayFormat),
			Entries:    regressionEntries(slot, sr.regressions),
			EntryCount: len(sr.regressions),
		}
		slots = append(slots, item)
	}
	return slots
}

func regressionEntries(slot int, regressions map[string]*testRegressions) []*EntryInfo {
	entries := []*EntryInfo{}
	for _, entry := range regressions {
		info := &EntryInfo{
			TestName:  entry.Name,
			TestType:  entry.testType.String(),
			TestCount: len(entry.regressions),
		}
		if slot == 0 {
			info.Tests = runURLs(entry.regressions)
		}
		entries = append(entries, info)
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
		entry := &SlotInfo{
			SlotName:   slotName(slot),
			Start:      sr.interval.start.Format(DayFormat),
			End:        sr.interval.end.Format(DayFormat),
			Entries:    flakyEntries(slot, sr.flakyTests),
			EntryCount: len(sr.flakyTests.rank),
		}
		slots = append(slots, entry)
	}
	return slots
}

func flakyEntries(slot int, flakyTests flakyRank) []*EntryInfo {
	entries := []*EntryInfo{}
	for _, entry := range flakyTests.rank {
		info := &EntryInfo{
			TestName:  entry.Name,
			TestType:  entry.testType.String(),
			TestCount: len(entry.tests),
		}
		if slot == 0 {
			info.Tests = jobURLs(entry.tests)
		}
		entries = append(entries, info)
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
	if slot == 0 {
		return "last week"
	}
	return fmt.Sprintf("%d weeks ago", slot+1)
}

func format(report *ReportInfo, err error) (string, error) {
	if os.Getenv("FORMAT") == "summary" {
		if err != nil {
			return "", err
		}
		return Summary(report), nil
	}
	return jsonize(report, err)
}

func Summary(report *ReportInfo) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Last %d weeks *%s* report *%s*\\n\\n",
		len(report.Slots), report.Type, report.Slots[0].End)
	totals := 0
	trend := []int{}
	for i := len(report.Slots) - 1; i >= 0; i-- {
		occurrences := 0
		for _, entry := range report.Slots[i].Entries {
			occurrences += entry.TestCount
		}
		trend = append(trend, occurrences)
		totals += occurrences
	}
	last := trend[len(report.Slots)-1]
	diff := trend[len(report.Slots)-2] - last
	decreasing := true
	direction := "Down"
	if diff < 0 {
		diff *= -1
		direction = "*UP*"
		decreasing = false
	}
	avg := float32(totals) / float32(len(report.Slots))
	avgDiff := avg - float32(last)
	level := "Below"
	below := true
	if avgDiff < 0 {
		avgDiff *= -1
		level = "*ABOVE*"
		below = false
	}
	good := decreasing && below
	perfect := good && trend[len(report.Slots)-1] == 0
	if perfect {
		fmt.Fprintf(&sb, "*PERFECT WEEK!*\\nStats:\\n")
	} else if good {
		fmt.Fprintf(&sb, "*Good, trending down and below average:*\\n")
	} else {
		fmt.Fprintf(&sb, "*NOT a good trend:*\\n")
	}
	for _, occurrences := range trend[0 : len(trend)-1] {
		fmt.Fprintf(&sb, "%d, ", occurrences)
	}
	fmt.Fprintf(&sb, "_*%d*_ <- last week\\n\\n", last)
	fmt.Fprintf(&sb, "- %s %d from last week.\\n", direction, diff)
	fmt.Fprintf(&sb, "- %.02f %s current average of %.02f per week.\\n", avgDiff, level, avg)

	if report.Slots[0].EntryCount > 0 {
		fmt.Fprintf(&sb, "Last week ranking:\\n\\n")
		fmt.Fprintf(&sb, "Top offender (make sure we have a jira in progress for this one):\\n\\n")
		entry := report.Slots[0].Entries[0]
		fmt.Fprintf(&sb, "- *%d* %s test: %s\\n", entry.TestCount, entry.TestType, entry.TestName)
		fmt.Fprintf(&sb, "\\nRest:\\n")
		for _, entry := range report.Slots[0].Entries[1:] {
			fmt.Fprintf(&sb, "- %d %s test: %s\\n", entry.TestCount, entry.TestType, entry.TestName)
		}

		fmt.Fprintf(&sb, "\\n\\nTop offender links:\\n\\n")
		for _, url := range report.Slots[0].Entries[0].Tests {
			fmt.Fprintf(&sb, "%s\\n", url)
		}
	}
	return sb.String()
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
