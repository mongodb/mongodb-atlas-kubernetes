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
	"errors"
	"net/http"
)

// Response is the paginated response containing the current page results and the total count.
// It is implemented by all supported SDK versions.
type Response[T any] interface {
	GetResults() []T
	GetTotalCount() int
}

// ListAll invokes the given pagination list function multiple times until the total count of responses is gathered.
// Once done, all paginated responses are returned.
// If an error occurs, the first error occurrence will be returned.
//
// This is taken over from https://github.com/mongodb/terraform-provider-mongodbatlas/blob/a5581ebb274dbcaffd43d330c5bfbbb329cae51d/internal/common/dsschema/page_request.go#L14-L31.
func ListAll[T any](ctx context.Context, listFunc func(ctx context.Context, pageNum int) (Response[T], *http.Response, error)) ([]T, error) {
	var results []T
	for currentPage := 1; ; currentPage++ {
		resp, _, err := listFunc(ctx, currentPage)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, errors.New("no response")
		}
		currentResults := resp.GetResults()
		results = append(results, currentResults...)
		if len(currentResults) == 0 || len(results) >= resp.GetTotalCount() {
			break
		}
	}
	return results, nil
}
