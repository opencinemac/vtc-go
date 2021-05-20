package rate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMustNewPanic tests that mustNew panics if we get a bad input.
func TestMustNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		mustNew(24, NTSC(100))
	})
}
