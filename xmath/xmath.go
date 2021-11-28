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

package xmath

import "math"

// MaxInt returns the maximum value.
func MaxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// MinInt returns minimum value.
func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

// MaxInt64 returns the maximum value.
func MaxInt64(a, b int64) int64 {
	return int64(math.Max(float64(a), float64(b)))
}

// MinInt64 returns minimum value.
func MinInt64(a, b int64) int64 {
	return int64(math.Min(float64(a), float64(b)))
}

// MaxFloat64 returns the maximum value.
func MaxFloat64(a, b float64) float64 {
	return math.Max(a, b)
}

// MinFloat64 returns the maximum value.
func MinFloat64(a, b float64) float64 {
	return math.Min(a, b)
}

// MaxFloat32 returns the maximum value.
func MaxFloat32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

// MinFloat32 returns the maximum value.
func MinFloat32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}
