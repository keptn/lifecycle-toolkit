package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObjectReference_IsNamespaceSet(t *testing.T) {
	o := ObjectReference{}

	require.False(t, o.IsNamespaceSet())

	o.Namespace = "ns"

	require.True(t, o.IsNamespaceSet())
}
