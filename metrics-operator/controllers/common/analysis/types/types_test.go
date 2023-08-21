package types

import (
	"testing"

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
		MaximumScore: 6,
		TotalScore:   3,
	}

	require.Equal(t, 50.0, a.GetAchievedPercentage())

	a.MaximumScore = 1.0
	a.TotalScore = 1.0

	require.Equal(t, 100.0, a.GetAchievedPercentage())

	a.MaximumScore = 0.0
	a.TotalScore = 0.0

	require.Equal(t, 0.0, a.GetAchievedPercentage())
}
