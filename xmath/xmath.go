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
