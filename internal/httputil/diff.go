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

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"go.uber.org/zap"
)

type TransportWithDiff struct {
	transport http.RoundTripper
	log       *zap.SugaredLogger
}

func NewTransportWithDiff(transport http.RoundTripper, log *zap.SugaredLogger) *TransportWithDiff {
	return &TransportWithDiff{
		transport: transport,
		log:       log,
	}
}

func (t *TransportWithDiff) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodPut || req.Method == http.MethodPatch {
		diffString, err := t.tryCalculateDiff(req,
			cleanLinksField,
			cleanCreatedField,
		)
		if err != nil {
			t.log.Debug("failed to calculate diff", zap.Error(err))
		} else {
			t.log.Desugar().Debug("JSON diff",
				zap.String("url", req.URL.String()),
				zap.String("diff", diffString),
			)
		}
	}
	return t.transport.RoundTrip(req)
}

type cleanupFunc func(map[string]interface{})

func cleanLinksField(data map[string]interface{}) {
	delete(data, "links")
}

func cleanCreatedField(data map[string]interface{}) {
	delete(data, "created")
}

func (t *TransportWithDiff) tryCalculateDiff(req *http.Request, cleanupFuncs ...cleanupFunc) (string, error) {
	var bodyCopy []byte
	if req.Body != nil {
		bodyCopy, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
	}

	defer func() {
		req.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
	}()

	getReq, _ := http.NewRequestWithContext(req.Context(), http.MethodGet, req.URL.String(), nil)
	getReq.Header = req.Header

	getResp, err := t.transport.RoundTrip(getReq)
	if err != nil {
		return "", fmt.Errorf("failed to GET original resource: %w", err)
	}
	defer getResp.Body.Close()

	payloadFromGet, _ := io.ReadAll(getResp.Body)

	var payloadFromGetParsed map[string]interface{}
	err = json.Unmarshal(payloadFromGet, &payloadFromGetParsed)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal payloadFromGetParsed JSON: %w", err)
	}

	for _, cFn := range cleanupFuncs {
		cFn(payloadFromGetParsed)
	}

	payloadBytes, err := json.Marshal(payloadFromGetParsed)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payloadFromGetParsed JSON: %w", err)
	}

	differ := gojsondiff.New()
	diff, err := differ.Compare(payloadBytes, bodyCopy)
	if err != nil {
		return "", fmt.Errorf("failed to compare JSON payloads: %w", err)
	}

	fmtr := formatter.NewAsciiFormatter(payloadFromGetParsed, formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       false,
	})

	diffString, err := fmtr.Format(diff)
	if err != nil {
		return "", fmt.Errorf("failed to format diff: %w", err)
	}

	return diffString, nil
}
