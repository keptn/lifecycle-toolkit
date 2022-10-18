package common

type KeptnPhase KeptnPhaseType

type KeptnPhaseType struct {
	LongName  string
	ShortName string
}

var (
	PhaseWorkloadPreDeployment  = KeptnPhaseType{LongName: "Workload Pre-Deployment", ShortName: "WorkloadPreDeploy"}
	PhaseWorkloadPostDeployment = KeptnPhaseType{LongName: "Workload Post-Deployment", ShortName: "WorkloadPostDeploy"}
	PhaseWorkloadPreEvaluation  = KeptnPhaseType{LongName: "Workload Pre-Evaluation", ShortName: "WorkloadPreEvaluation"}
	PhaseWorkloadPostEvaluation = KeptnPhaseType{LongName: "Workload Post-Evaluation", ShortName: "WorkloadPostEvaluation"}
	PhaseWorkloadDeployment     = KeptnPhaseType{LongName: "Workload Deployment", ShortName: "WorkloadDeploy"}
	PhaseAppPreDeployment       = KeptnPhaseType{LongName: "App Pre-Deployment", ShortName: "AppPreDeploy"}
	PhaseAppPostDeployment      = KeptnPhaseType{LongName: "App Post-Deployment", ShortName: "AppPostDeploy"}
	PhaseAppPreEvaluation       = KeptnPhaseType{LongName: "App Pre-Evaluation", ShortName: "AppPreEvaluation"}
	PhaseAppPostEvaluation      = KeptnPhaseType{LongName: "App Post-Evaluation", ShortName: "AppPostEvaluation"}
	PhaseAppDeployment          = KeptnPhaseType{LongName: "App Deployment", ShortName: "AppDeploy"}
)
