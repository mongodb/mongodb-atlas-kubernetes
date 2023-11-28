package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

func regressions() (string, error) {
	buf := bytes.NewBufferString("")
	start := time.Now()
	weeks := 7
	srr, err := QueryRegressions(NewDefaultQueryClient(), time.Now(), Weekly, weeks)
	if err != nil {
		return "", err
	}
	elapsed := time.Since(start)
	fmt.Fprintf(buf, "Query took %v\n", elapsed)
	for slot, sr := range srr {
		switch slot {
		case 0:
			fmt.Fprintf(buf, "%s -> %s this week:", sr.interval.start.Format(DayFormat), sr.interval.end.Format(DayFormat))
		case 1:
			fmt.Fprintf(buf, "%s -> %s last week:", sr.interval.start.Format(DayFormat), sr.interval.end.Format(DayFormat))
		default:
			fmt.Fprintf(buf, "%s -> %s %d weeks ago:", sr.interval.start.Format(DayFormat), sr.interval.end.Format(DayFormat), slot)
		}
		fmt.Fprintf(buf, " %d total regressions\n", sr.count())
		for _, tr := range sr.regressions {
			fmt.Fprintf(buf, "  => (%s) %s: %d regressions\n", tr.testType, tr.test, len(tr.regressions))
			for _, rid := range tr.regressions {
				fmt.Fprintf(buf, "     %s\n", rid)
			}
		}
	}
	fmt.Fprintf(buf, "%d total regressions in last %d weeks\n", srr.count(), weeks)
	return buf.String(), nil
}

func main() {
	if report, err := regressions(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprint(os.Stdout, report)
	}
}
