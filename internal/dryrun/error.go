package dryrun

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type DryRunError struct {
	GVK                    string
	Namespace, Name        string
	EventType, Reason, Msg string
}

func NewDryRunError(kind schema.ObjectKind, meta metav1.ObjectMetaAccessor, eventtype, reason, messageFmt string, args ...interface{}) error {
	gvk := "unknown"
	if kind != nil {
		gvk = kind.GroupVersionKind().String()
	}

	namespace, name := "unknown", "unknown"
	if meta != nil {
		namespace = meta.GetObjectMeta().GetNamespace()
		name = meta.GetObjectMeta().GetName()
	}

	msg := fmt.Sprintf(messageFmt, args...)

	return &DryRunError{
		GVK:       gvk,
		Namespace: namespace,
		Name:      name,
		EventType: eventtype,
		Reason:    reason,
		Msg:       msg,
	}
}

func (e *DryRunError) Error() string {
	return fmt.Sprintf(
		"DryRun event GVK=%v, Namespace=%v, Name=%v, EventType=%v, Reason=%v, Message=%v",
		e.GVK, e.Namespace, e.Name, e.EventType, e.Reason, e.Msg,
	)
}
