package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeptnState_IsCompleted(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateProgressing,
			Want:  false,
		},
		{
			State: StateFailed,
			Want:  true,
		},
		{
			State: StateSucceeded,
			Want:  true,
		},
		{
			State: StateCancelled,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsCompleted(), tt.Want)
		})
	}
}

func TestKeptnState_IsSucceeded(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateProgressing,
			Want:  false,
		},
		{
			State: StateSucceeded,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsSucceeded(), tt.Want)
		})
	}
}

func TestKeptnState_IsFailed(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateSucceeded,
			Want:  false,
		},
		{
			State: StateFailed,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsFailed(), tt.Want)
		})
	}
}

func TestKeptnState_IsCancelled(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateSucceeded,
			Want:  false,
		},
		{
			State: StateCancelled,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsCancelled(), tt.Want)
		})
	}
}

func TestKeptnKeptnState_IsPending(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateSucceeded,
			Want:  false,
		},
		{
			State: StatePending,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsPending(), tt.Want)
		})
	}
}

func Test_UpdateStatusSummary(t *testing.T) {
	emmptySummary := StatusSummary{0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		State KeptnState
		Want  StatusSummary
	}{
		{
			State: StateProgressing,
			Want:  StatusSummary{0, 1, 0, 0, 0, 0, 0},
		},
		{
			State: StateFailed,
			Want:  StatusSummary{0, 0, 1, 0, 0, 0, 0},
		},
		{
			State: StateSucceeded,
			Want:  StatusSummary{0, 0, 0, 1, 0, 0, 0},
		},
		{
			State: StatePending,
			Want:  StatusSummary{0, 0, 0, 0, 1, 0, 0},
		},
		{
			State: "",
			Want:  StatusSummary{0, 0, 0, 0, 1, 0, 0},
		},
		{
			State: StateUnknown,
			Want:  StatusSummary{0, 0, 0, 0, 0, 1, 0},
		},
		{
			State: StateCancelled,
			Want:  StatusSummary{0, 0, 0, 0, 0, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, UpdateStatusSummary(tt.State, emmptySummary), tt.Want)
		})
	}
}

func Test_GetTotalCount(t *testing.T) {
	summary := StatusSummary{2, 0, 2, 1, 0, 3, 5}
	require.Equal(t, summary.GetTotalCount(), 11)
}

func Test_GeOverallState(t *testing.T) {
	tests := []struct {
		Name    string
		Summary StatusSummary
		Want    KeptnState
	}{
		{
			Name:    "failed",
			Summary: StatusSummary{0, 0, 1, 0, 0, 0, 0},
			Want:    StateFailed,
		},
		{
			Name:    "cancelled",
			Summary: StatusSummary{0, 0, 0, 0, 0, 0, 1},
			Want:    StateFailed,
		},
		{
			Name:    "progressing",
			Summary: StatusSummary{0, 1, 0, 0, 0, 0, 0},
			Want:    StateProgressing,
		},
		{
			Name:    "pending",
			Summary: StatusSummary{0, 0, 0, 0, 1, 0, 0},
			Want:    StatePending,
		},
		{
			Name:    "unknown",
			Summary: StatusSummary{0, 0, 0, 0, 0, 1, 0},
			Want:    StateUnknown,
		},
		{
			Name:    "unknown totalcount",
			Summary: StatusSummary{5, 0, 0, 0, 0, 1, 0},
			Want:    StateUnknown,
		},
		{
			Name:    "succeeded",
			Summary: StatusSummary{1, 0, 0, 1, 0, 0, 0},
			Want:    StateSucceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			require.Equal(t, GetOverallState(tt.Summary), tt.Want)
		})
	}
}

func Test_TruncateString(t *testing.T) {
	tests := []struct {
		Input string
		Max   int
		Want  string
	}{
		{
			Input: "some_string",
			Max:   20,
			Want:  "some_string",
		},
		{
			Input: "some_string",
			Max:   5,
			Want:  "some_",
		},
		{
			Input: "",
			Max:   5,
			Want:  "",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, TruncateString(tt.Input, tt.Max), tt.Want)
		})
	}
}

func Test_GenerateTaskName(t *testing.T) {
	tests := []struct {
		Check CheckType
		Name  string
		Want  string
	}{
		{
			Check: PreDeploymentCheckType,
			Name:  "short-name",
			Want:  "pre-short-name-",
		},
		{
			Check: PreDeploymentCheckType,
			Name:  "",
			Want:  "pre--",
		},
		{
			Check: PreDeploymentCheckType,
			Name:  "loooooooooooooooooooooooooooooooooooooong_name",
			Want:  "pre-looooooooooooooooooooooooooooooo-",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.True(t, strings.HasPrefix(GenerateTaskName(tt.Check, tt.Name), tt.Want))
		})
	}
}

func Test_GenerateEvaluationName(t *testing.T) {
	tests := []struct {
		Check CheckType
		Name  string
		Want  string
	}{
		{
			Check: PreDeploymentEvaluationCheckType,
			Name:  "short-name",
			Want:  "pre-eval-short-name-",
		},
		{
			Check: PreDeploymentEvaluationCheckType,
			Name:  "",
			Want:  "pre-eval--",
		},
		{
			Check: PreDeploymentEvaluationCheckType,
			Name:  "loooooooooooooooooooooooooooooooooooooong_name",
			Want:  "pre-eval-loooooooooooooooooooooooooo-",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.True(t, strings.HasPrefix(GenerateEvaluationName(tt.Check, tt.Name), tt.Want))
		})
	}
}
