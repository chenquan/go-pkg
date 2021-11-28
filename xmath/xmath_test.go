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

import "testing"

func TestMaxInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1",
			args{
				a: 1,
				b: 2,
			},
			2,
		}, {"2",
			args{
				a: 2,
				b: 1,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxInt64(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			2,
		}, {
			"2",
			args{
				a: 2,
				b: 1,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			1,
		}, {
			"2",
			args{
				a: 2,
				b: 1,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt64(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			1,
		}, {
			"2",
			args{
				a: 2,
				b: 1,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			2,
		}, {
			"2",
			args{
				a: 2,
				b: 1,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			1,
		}, {
			"2",
			args{
				a: 2,
				b: 1,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxFloat32(t *testing.T) {
	type args struct {
		a float32
		b float32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"1",
			args{
				a: 1,
				b: 2,
			},
			2,
		}, {"2",
			args{
				a: 2,
				b: 1,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFloat32(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinFloat32(t *testing.T) {
	type args struct {
		a float32
		b float32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			"1",
			args{
				a: 1,
				b: 2,
			},
			1,
		}, {
			"2",
			args{
				a: 3.0,
				b: 1,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinFloat32(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}
