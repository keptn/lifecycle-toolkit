package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLessOrEqual(t *testing.T) {
	require.True(t, isLessOrEqual("<"))
	require.True(t, isLessOrEqual("<="))
	require.False(t, isLessOrEqual(">"))
	require.False(t, isLessOrEqual(">="))
}

func TestIsGreaterOrEqual(t *testing.T) {
	require.False(t, isGreaterOrEqual("<"))
	require.False(t, isGreaterOrEqual("<="))
	require.True(t, isGreaterOrEqual(">"))
	require.True(t, isGreaterOrEqual(">="))
}
