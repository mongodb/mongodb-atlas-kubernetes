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
