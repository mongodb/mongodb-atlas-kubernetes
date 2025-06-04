package ginkgo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	ginkgotypes "github.com/onsi/ginkgo/v2/types"
)

func ReportDuration(ctx ginkgo.SpecContext, report ginkgotypes.Report) {
	durationReport := struct {
		GrafanaURL    string `json:"grafana_url"`
		PrometheusURL string `json:"prometheus_url"`
	}{}

	grafanaURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:30001",
		// take the uid from grafana-config.yaml's .data.ginkgo.json "uid" JSON field.
		Path: "/d/12345678-1234-1234-1234-123456789000/ginkgo-status",
		RawQuery: (&url.URL{
			Path: fmt.Sprintf(
				"orgId=1&from=%v&to=%v",
				report.StartTime.UnixMilli(),
				report.EndTime.UnixMilli(),
			),
		}).EscapedPath(),
	}

	promURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:30000",
		Path:   "/graph",
		RawQuery: (&url.URL{
			Path: fmt.Sprintf(
				"g0.expr=up&g0.tab=0&g0.display_mode=lines&g0.show_exemplars=0&g0.range_input=%v&g0.end_input=%v",
				report.RunTime.Round(time.Second),
				report.EndTime.UTC().Format("2006-01-02 15:04:05"),
			),
		}).EscapedPath(),
	}

	durationReport.PrometheusURL = promURL.String()
	durationReport.GrafanaURL = grafanaURL.String()

	out, err := exec.Command("go", "list", "-m", "-f", `{{.Dir}}`).Output()
	if err != nil {
		ginkgo.GinkgoWriter.Printf("[ERROR] error executing command: %v\n", err)
		return
	}

	reportPath := filepath.Join(strings.TrimSpace(string(out)), "urls.json")
	file, _ := os.OpenFile(reportPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	encoder.Encode(durationReport)
}
