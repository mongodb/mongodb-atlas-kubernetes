package result

import (
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func NextState(s state.ResourceState, msg string) (ctrlstate.Result, error) {
	if !strings.HasSuffix(msg, ".") {
		msg = msg + "."
	}

	switch s {
	case state.StateCreated:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateImported:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateUpdated:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateDeleted:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateInitial:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateImportRequested:
		return ctrlstate.Result{NextState: s, StateMsg: msg}, nil

	case state.StateCreating:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: 15 * time.Second},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateUpdating:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: 15 * time.Second},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateDeleting:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: 15 * time.Second},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateDeletionRequested:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: 15 * time.Second},
			NextState: s,
			StateMsg:  msg,
		}, nil

	default:
		return ctrlstate.Result{}, fmt.Errorf("unknown state %v", s)
	}
}

func Error(s state.ResourceState, err error) (ctrlstate.Result, error) {
	return ctrlstate.Result{
		NextState: s,
	}, err
}
