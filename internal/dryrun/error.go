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
	"fmt"
)

const dryRunErrorPrefix = "DryRun event: "

type DryRunError struct {
	Msg string
}

func NewDryRunError(messageFmt string, args ...any) error {
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
