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

package validate

import (
	"errors"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func Tags(tags []*akov2.TagSpec) error {
	tagsMap := make(map[string]struct{}, len(tags))

	for _, currTag := range tags {
		if _, ok := tagsMap[currTag.Key]; ok {
			return errors.New("duplicate keys found in tags, this is forbidden")
		}

		tagsMap[currTag.Key] = struct{}{}
	}

	return nil
}
