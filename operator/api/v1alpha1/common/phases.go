package common

type KeptnPhase KeptnPhaseType

type KeptnPhaseType struct {
	LongName  string
	ShortName string
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
)
