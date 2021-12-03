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

package xstring

import "strings"

type (
	Joiner struct {
		b        strings.Builder
		notEmpty bool
		opts     *JoinerOptions
	}
	JoinerOptions struct {
		prefix string
		step   string
		suffix string
	}
	JoinerOption func(*JoinerOptions)
)

func WithJoinStep(step string) JoinerOption {
	return func(options *JoinerOptions) {
		options.step = step
	}
}

func WithJoinPrefix(prefix string) JoinerOption {
	return func(options *JoinerOptions) {
		options.prefix = prefix
	}
}

func WithJoinSuffix(suffix string) JoinerOption {
	return func(options *JoinerOptions) {
		options.suffix = suffix
	}
}

func WithJoin(prefix, step, suffix string) JoinerOption {
	return func(options *JoinerOptions) {
		options.prefix = prefix
		options.step = step
		options.suffix = suffix
	}
}

func NewJoin(opts ...JoinerOption) *Joiner {
	j := &Joiner{}
	j.loadOpts(opts...)
	return j
}

func (j *Joiner) loadOpts(opts ...JoinerOption) {
	op := new(JoinerOptions)
	for _, opt := range opts {
		opt(op)
	}
	j.opts = op
}

func (j *Joiner) WriteRune(r rune) (int, error) {
	j.tryWriteStep()
	n, _ := j.b.WriteRune(r)
	return n, nil
}

func (j *Joiner) WriteString(s string) (int, error) {
	j.tryWriteStep()
	n, _ := j.b.WriteString(s)
	return n, nil
}

func (j *Joiner) WriteByte(b byte) error {
	j.tryWriteStep()
	_ = j.b.WriteByte(b)
	return nil
}

func (j *Joiner) Write(p []byte) (int, error) {
	j.tryWriteStep()
	n, _ := j.b.Write(p)
	return n, nil
}

func (j *Joiner) String() string {
	return j.opts.prefix + j.b.String() + j.opts.suffix
}

func (j *Joiner) tryWriteStep() {
	if j.notEmpty {
		j.b.WriteString(j.opts.step)
	} else {
		j.notEmpty = true
	}
	return
}

func (j *Joiner) Grow(n int) {
	n = n - len(j.opts.prefix) - len(j.opts.suffix)
	j.b.Grow(n)
}

func (j *Joiner) Len() int {
	return len(j.opts.prefix) + j.b.Len() + len(j.opts.suffix)
}
