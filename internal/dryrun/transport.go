package dryrun

import (
	"errors"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
)

const DryRunReason = "DryRun"

var verbMap = map[string]string{
	http.MethodPost:   "create (" + http.MethodPost + ")",
	http.MethodPut:    "update (" + http.MethodPut + ")",
	http.MethodPatch:  "update (" + http.MethodPatch + ")",
	http.MethodDelete: "delete (" + http.MethodDelete + ")",
}

type DryRunTransport struct {
	Recorder record.EventRecorder
	Delegate http.RoundTripper
}

func NewDryRunTransport(recorder record.EventRecorder, delegate http.RoundTripper) *DryRunTransport {
	return &DryRunTransport{
		Recorder: recorder,
		Delegate: delegate,
	}
}

func (t *DryRunTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodGet:
	case http.MethodConnect:
	case http.MethodTrace:
	case http.MethodHead:
	default:
		obj, ok := runtimeObjectFrom(req.Context())
		if !ok {
			return nil, errors.New("no object present in context: cannot record dry run without an object")
		}

		meta, ok := obj.(metav1.ObjectMetaAccessor)
		if !ok {
			return nil, errors.New("object does not implement ObjectMetaAccessor")
		}

		verb, ok := verbMap[req.Method]
		if !ok {
			verb = "execute " + req.Method
		}
		msg := "Would %v %v"
		t.Recorder.Eventf(obj, v1.EventTypeNormal, DryRunReason, msg, verb, req.URL.Path)

		return nil, NewDryRunError(obj.GetObjectKind(), meta, v1.EventTypeNormal, DryRunReason, msg, verb, req.URL.Path)
	}

	return t.Delegate.RoundTrip(req)
}
