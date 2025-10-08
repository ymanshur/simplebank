package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidUsername(t *testing.T) {
	err := ValidUsername("")
	require.Nil(t, err)
}
