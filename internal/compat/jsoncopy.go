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

package compat

import "encoding/json"

// JSONCopy will copy src to dst via JSON serialization/deserialization.
func JSONCopy(dst, src any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dst)
	if err != nil {
		return err
	}

	return nil
}
