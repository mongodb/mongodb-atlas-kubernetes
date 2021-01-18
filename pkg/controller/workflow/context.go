package workflow

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

// Context is a container for some information that is needed on all levels of function calls during reconciliation.
// It's mutable by design.
// Note, that it's completely different from the Go Context
type Context struct {
	// Log is the root logger used in the reconciliation. Used just for convenience to avoid passing log to each
	// method.
	// Is not supposed to be mutated!
	Log *zap.SugaredLogger

	status Status

	// This is the condition happened the last (most of all it contains the most important information that needs
	// to be logged)
	lastCondition *status.Condition
}

func NewContext(log *zap.SugaredLogger, conditions []status.Condition) *Context {
	return &Context{
		status: NewStatus(conditions),
		Log:    log,
	}
}

func (c Context) Conditions() []status.Condition {
	return c.status.conditions
}

func (c Context) StatusOptions() []status.Option {
	return c.status.options
}

func (c Context) LastCondition() *status.Condition {
	return c.lastCondition
}

func (c *Context) EnsureStatusOption(option status.Option) *Context {
	c.status.EnsureOption(option)
	return c
}

func (c *Context) EnsureCondition(condition status.Condition) *Context {
	c.status.EnsureCondition(condition)
	c.lastCondition = &condition
	return c
}

func (c *Context) SetConditionFromResult(conditionType status.ConditionType, result Result) *Context {
	c.EnsureCondition(status.Condition{
		Type:    conditionType,
		Status:  corev1.ConditionFalse,
		Reason:  string(result.reason),
		Message: result.message,
	})
	return c
}

func (c *Context) SetConditionFalse(conditionType status.ConditionType) *Context {
	c.EnsureCondition(status.Condition{
		Type:   conditionType,
		Status: corev1.ConditionFalse,
	})
	return c
}

func (c *Context) SetConditionTrue(conditionType status.ConditionType) *Context {
	c.EnsureCondition(status.Condition{
		Type:   conditionType,
		Status: corev1.ConditionTrue,
	})
	return c
}
