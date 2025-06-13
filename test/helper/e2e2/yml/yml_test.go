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

package yml_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

//go:embed samples/*
var samples embed.FS

func TestParseObjects(t *testing.T) {
	in, err := samples.Open("samples/sample.yml")
	require.NoError(t, err)
	defer in.Close()

	objs, err := yml.ParseObjects(in)
	require.NoError(t, err)
	assert.Len(t, objs, 2)
	assert.IsType(t, &corev1.Secret{}, objs[0])
	assert.Equal(t, &akov2.AtlasProject{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasProject",
			APIVersion: "atlas.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "atlas-project",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "atlas-project",
		},
	}, objs[1])
}
