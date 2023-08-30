package analysis

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
)

func TestAssignTasks(t *testing.T) {
	tests := []struct {
		name       string
		tasks      []metricsapi.Objective
		numWorkers int
		expected   map[int][]metricsapi.Objective
	}{
		{
			name:       "Equal Workers and Tasks",
			tasks:      []metricsapi.Objective{{}, {}, {}},
			numWorkers: 3,
			expected: map[int][]metricsapi.Objective{
				1: {{}},
				2: {{}},
				3: {{}},
			},
		},
		{
			name:       "More Workers than Tasks",
			tasks:      []metricsapi.Objective{{}, {}},
			numWorkers: 3,
			expected: map[int][]metricsapi.Objective{
				1: {{}},
				2: {{}},
			},
		},
		{
			name:       "More Tasks",
			tasks:      []metricsapi.Objective{{}, {}, {}},
			numWorkers: 2,
			expected: map[int][]metricsapi.Objective{
				1: {{}, {}},
				2: {{}},
			},
		}, {
			name:       "No Tasks",
			tasks:      []metricsapi.Objective{},
			numWorkers: 2,
			expected:   map[int][]metricsapi.Objective{},
		}, {
			name:       "No Workers",
			tasks:      []metricsapi.Objective{},
			numWorkers: 0,
			expected:   map[int][]metricsapi.Objective{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assigner := TaskAssigner{
				numWorkers: test.numWorkers,
				tasks:      test.tasks,
			}
			result := assigner.AssignTasks()
			require.Equal(t, test.expected, result)
		})
	}
}
