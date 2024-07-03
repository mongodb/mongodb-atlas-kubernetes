package indexer

import (
	"context"
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"go.uber.org/zap/zaptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type managerMock struct {
	manager.Manager
	client.FieldIndexer

	fields sets.Set[string]
}

func (m *managerMock) GetFieldIndexer() client.FieldIndexer {
	return m
}

func (m *managerMock) IndexField(ctx context.Context, obj client.Object, field string, extractValue client.IndexerFunc) error {
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
