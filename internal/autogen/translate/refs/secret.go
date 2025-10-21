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
//

package refs

import (
	"encoding/base64"
	"fmt"
)

func secretDecode(value string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}
	return string(bytes), nil
}

func secretEncode(value string) string {
	return base64.StdEncoding.EncodeToString(([]byte)(value))
}
