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

package result

import (
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	DefaultRequeueTIme = 15 * time.Second
)

func NextState(s state.ResourceState, msg string) (ctrlstate.Result, error) {
	if len(msg) > 0 && !strings.HasSuffix(msg, ".") {
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
			Result:    reconcile.Result{RequeueAfter: DefaultRequeueTIme},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateUpdating:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: DefaultRequeueTIme},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateDeleting:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: DefaultRequeueTIme},
			NextState: s,
			StateMsg:  msg,
		}, nil

	case state.StateDeletionRequested:
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: DefaultRequeueTIme},
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
