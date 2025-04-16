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

// TODO: is this even used?
package utils

import (
	"bufio"
	"os"
	"regexp"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
)

type Secrets map[string]string

// GetSecretsFromActrc act use .actrc file for define the secrets, we can use the same file
func GetSecretsFromActrc() (Secrets, error) {
	s := make(map[string]string)
	file, err := os.Open(config.ActrcPath)
	if err != nil {
		return s, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("^-s (?P<key>.+)=(?P<value>.+)$")
		if re.MatchString(line) {
			matchers := re.FindStringSubmatch(line)
			s[matchers[1]] = matchers[2]
		}
	}
	return s, err
}

func GetSecretEnvOrActrc(keys []string) Secrets {
	actrcSecrets, _ := GetSecretsFromActrc()
	finalSecrets := make(map[string]string)
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			if actValue, ok := actrcSecrets[key]; ok {
				finalSecrets[key] = actValue
			}
		} else {
			finalSecrets[key] = value
		}
	}
	return finalSecrets
}
