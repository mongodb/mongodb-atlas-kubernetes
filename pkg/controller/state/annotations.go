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

package state

import "strings"

const (
	AnnotationReapplyTimestamp = "mongodb.internal.com/reapply-timestamp"
	AnnotationStateTracker     = "mongodb.internal.com/state-tracker"
)

var jsonPatchReplacer = strings.NewReplacer("/", "~1", "~", "~0")
