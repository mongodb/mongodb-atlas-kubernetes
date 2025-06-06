package ginkgo

import (
	"github.com/onsi/ginkgo/v2"
)

func RegisterCallbacks() {
	ginkgo.ReportAfterSuite("Ginkgo Metrics", UpdateMetricsAfterSuite)
	ginkgo.ReportAfterSuite("Duration Reporter", ReportDuration)
}
