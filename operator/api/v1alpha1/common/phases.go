package common

import (
	"strings"

	"go.opentelemetry.io/otel/propagation"
)

type KeptnPhase KeptnPhaseType

type KeptnPhaseType struct {
	LongName  string
	ShortName string
}

func (p KeptnPhaseType) IsEvaluation() bool {
	return strings.Contains(p.ShortName, "DeployEvaluations")
}

func (p KeptnPhaseType) IsPreEvaluation() bool {
	return strings.Contains(p.ShortName, "PreDeployEvaluations")
}

func (p KeptnPhaseType) IsPostEvaluation() bool {
	return strings.Contains(p.ShortName, "PostDeployEvaluations")
}

func (p KeptnPhaseType) IsTask() bool {
	return strings.Contains(p.ShortName, "DeployTasks")
}

func (p KeptnPhaseType) IsPreTask() bool {
	return strings.Contains(p.ShortName, "PreDeployTasks")
}

func (p KeptnPhaseType) IsPostTask() bool {
	return strings.Contains(p.ShortName, "PostDeployTasks")
}

var (
	PhaseWorkloadPreDeployment  = KeptnPhaseType{LongName: "Workload Pre-Deployment Tasks", ShortName: "WorkloadPreDeployTasks"}
	PhaseWorkloadPostDeployment = KeptnPhaseType{LongName: "Workload Post-Deployment Tasks", ShortName: "WorkloadPostDeployTasks"}
	PhaseWorkloadPreEvaluation  = KeptnPhaseType{LongName: "Workload Pre-Deployment Evaluations", ShortName: "WorkloadPreDeployEvaluations"}
	PhaseWorkloadPostEvaluation = KeptnPhaseType{LongName: "Workload Post-Deployment Evaluations", ShortName: "WorkloadPostDeployEvaluations"}
	PhaseWorkloadDeployment     = KeptnPhaseType{LongName: "Workload Deployment", ShortName: "WorkloadDeploy"}
	PhaseAppPreDeployment       = KeptnPhaseType{LongName: "App Pre-Deployment Tasks", ShortName: "AppPreDeployTasks"}
	PhaseAppPostDeployment      = KeptnPhaseType{LongName: "App Post-Deployment Tasks", ShortName: "AppPostDeployTasks"}
	PhaseAppPreEvaluation       = KeptnPhaseType{LongName: "App Pre-Deployment Evaluations", ShortName: "AppPreDeployEvaluations"}
	PhaseAppPostEvaluation      = KeptnPhaseType{LongName: "App Post-Deployment Evaluations", ShortName: "AppPostDeployEvaluations"}
	PhaseAppDeployment          = KeptnPhaseType{LongName: "App Deployment", ShortName: "AppDeploy"}
	PhaseCompleted              = KeptnPhaseType{LongName: "Completed", ShortName: "Completed"}
	PhaseCancelled              = KeptnPhaseType{LongName: "Cancelled", ShortName: "Cancelled"}
)

type PhaseTraceID map[string]propagation.MapCarrier

func (pid PhaseTraceID) SetPhaseTraceID(phase string, carrier propagation.MapCarrier) {
	pid[phase] = carrier
}

func (pid PhaseTraceID) GetPhaseTraceID(phase string) propagation.MapCarrier {
	return pid[phase]
}
