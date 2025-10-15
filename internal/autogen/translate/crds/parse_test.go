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

package crds_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/crds"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	scanner := bufio.NewScanner(bytes.NewBuffer(testdata.SampleCRDs))
	for _ = range 2 { // CRDs sample file has at least 2 CRDs
		def, err := crds.Parse(scanner)
		require.NoError(t, err)
		assert.NotNil(t, def)
	}
}
