package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ExtraLongName = "loooooooooooooooooooooo00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ooooooo01234567891234567890123456789"

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
			State: StateDeprecated,
			Want:  true,
		},
		{
			State: StateCancelled,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.Want, tt.State.IsCompleted())
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

func TestHash(t *testing.T) {
	tests := []struct {
		in  int64
		out string
	}{
		{
			in:  int64(1),
			out: "6b86b273",
		},
		{
			in:  int64(2),
			out: "d4735e3a",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.out, Hash(tt.in))
		})
	}
}

func TestKeptnState_IsDeprecated(t *testing.T) {
	tests := []struct {
		State KeptnState
		Want  bool
	}{
		{
			State: StateSucceeded,
			Want:  false,
		},
		{
			State: StateDeprecated,
			Want:  true,
		},
		{
			State: StateCancelled,
			Want:  true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsDeprecated(), tt.Want)
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
			State: StateDeprecated,
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
			Name:    "deprecated",
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
			Name:  ExtraLongName,
			Want:  "pre-loooooooooooooooooooooo00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ooooooo0123456789123456-",
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
			Name:  ExtraLongName,
			Want:  "pre-eval-loooooooooooooooooooooo00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ooooooo01234567891-",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.True(t, strings.HasPrefix(GenerateEvaluationName(tt.Check, tt.Name), tt.Want))
		})
	}
}

func Test_GenerateJobName(t *testing.T) {
	tests := []struct {
		Name string
		Want string
	}{
		{
			Name: "short-name",
			Want: "short-name-",
		},
		{
			Name: "",
			Want: "-",
		},
		{
			Name: ExtraLongName,
			Want: "loooooooooooooooooooooo00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ooooooo01234567891234567890-",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.True(t, strings.HasPrefix(GenerateJobName(tt.Name), tt.Want))
		})
	}
}

func Test_MergeMaps(t *testing.T) {
	tests := []struct {
		In1  map[string]string
		In2  map[string]string
		Want map[string]string
	}{
		{
			In1:  nil,
			In2:  nil,
			Want: map[string]string{},
		},
		{
			In1: nil,
			In2: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
			},
			Want: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
			},
		},
		{
			In1: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
			},
			In2: nil,
			Want: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
			},
		},
		{
			In1: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
			},
			In2: map[string]string{
				"ll5": "ll6",
				"ll7": "ll8",
			},
			Want: map[string]string{
				"ll1": "ll2",
				"ll3": "ll4",
				"ll5": "ll6",
				"ll7": "ll8",
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, MergeMaps(tt.In1, tt.In2), tt.Want)
		})
	}
}

func TestIsOwnerSupported(t *testing.T) {
	type args struct {
		owner v1.OwnerReference
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Deployment -> true",
			args: args{
				owner: v1.OwnerReference{
					Kind: "Deployment",
				},
			},
			want: true,
		},
		{
			name: "DaemonSet-> true",
			args: args{
				owner: v1.OwnerReference{
					Kind: "DaemonSet",
				},
			},
			want: true,
		},
		{
			name: "ReplicaSet-> true",
			args: args{
				owner: v1.OwnerReference{
					Kind: "ReplicaSet",
				},
			},
			want: true,
		},
		{
			name: "StatefulSet-> true",
			args: args{
				owner: v1.OwnerReference{
					Kind: "StatefulSet",
				},
			},
			want: true,
		},
		{
			name: "Rollout-> true",
			args: args{
				owner: v1.OwnerReference{
					Kind: "Rollout",
				},
			},
			want: true,
		},
		{
			name: "Job-> false",
			args: args{
				owner: v1.OwnerReference{
					Kind: "Job",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOwnerSupported(tt.args.owner); got != tt.want {
				t.Errorf("IsOwnerSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}
