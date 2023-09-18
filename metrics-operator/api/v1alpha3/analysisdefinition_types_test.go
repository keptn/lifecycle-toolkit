package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestOperatorValue_GetFloatValue(t *testing.T) {
	o := OperatorValue{
		FixedValue: *resource.NewQuantity(15, resource.DecimalSI),
	}

	require.Equal(t, 15.0, o.GetFloatValue())
}
