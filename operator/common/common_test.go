package operatorcommon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateResourceName(t *testing.T) {
	tests := []struct {
		Name  string
		Input []string
		Max   int
		Min   int
		Want  string
	}{
		{
			Name: "parts not exceeding max, not min",
			Input: []string{
				"str1",
				"str2",
				"str3",
			},
			Max:  20,
			Min:  5,
			Want: "str1-str2-str3",
		},
		{
			Name: "1 part exceeding max",
			Input: []string{
				"str1111111111111111111111",
				"str2",
				"str3",
			},
			Max:  20,
			Min:  5,
			Want: "str1111111-str2-str3",
		},
		{
			Name: "2 part exceeding max",
			Input: []string{
				"str1",
				"str222222222222222222222222",
				"str3",
			},
			Max:  20,
			Min:  5,
			Want: "str1-str2222222-str3",
		},
		{
			Name: "1 and 2 part exceeding max",
			Input: []string{
				"str111111111111111111111",
				"str22222222",
				"str3",
			},
			Max:  20,
			Min:  5,
			Want: "str11-str222222-str3",
		},
		{
			Name: "1 and 2 part exceeding max, min needs to be reduced",
			Input: []string{
				"str111111111111111111111",
				"str22222222",
				"str3",
			},
			Max:  20,
			Min:  10,
			Want: "str11-str222222-str3",
		},
		{
			Name: "1 and 2 part exceeding max, min needs to be reduced",
			Input: []string{
				"str111111111111111111111",
				"str22222222",
				"str3",
			},
			Max:  20,
			Min:  20,
			Want: "str11-str222222-str3",
		},
		{
			Name: "1 and 2 part exceeding max, min needs to be reduced",
			Input: []string{
				"str111111111111111111111",
				"str22222222",
				"str3",
			},
			Max:  20,
			Min:  100,
			Want: "str111-str22222-str3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			require.Equal(t, tt.Want, CreateResourceName(tt.Max, tt.Min, tt.Input...))
		})
	}
}
