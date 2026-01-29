// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boilerplate

import (
	"github.com/dave/jennifer/jen"
)

// AddLicenseHeader adds the standard license header to generated files.
func AddLicenseHeader(f *jen.File) {
	f.HeaderComment("Copyright 2025 MongoDB Inc")
	f.HeaderComment("")
	f.HeaderComment("Licensed under the Apache License, Version 2.0 (the \"License\");")
	f.HeaderComment("you may not use this file except in compliance with the License.")
	f.HeaderComment("You may obtain a copy of the License at")
	f.HeaderComment("")
	f.HeaderComment("\thttp://www.apache.org/licenses/LICENSE-2.0")
	f.HeaderComment("")
	f.HeaderComment("Unless required by applicable law or agreed to in writing, software")
	f.HeaderComment("distributed under the License is distributed on an \"AS IS\" BASIS,")
	f.HeaderComment("WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.")
	f.HeaderComment("See the License for the specific language governing permissions and")
	f.HeaderComment("limitations under the License.")
}
