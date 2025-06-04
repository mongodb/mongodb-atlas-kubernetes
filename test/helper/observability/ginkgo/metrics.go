package ginkgo

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/prometheus/prometheus/prompb"
)

func sendMetrics(report types.Report) error {
	tss := make([]prompb.TimeSeries, 0, len(report.SpecReports))
	for _, r := range report.SpecReports {
		start := r.StartTime.UnixMilli()
		if start < 0 {
			continue
		}

		node := strings.Join(append(r.ContainerHierarchyTexts, r.LeafNodeText), " ")
		if node == "" {
			node = r.LeafNodeLocation.FileName + ":" + strconv.Itoa(r.LeafNodeLocation.LineNumber)
		}

		if r.LeafNodeLocation.CustomMessage != "" {
			node = node + " (" + r.LeafNodeLocation.CustomMessage + ")"
		}

		typ := nodeType(r.LeafNodeType)

		ts := prompb.TimeSeries{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "ginkgo_spec"},
				{Name: "node", Value: node},
				{Name: "type", Value: typ},
				{Name: "parallel_process", Value: strconv.Itoa(r.ParallelProcess)},
			},
			Samples: []prompb.Sample{
				{Value: float64(r.State), Timestamp: start},
			},
		}

		if r.EndTime.After(r.StartTime) {
			ts.Samples = append(ts.Samples, prompb.Sample{
				Value: float64(0.0), Timestamp: r.EndTime.UnixMilli(),
			})
		}

		tss = append(tss, ts)
	}

	var suiteSucceeded float64
	if report.SuiteSucceeded {
		suiteSucceeded = 1.0
	}
	suite := prompb.TimeSeries{
		Labels: []prompb.Label{
			{Name: "__name__", Value: "ginkgo_suite"},
			{Name: "description", Value: report.SuiteDescription},
			{Name: "path", Value: report.SuitePath},
			{Name: "focus", Value: strings.Join(report.SuiteConfig.FocusStrings, ", ")},
		},
		Samples: []prompb.Sample{
			{Value: suiteSucceeded, Timestamp: report.StartTime.UnixMilli()},
		},
	}
	if report.EndTime.After(report.StartTime) {
		suite.Samples = append(suite.Samples, prompb.Sample{
			Value: float64(0.0), Timestamp: report.EndTime.UnixMilli(),
		})
	}
	tss = append(tss, suite)

	req := &prompb.WriteRequest{
		Timeseries: tss,
	}

	data, err := proto.Marshal(req)
	compressed := snappy.Encode(nil, data)

	httpReq, err := http.NewRequest("POST", "http://localhost:30000/api/v1/write", bytes.NewReader(compressed))
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode/100 != 2 {
		err = fmt.Errorf("server returned HTTP status %s", httpResp.Status)
	}

	return err
}

func UpdateMetricsAfterSuite(ctx ginkgo.SpecContext, report types.Report) {
	if err := sendMetrics(report); err != nil {
		ginkgo.GinkgoWriter.Printf("[ERROR] error sending metrics: %v\n", err)
	}
}

func nodeType(leafNodeType types.NodeType) string {
	var typ string
	switch leafNodeType {
	case types.NodeTypeInvalid:
		typ = "NodeTypeInvalid"
	case types.NodeTypeContainer:
		typ = "NodeTypeContainer"
	case types.NodeTypeIt:
		typ = "NodeTypeIt"
	case types.NodeTypeBeforeEach:
		typ = "NodeTypeBeforeEach"
	case types.NodeTypeJustBeforeEach:
		typ = "NodeTypeJustBeforeEach"
	case types.NodeTypeAfterEach:
		typ = "NodeTypeAfterEach"
	case types.NodeTypeJustAfterEach:
		typ = "NodeTypeJustAfterEach"
	case types.NodeTypeBeforeAll:
		typ = "NodeTypeBeforeAll"
	case types.NodeTypeAfterAll:
		typ = "NodeTypeAfterAll"
	case types.NodeTypeBeforeSuite:
		typ = "NodeTypeBeforeSuite"
	case types.NodeTypeSynchronizedBeforeSuite:
		typ = "NodeTypeSynchronizedBeforeSuite"
	case types.NodeTypeAfterSuite:
		typ = "NodeTypeAfterSuite"
	case types.NodeTypeSynchronizedAfterSuite:
		typ = "NodeTypeSynchronizedAfterSuite"
	case types.NodeTypeReportBeforeEach:
		typ = "NodeTypeReportBeforeEach"
	case types.NodeTypeReportAfterEach:
		typ = "NodeTypeReportAfterEach"
	case types.NodeTypeReportBeforeSuite:
		typ = "NodeTypeReportBeforeSuite"
	case types.NodeTypeReportAfterSuite:
		typ = "NodeTypeReportAfterSuite"
	case types.NodeTypeCleanupInvalid:
		typ = "NodeTypeCleanupInvalid"
	case types.NodeTypeCleanupAfterEach:
		typ = "NodeTypeCleanupAfterEach"
	case types.NodeTypeCleanupAfterAll:
		typ = "NodeTypeCleanupAfterAll"
	case types.NodeTypeCleanupAfterSuite:
		typ = "NodeTypeCleanupAfterSuite"
	}
	return typ
}
