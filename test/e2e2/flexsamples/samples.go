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

package flexsamples

import _ "embed"

//go:embed with_groupid_create.yaml
var WithGroupIdCreate []byte

//go:embed with_groupid_update.yaml
var WithGroupIdUpdate []byte

//go:embed with_groupref_create.yaml
var WithGroupRefCreate []byte

//go:embed with_groupref_update.yaml
var WithGroupRefUpdate []byte

//go:embed test_group.yaml
var TestGroup []byte
