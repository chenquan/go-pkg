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

package xsync

import "bytes"

type (
	// A BatchError is an error that can hold multiple errors.
	BatchError struct {
		errs errorArray
	}

	errorArray []error
)

// Add adds err to be.
func (be *BatchError) Add(err error) {
	if err != nil {
		be.errs = append(be.errs, err)
	}
}

// Err returns an error that represents all errors.
func (be *BatchError) Err() error {
	switch len(be.errs) {
	case 0:
		return nil
	case 1:
		return be.errs[0]
	default:
		return be.errs
	}
}

// NotNil checks if any error inside.
func (be *BatchError) NotNil() bool {
	return len(be.errs) > 0
}

// Error returns a string that represents inside errors.
func (ea errorArray) Error() string {
	var buf bytes.Buffer

	for i := range ea {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(ea[i].Error())
	}

	return buf.String()
}
