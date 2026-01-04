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

package operator

import (
	"context"
	"flag"
	"os"
	"strconv"
	"sync"
)

func RunEmbeddedSet() bool {
	envSet, _ := strconv.ParseBool(os.Getenv("RUN_EMBEDDED"))
	return envSet
}

type RunnerFunc func(context.Context, *flag.FlagSet, []string) error

type EmbeddedOperator struct {
	runnerFunc RunnerFunc
	mutex      sync.Mutex
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFn   context.CancelFunc
	args       []string
}

func NewEmbeddedOperator(runnerFunc RunnerFunc, args []string) *EmbeddedOperator {
	return &EmbeddedOperator{runnerFunc: runnerFunc, args: args}
}

func (e *EmbeddedOperator) Start(t testingT) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	t.Logf("starting operator in-process with args: %v", e.args)

	if e.ctx != nil {
		return
	}
	e.ctx, e.cancelFn = context.WithCancel(context.Background())
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		e.runnerFunc(e.ctx, fs, e.args)
	}()
}

func (e *EmbeddedOperator) Running() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.ctx != nil
}

func (e *EmbeddedOperator) Wait(t testingT) {
	t.Logf("waiting for operator goroutines to stop")
	e.wg.Wait()
}

func (e *EmbeddedOperator) Stop(t testingT) {
	t.Logf("canceling operator context to force it to stop")
	e.cancelFn()
	e.Wait(t)
}
