package workflow

import (
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
)

// Context is a container for some information that is needed on all levels of function calls during reconciliation.
// It's mutable by design.
// Note, that it's completely different from the Go Context
type Context struct {
	// Log is the root logger used in the reconciliation. Used just for convenience to avoid passing log to each
	// method.
	// Is not supposed to be mutated!
	Log *zap.SugaredLogger

	// Client is a mongodb atlas client used to make v1.0 API calls
	Client mongodbatlas.Client

	// Connection is an object encapsulating information about connecting to Atlas using API
	Connection atlas.Connection

	status Status

	// This is the condition happened the last (most of all it contains the most important information that needs
	// to be logged)
	lastCondition *status.Condition

	// lastConditionWarn indicates if the last "terminal" condition was expected (for example wait for some resource)
	// or unexpected (any errors)
	lastConditionWarn bool
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

func (c Context) LastConditionWarn() bool {
	return c.lastConditionWarn
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
	condition := status.Condition{
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
