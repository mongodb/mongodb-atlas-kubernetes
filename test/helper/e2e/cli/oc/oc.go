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

package oc

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	cli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli"
)

func Version() {
	session := cli.Execute("oc", "version")
	session.Wait()
}

func Login(code, serverAPI string) {
	session := cli.ExecuteWithoutWriter("oc", "login", "--token="+code, "--server="+serverAPI, "--insecure-skip-tls-verify")
	EventuallyWithOffset(1, session).Should(Say("Logged into"), "Can not login to "+serverAPI)
}

func Apply(path string) {
	session := cli.Execute("oc", "apply", "-f", path)
	session.Wait("2m")
}

func Delete(path string) {
	session := cli.Execute("oc", "delete", "-f", path)
	session.Wait("2m")
}
