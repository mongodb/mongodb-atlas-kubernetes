// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultAPIPageSize = 100
)

type Exporter interface {
	Export(ctx context.Context) ([]client.Object, error)
}

type getPageFn[T any] func(int, int) ([]T, error)

func AllPages[T any](getPageFn getPageFn[T]) ([]T, error) {
	allPages := []T{}
	pageNum := 1
	itemsPerPage := DefaultAPIPageSize
	for {
		page, err := getPageFn(pageNum, itemsPerPage)
		if err != nil {
			return nil, fmt.Errorf("failed to get all pages: %w", err)
		}
		allPages = append(allPages, page...)
		if len(page) < itemsPerPage {
			break
		}
		pageNum += 1
	}
	return allPages, nil
}
