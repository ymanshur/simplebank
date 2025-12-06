package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymanshur/simplebank/internal/typex"
)

func TestJoinErrors(t *testing.T) {
	errA := errors.New("error A")
	var notFoundErr typex.ErrDataNotFound
	errB := typex.NewErrDataNotFound("aaa")
	err := JoinErrors(errA, errB)
	assert.True(t, errors.Is(err, errA))
	assert.True(t, errors.As(err, &notFoundErr))
}
