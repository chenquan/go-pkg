package stream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithOption(t *testing.T) {
	options := &Options{workSize: 1}
	option := WithOption(options)
	ops := new(Options)
	option(ops)
	assert.Equal(t, options, ops)
}
func TestWithWorkSize(t *testing.T) {
	withWorkSize := WithWorkSize(1)
	ops := new(Options)
	withWorkSize(ops)
	assert.Equal(t, &Options{workSize: 1}, ops)
}
