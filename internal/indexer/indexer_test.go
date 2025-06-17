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

package indexer

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type managerMock struct {
	manager.Manager
	client.FieldIndexer

	fields sets.Set[string]
}

func (m *managerMock) GetClient() client.Client {
	return fake.NewFakeClient()
}

func (m *managerMock) GetFieldIndexer() client.FieldIndexer {
	return m
}

func (m *managerMock) IndexField(_ context.Context, obj client.Object, field string, _ client.IndexerFunc) error {
	if field == "" {
		return fmt.Errorf("error adding indexer for type %T: field is empty", obj)
	}

	if m.fields.Has(field) {
		return fmt.Errorf("error indexing field %q: field is already registered", field)
	}

	m.fields.Insert(field)

	return nil
}

func TestRegisterAll(t *testing.T) {
	err := RegisterAll(context.Background(), &managerMock{fields: sets.New[string]()}, zaptest.NewLogger(t))
	if err != nil {
		t.Error(err)
	}
}
