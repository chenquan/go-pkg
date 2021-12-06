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

package xtask

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

// Do fn with ctx control.
func Do(ctx context.Context, do func() error) (err error) {
	doneChan := make(chan error, 1)
	panicChan := make(chan interface{}, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- fmt.Sprintf("%+v\n\n%s", r, strings.TrimSpace(string(debug.Stack())))
			}
		}()

		doneChan <- do()
	}()

	select {
	case p := <-panicChan:
		panic(p)
	case err = <-doneChan:
		return
	case <-ctx.Done():

		select {
		case p := <-panicChan:
			panic(p)
		default:
		}

		err = ctx.Err()
	}
	return
}

// DoWithTimeout fn with timeout control.
func DoWithTimeout(timeout time.Duration, do func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Do(ctx, do)
}
