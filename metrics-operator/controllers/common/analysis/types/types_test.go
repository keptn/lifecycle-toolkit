package types

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
)

func TestTargetResult(t *testing.T) {
	target := TargetResult{
		WarningResult: OperatorResult{
			Fulfilled: true,
		},
		FailureResult: OperatorResult{
			Fulfilled: false,
		},
	}

	require.True(t, target.IsWarning())
	require.False(t, target.IsFailure())
}

func TestObjectiveResult(t *testing.T) {
	o := ObjectiveResult{
		Score: 1.0,
		Result: TargetResult{
			Warning: true,
		},
	}

	require.True(t, o.IsWarning())
	require.False(t, o.IsFailure())
	require.False(t, o.IsPass())
}

func TestAnalysisResult(t *testing.T) {
	a := AnalysisResult{
		MaximumScore: 3,
		TotalScore:   2,
	}

	a.CountScores(v1alpha3.Objective{
		Weight: 3,
	},
		ObjectiveResult{
			Score: 1,
		})

	require.Equal(t, 6.0, a.MaximumScore)
	require.Equal(t, 3.0, a.TotalScore)
	require.Equal(t, 50.0, a.GetAchievedPercentage())

	a.HandlePercentageScore(v1alpha3.AnalysisDefinition{
		Spec: v1alpha3.AnalysisDefinitionSpec{
			TotalScore: v1alpha3.TotalScore{
				PassPercentage:    80,
				WarningPercentage: 50,
			},
		},
	})

	require.True(t, a.Warning)
	require.False(t, a.Pass)

	a.MaximumScore = 0.0
	a.TotalScore = 0.0

	require.Equal(t, 100.0, a.GetAchievedPercentage())
}
