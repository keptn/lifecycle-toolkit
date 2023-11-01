package common

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestGetRequestInfo(t *testing.T) {
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "example",
			Namespace: "test-namespace",
		}}

	info := GetRequestInfo(req)
	expected := map[string]string{
		"name":      "example",
		"namespace": "test-namespace",
	}
	require.Equal(t, expected, info)
}
