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

package paging

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type page struct {
	results    []string
	totalCount int
}

func (r *page) GetResults() []string {
	if r == nil {
		return nil
	}
	return r.results
}

func (r *page) GetTotalCount() int {
	if r == nil {
		return 0
	}
	return r.totalCount
}

func responder(pages []*page) func(ctx context.Context, pageNum int) (Response[string], *http.Response, error) {
	totalCount := 0
	for _, p := range pages {
		if p == nil {
			continue
		}
		totalCount = totalCount + len(p.results)
	}

	for _, p := range pages {
		if p == nil {
			continue
		}
		p.totalCount = totalCount
	}

	return func(ctx context.Context, pageNum int) (Response[string], *http.Response, error) {
		if len(pages) == 0 {
			return nil, nil, nil
		}
		return pages[pageNum-1], nil, nil
	}
}

func TestAll(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name       string
		pages      []*page
		wantErr    string
		wantResult []string
	}{
		{
			name:    "no response",
			wantErr: "no response",
		},
		{
			name:    "empty response",
			pages:   []*page{},
			wantErr: "no response",
		},
		{
			name: "empty results",
			pages: []*page{
				{results: []string{}},
			},
			wantResult: nil,
		},
		{
			name: "single result",
			pages: []*page{
				{results: []string{"a"}},
			},
			wantResult: []string{"a"},
		},
		{
			name: "multiple results",
			pages: []*page{
				{results: []string{"a", "b"}},
			},
			wantResult: []string{"a", "b"},
		},
		{
			name: "one additional nil page",
			pages: []*page{
				{results: []string{"a", "b"}},
				nil,
			},
			wantResult: []string{"a", "b"},
		},
		{
			name: "one additional empty results page",
			pages: []*page{
				{results: []string{"a", "b"}},
				{results: []string{}},
			},
			wantResult: []string{"a", "b"},
		},
		{
			name: "multiple results",
			pages: []*page{
				{results: []string{"a", "b"}},
				{results: []string{"c", "d"}},
			},
			wantResult: []string{"a", "b", "c", "d"},
		},
		{
			name: "multiple results with nil page",
			pages: []*page{
				{results: []string{"a", "b"}},
				nil,
				{results: []string{"c", "d"}},
			},
			wantResult: []string{"a", "b"},
		},
		{
			name: "multiple results with empty results",
			pages: []*page{
				{results: []string{"a", "b"}},
				{results: []string{}},
				{results: []string{"c", "d"}},
			},
			wantResult: []string{"a", "b"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response, err := ListAll(ctx, responder(tc.pages))
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tc.wantErr, gotErr)
			require.Equal(t, tc.wantResult, response)
		})
	}
}
