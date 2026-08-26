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
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/yudai/gojsondiff"
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

type cleanupFunc func(map[string]any)

func cleanLinksField(data map[string]any) {
	delete(data, "links")
}

func cleanCreatedField(data map[string]any) {
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

	getReq := req.Clone(req.Context())
	getReq.Method = http.MethodGet
	getReq.Body = nil
	getReq.GetBody = nil
	getReq.ContentLength = 0

	getResp, err := t.transport.RoundTrip(getReq)
	if err != nil {
		return "", fmt.Errorf("failed to GET original resource: %w", err)
	}
	defer getResp.Body.Close()

	payloadFromGet, _ := io.ReadAll(getResp.Body)

	var payloadFromGetParsed map[string]any
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

	return formatChangedPaths(diff), nil
}

func formatChangedPaths(diff gojsondiff.Diff) string {
	if !diff.Modified() {
		return "no changes"
	}

	paths := collectChangedPaths(nil, diff.Deltas())
	sort.Strings(paths)

	return strings.Join(paths, "\n")
}

func collectChangedPaths(prefix []string, deltas []gojsondiff.Delta) []string {
	var paths []string

	for _, delta := range deltas {
		switch d := delta.(type) {
		case *gojsondiff.Object:
			paths = append(paths, collectChangedPaths(appendPath(prefix, d.PostPosition()), d.Deltas)...)
		case *gojsondiff.Array:
			paths = append(paths, collectChangedPaths(appendPath(prefix, d.PostPosition()), d.Deltas)...)
		case *gojsondiff.Added:
			paths = append(paths, expandValuePaths("added", appendPath(prefix, d.PostPosition()), d.Value)...)
		case *gojsondiff.Modified:
			paths = append(paths, describePath("modified", prefix, d.PostPosition()))
		case *gojsondiff.TextDiff:
			paths = append(paths, describePath("modified", prefix, d.PostPosition()))
		case *gojsondiff.Deleted:
			paths = append(paths, expandValuePaths("deleted", appendPath(prefix, d.PrePosition()), d.Value)...)
		case *gojsondiff.Moved:
			paths = append(paths, fmt.Sprintf("moved: %s -> %s",
				renderPath(appendPath(prefix, d.PrePosition())),
				renderPath(appendPath(prefix, d.PostPosition()))))
		default:
			paths = append(paths, fmt.Sprintf("changed: %s", renderPath(prefix)))
		}
	}

	return paths
}

func expandValuePaths(kind string, path []string, value any) []string {
	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			return []string{fmt.Sprintf("%s: %s", kind, renderPath(path))}
		}

		var paths []string
		for key, child := range v {
			paths = append(paths, expandValuePaths(kind, append(slices.Clone(path), key), child)...)
		}
		return paths

	case []any:
		if len(v) == 0 {
			return []string{fmt.Sprintf("%s: %s", kind, renderPath(path))}
		}

		var paths []string
		for i, child := range v {
			paths = append(paths, expandValuePaths(kind, append(slices.Clone(path), "["+strconv.Itoa(i)+"]"), child)...)
		}
		return paths

	default:
		return []string{fmt.Sprintf("%s: %s", kind, renderPath(path))}
	}
}

func describePath(kind string, prefix []string, position gojsondiff.Position) string {
	return fmt.Sprintf("%s: %s", kind, renderPath(appendPath(prefix, position)))
}

func appendPath(prefix []string, position gojsondiff.Position) []string {
	path := make([]string, 0, len(prefix)+1)
	path = append(path, prefix...)

	switch p := position.(type) {
	case gojsondiff.Index:
		return append(path, "["+strconv.Itoa(int(p))+"]")
	default:
		return append(path, position.String())
	}
}

func renderPath(path []string) string {
	if len(path) == 0 {
		return "."
	}

	var sb strings.Builder
	for _, segment := range path {
		if strings.HasPrefix(segment, "[") {
			sb.WriteString(segment)
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(segment)
	}

	return sb.String()
}
