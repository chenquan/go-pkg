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

import "github.com/chenquan/go-pkg/xerror"

const defaultRetryTimes = 3

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times int // number of retries.
	}
)

// WithRetry customize a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var batchError xerror.BatchError
	for i := 0; i <= options.times; i++ {
		if err := fn(); err != nil {
			batchError.Add(err)
		} else {
			return nil
		}
	}

	return batchError.Err()
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
