/*
 *
 *     Copyright 2021 chenquan
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */

package xstream

// Options defines the struct to customize a Stream.
type Options struct {
	workSize int
}

// Option defines the method to customize a Stream.
type Option func(options *Options)

// loadOptions return a Options
func loadOptions(options ...Option) *Options {
	op := new(Options)
	for _, option := range options {
		option(op)
	}
	// set the default pool size
	if op.workSize <= 0 {
		op.workSize = 1
	}
	return op
}

// WithOption return a Option interface
func WithOption(options *Options) Option {
	return func(ops *Options) {
		*ops = *options
	}
}

// WithWorkSize return a Option that set size of work
func WithWorkSize(size int) Option {
	return func(options *Options) {
		options.workSize = size
	}
}
