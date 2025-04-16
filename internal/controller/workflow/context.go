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

package workflow

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

// Context is a container for some information that is needed on all levels of function calls during reconciliation.
// It's mutable by design.
// Note, that it's NOT a Go Context but can carry one
type Context struct {
	// Log is the root logger used in the reconciliation. Used just for convenience to avoid passing log to each
	// method.
	// Is not supposed to be mutated!
	Log *zap.SugaredLogger

	// OrgID is the identifier of the Organization which the Atlas client was configured for
	OrgID string

	// Client is a mongodb atlas client used to make v1.0 API calls
	Client       *mongodbatlas.Client
	SdkClientSet *atlas.ClientSet

	status Status

	// This is the condition happened the last (most of all it contains the most important information that needs
	// to be logged)
	lastCondition *api.Condition

	// lastConditionWarn indicates if the last "terminal" condition was expected (for example wait for some resource)
	// or unexpected (any errors)
	lastConditionWarn bool

	// Go context, when appropriate
	Context context.Context
}

func NewContext(log *zap.SugaredLogger, conditions []api.Condition, context context.Context, obj runtime.Object) *Context {
	return &Context{
		status:  NewStatus(conditions),
		Log:     log,
		Context: context,
	}
}

func (c Context) Conditions() []api.Condition {
	return c.status.conditions
}

func (c Context) GetCondition(conditionType api.ConditionType) (condition api.Condition, found bool) {
	return c.status.GetCondition(conditionType)
}

func (c Context) StatusOptions() []api.Option {
	return c.status.options
}

func (c Context) LastCondition() *api.Condition {
	return c.lastCondition
}

func (c Context) LastConditionWarn() bool {
	return c.lastConditionWarn
}

func (c *Context) EnsureStatusOption(option api.Option) *Context {
	c.status.EnsureOption(option)
	return c
}

func (c *Context) EnsureCondition(condition api.Condition) *Context {
	c.status.EnsureCondition(condition)
	c.lastCondition = &condition
	return c
}

func (c *Context) SetConditionFromResult(conditionType api.ConditionType, result Result) *Context {
	condition := api.Condition{
		Type:    conditionType,
		Status:  corev1.ConditionFalse,
		Reason:  string(result.reason),
		Message: result.message,
	}
	if result.IsOk() {
		condition.Status = corev1.ConditionTrue
	}
	c.EnsureCondition(condition)
	c.lastConditionWarn = result.warning
	return c
}

func (c *Context) SetConditionFalse(conditionType api.ConditionType) *Context {
	c.EnsureCondition(api.Condition{
		Type:   conditionType,
		Status: corev1.ConditionFalse,
	})
	return c
}

func (c *Context) SetConditionFalseMsg(conditionType api.ConditionType, msg string) *Context {
	c.EnsureCondition(api.Condition{
		Type:    conditionType,
		Status:  corev1.ConditionFalse,
		Message: msg,
	})
	return c
}

func (c *Context) SetConditionTrue(conditionType api.ConditionType) *Context {
	c.EnsureCondition(api.Condition{
		Type:   conditionType,
		Status: corev1.ConditionTrue,
	})
	return c
}

func (c *Context) HasReason(reason ConditionReason) bool {
	for _, condition := range c.Conditions() {
		if condition.Reason == string(reason) {
			return true
		}
	}
	return false
}

func (c *Context) SetConditionTrueMsg(conditionType api.ConditionType, msg string) *Context {
	c.EnsureCondition(api.Condition{
		Type:    conditionType,
		Status:  corev1.ConditionTrue,
		Message: msg,
	})
	return c
}

func (c *Context) UnsetCondition(conditionType api.ConditionType) *Context {
	c.status.RemoveCondition(conditionType)
	return c
}
