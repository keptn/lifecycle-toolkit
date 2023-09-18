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

func TestObjectReference_GetNamespace(t *testing.T) {
	o := ObjectReference{}

	require.Equal(t, "default", o.GetNamespace("default"))

	o.Namespace = "ns"

	require.Equal(t, "ns", o.GetNamespace("default"))
}
