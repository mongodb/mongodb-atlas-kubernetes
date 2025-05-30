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

package e2e2_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e2"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

//go:embed samples/*
var samples embed.FS

func TestParseCRs(t *testing.T) {
	in, err := samples.Open("samples/sample.yml")
	require.NoError(t, err)
	defer in.Close()

	objs, err := e2e2.ParseCRs(in)
	require.NoError(t, err)
	assert.Len(t, objs, 2)
	assert.IsType(t, &corev1.Secret{}, objs[0])
	assert.IsType(t, &akov2next.AtlasThirdPartyIntegration{}, objs[1])
}
