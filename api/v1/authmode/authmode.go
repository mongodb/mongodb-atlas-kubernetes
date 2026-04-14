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

package authmode

import "slices"

type AuthMode string

const (
	Scram AuthMode = "SCRAM"
	X509  AuthMode = "X509"
)

type AuthModes []AuthMode

func (authModes AuthModes) CheckAuthMode(modeToCheck AuthMode) bool {
	return slices.Contains(authModes, modeToCheck)
}

func (authModes *AuthModes) AddAuthMode(modeToAdd AuthMode) {
	found := slices.Contains(*authModes, modeToAdd)

	if !found {
		*authModes = append(*authModes, modeToAdd)
	}
}

func (authModes *AuthModes) RemoveAuthMode(modeToRemove AuthMode) {
	var result AuthModes
	for _, mode := range *authModes {
		if mode != modeToRemove {
			result = append(result, mode)
		}
	}
	*authModes = result
}
