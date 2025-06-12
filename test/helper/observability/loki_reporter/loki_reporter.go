// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loki_reporter

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/loki/v3/clients/pkg/promtail/api"
	"github.com/grafana/loki/v3/clients/pkg/promtail/client"
	"github.com/grafana/loki/v3/pkg/logproto"
	lokiflag "github.com/grafana/loki/v3/pkg/util/flagext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

type LokiReporter struct {
	lokiClient client.Client
}

func (l *LokiReporter) Write(p []byte) (n int, err error) {
	l.lokiClient.Chan() <- api.Entry{
		Entry: logproto.Entry{
			Timestamp: time.Now(),
			Line:      string(p),
		},
	}

	return len(p), nil
}

func (l *LokiReporter) Stop() {
	l.lokiClient.Stop()
}

func New(lokiURL, job string, loggerWriter io.Writer) (*LokiReporter, error) {
	u, err := url.Parse(lokiURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing loki URL: %w", err)
	}
	u.Path = "/loki/api/v1/push"

	logger := log.NewLogfmtLogger(log.NewSyncWriter(loggerWriter))

	lokiClient, err := client.New(
		client.NewMetrics(prometheus.DefaultRegisterer),
		client.Config{
			URL:       flagext.URLValue{URL: u},
			BatchWait: client.BatchWait,
			BatchSize: client.BatchSize,
			ExternalLabels: lokiflag.LabelSet{
				LabelSet: map[model.LabelName]model.LabelValue{
					"job": model.LabelValue(job),
				},
			},
			Timeout: client.Timeout,
		},
		0,     // disable max streams cap
		0,     // disable max line size
		false, // disable line truncation
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating loki client: %w", err)
	}
	return &LokiReporter{lokiClient: lokiClient}, nil
}
