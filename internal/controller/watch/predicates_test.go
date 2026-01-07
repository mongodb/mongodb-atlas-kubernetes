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

package watch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
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
			f := watch.SelectNamespacesPredicate(tt.namespaces)
			assert.Equal(t, tt.expect, f.CreateFunc(tt.createEvent))
			assert.Equal(t, tt.expect, f.UpdateFunc(tt.updateEvent))
			assert.Equal(t, tt.expect, f.DeleteFunc(tt.deleteEvent))
			assert.Equal(t, tt.expect, f.GenericFunc(tt.genericEvent))
		})
	}
}

func TestDeprecatedCommonPredicates(t *testing.T) {
	for _, tc := range []struct {
		title string
		old   *akov2.AtlasProject
		new   *akov2.AtlasProject
		want  bool
	}{
		{
			title: "no changes - resync",
			old:   sampleObj(),
			new:   sampleObj(),
			want:  true,
		},
		{
			title: "no gen change",
			old:   sampleObj(resourceVersion("0")),
			new:   sampleObj(resourceVersion("1")),
			want:  false,
		},
		{
			title: "finalizers changed",
			old:   sampleObj(resourceVersion("0")),
			new:   sampleObj(resourceVersion("1"), finalizers([]string{"finalize"})),
			want:  true,
		},
		{
			title: "skipped",
			old:   sampleObj(resourceVersion("0")),
			new:   sampleObj(resourceVersion("1"), skipAnnotation()),
			want:  false,
		},
		{
			title: "no longer skipped",
			old:   sampleObj(resourceVersion("0"), skipAnnotation()),
			new: sampleObj(
				resourceVersion("1"),
			),
			want: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			f := watch.DeprecatedCommonPredicates[client.Object]()
			assert.Equal(t, tc.want, f.Update(
				event.UpdateEvent{ObjectOld: tc.old, ObjectNew: tc.new}))
		})
	}
}

func TestDefaultPredicates(t *testing.T) {
	for _, tc := range []struct {
		title       string
		old         *akov2.AtlasProject
		new         *akov2.AtlasProject
		wantCreate  bool
		wantUpdate  bool
		wantDelete  bool
		wantGeneric bool
	}{
		{
			title:       "no changes",
			old:         sampleObj(),
			new:         sampleObj(),
			wantCreate:  true,
			wantUpdate:  true,
			wantDelete:  false,
			wantGeneric: true,
		},
		{
			title:       "no gen change",
			old:         sampleObj(resourceVersion("0")),
			new:         sampleObj(resourceVersion("1")),
			wantCreate:  true,
			wantUpdate:  false,
			wantDelete:  false,
			wantGeneric: true,
		},
		{
			title:       "finalizers set",
			old:         sampleObj(resourceVersion("0")),
			new:         sampleObj(resourceVersion("1"), finalizers([]string{"finalize"})),
			wantCreate:  true,
			wantUpdate:  false, // finalizer changes do not trigger updates unlike with deprecated
			wantDelete:  false,
			wantGeneric: true,
		},
		{
			title: "skipped",
			old: sampleObj(
				resourceVersion("0"),
			),
			new: sampleObj(
				resourceVersion("1"),
				skipAnnotation(),
			),
			wantCreate:  true,
			wantUpdate:  false,
			wantDelete:  false,
			wantGeneric: true,
		},
		{
			title:       "no longer skipped",
			old:         sampleObj(resourceVersion("0"), skipAnnotation()),
			new:         sampleObj(resourceVersion("1")),
			wantCreate:  true,
			wantUpdate:  true,
			wantDelete:  false,
			wantGeneric: true,
		},
		{
			title:       "finalizers removed",
			old:         sampleObj(resourceVersion("0"), finalizers([]string{"finalize"})),
			new:         sampleObj(resourceVersion("1")),
			wantCreate:  true,
			wantUpdate:  false,
			wantDelete:  false,
			wantGeneric: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			f := watch.DefaultPredicates[*akov2.AtlasProject]()
			assert.Equal(t, tc.wantCreate,
				f.Create(event.TypedCreateEvent[*akov2.AtlasProject]{Object: tc.new}), "on create")
			assert.Equal(t, tc.wantUpdate,
				f.Update(event.TypedUpdateEvent[*akov2.AtlasProject]{
					ObjectOld: tc.old, ObjectNew: tc.new,
				}), "on update")
			assert.Equal(t, tc.wantDelete,
				f.Delete(event.TypedDeleteEvent[*akov2.AtlasProject]{Object: tc.new}), "on delete")
			assert.Equal(t, tc.wantGeneric,
				f.Generic(event.TypedGenericEvent[*akov2.AtlasProject]{Object: tc.new}), "on generic event")
		})
	}
}

type optionFunc func(*akov2.AtlasProject) *akov2.AtlasProject

func sampleObj(opts ...optionFunc) *akov2.AtlasProject {
	p := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "project",
			Namespace:   "ns",
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "atlas-project",
		},
	}
	for _, opt := range opts {
		p = opt(p)
	}
	return p
}

func resourceVersion(rs string) optionFunc {
	return func(p *akov2.AtlasProject) *akov2.AtlasProject {
		p.ResourceVersion = rs
		return p
	}
}

func skipAnnotation() optionFunc {
	return func(p *akov2.AtlasProject) *akov2.AtlasProject {
		p.Annotations[customresource.ReconciliationPolicyAnnotation] =
			customresource.ReconciliationPolicySkip
		return p
	}
}

func finalizers(f []string) optionFunc {
	return func(p *akov2.AtlasProject) *akov2.AtlasProject {
		p.Finalizers = f
		return p
	}
}
