package analysis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_generateQuery(t *testing.T) {

	tests := []struct {
		name      string
		query     string
		selectors map[string]string
		want      string
		wanterror string
	}{
		{
			name:  "successful, all args exist",
			query: "this is a {{.good}} query{{.dot}}",
			selectors: map[string]string{
				"good": "good",
				"dot":  ".",
			},
			want: "this is a good query.",
		},
		{
			name:  "no substitution, all args missing",
			query: "this is a {{.good}} query{{.dot}}",
			selectors: map[string]string{
				"bad":    "good",
				"dotted": ".",
			},
			want: "this is a <no value> query<no value>",
		},
		{
			name:  "no substitution, bad template",
			query: "this is a {{.good} query{{.dot}}",
			selectors: map[string]string{
				"bad":    "good",
				"dotted": ".",
			},
			want:      "",
			wanterror: "could not create a template:",
		},
		{
			name:      "nothing to do",
			query:     "this is a query",
			selectors: map[string]string{},
			want:      "this is a query",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateQuery(tt.query, tt.selectors)
			if tt.wanterror != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wanterror)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
