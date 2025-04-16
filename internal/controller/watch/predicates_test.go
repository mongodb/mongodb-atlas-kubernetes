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

package watch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestSelectNamespacesPredicate(t *testing.T) {
	tests := map[string]struct {
		namespaces   []string
		createEvent  event.CreateEvent
		updateEvent  event.UpdateEvent
		deleteEvent  event.DeleteEvent
		genericEvent event.GenericEvent
		expect       bool
	}{
		"should return true when there are no namespace to filter": {
			namespaces:   []string{},
			createEvent:  event.CreateEvent{},
			updateEvent:  event.UpdateEvent{},
			deleteEvent:  event.DeleteEvent{},
			genericEvent: event.GenericEvent{},
			expect:       true,
		},
		"should return true when matching namespace": {
			namespaces:   []string{"test"},
			createEvent:  event.CreateEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			updateEvent:  event.UpdateEvent{ObjectNew: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			deleteEvent:  event.DeleteEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			genericEvent: event.GenericEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			expect:       true,
		},
		"should return false when not matching namespace": {
			namespaces:   []string{"other"},
			createEvent:  event.CreateEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			updateEvent:  event.UpdateEvent{ObjectNew: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			deleteEvent:  event.DeleteEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			genericEvent: event.GenericEvent{Object: &akov2.AtlasProject{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}},
			expect:       false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := SelectNamespacesPredicate(tt.namespaces)
			assert.Equal(t, tt.expect, f.CreateFunc(tt.createEvent))
			assert.Equal(t, tt.expect, f.UpdateFunc(tt.updateEvent))
			assert.Equal(t, tt.expect, f.DeleteFunc(tt.deleteEvent))
			assert.Equal(t, tt.expect, f.GenericFunc(tt.genericEvent))
		})
	}
}
