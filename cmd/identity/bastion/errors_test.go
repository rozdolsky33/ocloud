package bastion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrAborted(t *testing.T) {
	assert.NotNil(t, ErrAborted)
	assert.Equal(t, "aborted by user", ErrAborted.Error())
	assert.IsType(t, &ErrAborted, &ErrAborted)
}
