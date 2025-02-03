package dryrun

import (
	"errors"
	"fmt"
)

const dryRunErrorPrefix = "DryRun event: "

type DryRunError struct {
	Msg string
}

func NewDryRunError(messageFmt string, args ...interface{}) error {
	msg := fmt.Sprintf(messageFmt, args...)

	return &DryRunError{
		Msg: msg,
	}
}

func (e *DryRunError) Error() string {
	return dryRunErrorPrefix + e.Msg
}

// containsDryRunErrors returns true if the given error contains at least one DryRunError.
//
// Note: we DO NOT want to export this as we do not want "special dry-run" cases in reconcilers.
// Reconcilers should behave exactly the same during dry-run as during regular reconciles.
func containsDryRunErrors(err error) bool {
	dErr := &DryRunError{}
	return errors.As(err, &dErr)
}
