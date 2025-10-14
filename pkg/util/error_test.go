package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinErrors(t *testing.T) {
	errA := errors.New("error A")
	errB := errors.New("error B")
	err := JoinErrors(errA, errB)
	assert.True(t, errors.Is(err, errA))
	assert.True(t, errors.Is(err, errB))
}
