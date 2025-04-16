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

package dryrun

import (
	"errors"
	"sync"
)

type errorQueue struct {
	mu sync.Mutex // protects fields below

	active bool
	errs   []error
}

var reconcileErrors = &errorQueue{}

func AddTerminationError(err error) {
	reconcileErrors.mu.Lock()
	defer reconcileErrors.mu.Unlock()

	if !reconcileErrors.active {
		return
	}

	reconcileErrors.errs = append(reconcileErrors.errs, err)
}

func terminationError() error {
	reconcileErrors.mu.Lock()
	defer reconcileErrors.mu.Unlock()

	result := make([]error, 0, len(reconcileErrors.errs))
	result = append(result, reconcileErrors.errs...)

	return errors.Join(result...)
}

func clearTerminationErrors() {
	reconcileErrors.mu.Lock()
	defer reconcileErrors.mu.Unlock()

	reconcileErrors.errs = nil
}

func enableErrors() {
	reconcileErrors.mu.Lock()
	defer reconcileErrors.mu.Unlock()

	reconcileErrors.active = true
}
