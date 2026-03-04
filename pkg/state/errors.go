// Copyright 2026 MongoDB Inc
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

package state

import (
	"errors"
	"fmt"
)

var ErrOptionalAtlasAccessLost = errors.New("lost optional access to Atlas")

// AtlasAccessLostError returns an error with the appropriate message based on
// whether the Atlas access is optional or required.
// When deletion protection is enabled, the error is wrapped as optional.
func AtlasAccessLostError(err error, optional bool) error {
	if optional {
		return fmt.Errorf("%w: %w", ErrOptionalAtlasAccessLost, err)
	}
	return err
}
