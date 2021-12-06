/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package xworker

import (
	"context"
	"github.com/chenquan/go-pkg/xtask"
)

// Worker is used to control the concurrency of goroutines.
type Worker struct {
	c chan struct{}
}

// NewWorker returns a Worker.
func NewWorker(size int) *Worker {
	return &Worker{c: make(chan struct{}, size)}
}

// Run executes function with ctx.
func (w *Worker) Run(ctx context.Context, run func()) {
	w.c <- struct{}{}
	defer func() {
		<-w.c
	}()
	_ = xtask.Do(ctx, func() error {
		run()
		return nil
	})
}
