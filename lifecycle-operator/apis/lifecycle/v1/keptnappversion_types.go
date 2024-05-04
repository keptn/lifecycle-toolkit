/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnAppVersionSpec defines the desired state of KeptnAppVersion
type KeptnAppVersionSpec struct {
	KeptnAppContextSpec `json:",inline"`
	KeptnAppSpec        `json:",inline"`
	// AppName is the name of the KeptnApp.
	AppName string `json:"appName"`
	// PreviousVersion is the version of the KeptnApp that has been deployed prior to this version.
	// +optional
	PreviousVersion string `json:"previousVersion,omitempty"`
	// TraceId contains the OpenTelemetry trace ID.
	// +optional
	TraceId map[string]string `json:"traceId,omitempty"`
}

// KeptnAppVersionStatus defines the observed state of KeptnAppVersion
type KeptnAppVersionStatus struct {
	// PreDeploymentStatus indicates the current status of the KeptnAppVersion's PreDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// PostDeploymentStatus indicates the current status of the KeptnAppVersion's PostDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// PromotionStatus indicates the current status of the KeptnAppVersion's Promotion phase.
	// +kubebuilder:default:=Pending
	// +optional
	PromotionStatus common.KeptnState `json:"promotionStatus,omitempty"`
	// PreDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PreDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// PostDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PostDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// WorkloadOverallStatus indicates the current status of the KeptnAppVersion's Workload deployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	WorkloadOverallStatus common.KeptnState `json:"workloadOverallStatus,omitempty"`
	// WorkloadStatus contains the current status of each KeptnWorkload that is part of the KeptnAppVersion.
	// +optional
	WorkloadStatus []WorkloadStatus `json:"workloadStatus,omitempty"`
	// CurrentPhase indicates the current phase of the KeptnAppVersion.
	// +optional
	CurrentPhase string `json:"currentPhase,omitempty"`
	// PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnAppVersion.
	// +optional
	PreDeploymentTaskStatus []ItemStatus `json:"preDeploymentTaskStatus,omitempty"`
	// PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnAppVersion.
	// +optional
	PostDeploymentTaskStatus []ItemStatus `json:"postDeploymentTaskStatus,omitempty"`
	// PromotionTaskStatus indicates the current state of each promotionTask of the KeptnAppVersion.
	// +optional
	PromotionTaskStatus []ItemStatus `json:"promotionTaskStatus,omitempty"`
	// PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnAppVersion.
	// +optional
	PreDeploymentEvaluationTaskStatus []ItemStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	// PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnAppVersion.
	// +optional
	PostDeploymentEvaluationTaskStatus []ItemStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	// PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnAppVersion.
	// +optional
	PhaseTraceIDs common.PhaseTraceID `json:"phaseTraceIDs,omitempty"`
	// Status represents the overall status of the KeptnAppVersion.
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`

	// StartTime represents the time at which the deployment of the KeptnAppVersion started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the deployment of the KeptnAppVersion finished.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
}

type WorkloadStatus struct {
	// Workload refers to a KeptnWorkload that is part of the KeptnAppVersion.
	// +optional
	Workload KeptnWorkloadRef `json:"workload,omitempty"`
	// Status indicates the current status of the KeptnWorkload.
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=keptnappversions,shortName=kav
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.appName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.currentPhase`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",priority=1,type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="PreDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.preDeploymentEvaluationStatus`
// +kubebuilder:printcolumn:name="WorkloadOverallStatus",priority=1,type=string,JSONPath=`.status.workloadOverallStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",priority=1,type=string,JSONPath=`.status.postDeploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.postDeploymentEvaluationStatus`
// +kubebuilder:printcolumn:name="PromotionStatus",priority=1,type=string,JSONPath=`.status.promotionStatus`

// KeptnAppVersion is the Schema for the keptnappversions API
type KeptnAppVersion struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnAppVersion.
	// +optional
	Spec KeptnAppVersionSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnAppVersion.
	// +optional
	Status KeptnAppVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnAppVersionList contains a list of KeptnAppVersion
type KeptnAppVersionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnAppVersion `json:"items"`
}

func (a KeptnAppVersionList) GetItems() []client.Object {
	b := make([]client.Object, 0, len(a.Items))
	for i := 0; i < len(a.Items); i++ {
		b = append(b, &a.Items[i])
	}
	return b
}

func (a *KeptnAppVersionList) RemoveDeprecated() {
	b := make([]KeptnAppVersion, 0, len(a.Items))
	for i := 0; i < len(a.Items); i++ {
		if a.Items[i].Status.Status != common.StateDeprecated {
			b = append(b, a.Items[i])
		}
	}
	a.Items = b
}

func init() {
	SchemeBuilder.Register(&KeptnAppVersion{}, &KeptnAppVersionList{})
}

func (a KeptnAppVersion) IsPreDeploymentCompleted() bool {
	return a.Status.PreDeploymentStatus.IsCompleted()
}

func (a KeptnAppVersion) IsPreDeploymentEvaluationCompleted() bool {
	return a.Status.PreDeploymentEvaluationStatus.IsCompleted()
}

func (a KeptnAppVersion) IsPreDeploymentSucceeded(isBlocking bool) bool {
	if isBlocking {
		return a.Status.PreDeploymentStatus.IsSucceeded()
	}
	return a.Status.PreDeploymentStatus.IsSucceeded() || a.Status.PreDeploymentStatus.IsWarning()
}

func (a KeptnAppVersion) IsPreDeploymentFailed() bool {
	return a.Status.PreDeploymentStatus.IsFailed()
}

func (a KeptnAppVersion) IsPreDeploymentEvaluationSucceeded(isBlocking bool) bool {
	if isBlocking {
		return a.Status.PreDeploymentEvaluationStatus.IsSucceeded()
	}
	return a.Status.PreDeploymentEvaluationStatus.IsSucceeded() || a.Status.PreDeploymentEvaluationStatus.IsWarning()
}

func (a KeptnAppVersion) IsPreDeploymentEvaluationFailed() bool {
	return a.Status.PreDeploymentEvaluationStatus.IsFailed()
}

func (a KeptnAppVersion) IsPostDeploymentCompleted() bool {
	return a.Status.PostDeploymentStatus.IsCompleted()
}

func (a KeptnAppVersion) IsPromotionCompleted() bool {
	return a.Status.PromotionStatus.IsCompleted()
}

func (a KeptnAppVersion) IsPostDeploymentEvaluationCompleted() bool {
	return a.Status.PostDeploymentEvaluationStatus.IsCompleted()
}

func (a KeptnAppVersion) IsPostDeploymentFailed() bool {
	return a.Status.PostDeploymentStatus.IsFailed()
}

func (a KeptnAppVersion) IsPromotionFailed() bool {
	return a.Status.PromotionStatus.IsFailed()
}

func (a KeptnAppVersion) IsPostDeploymentEvaluationSucceeded(isBlocking bool) bool {
	if isBlocking {
		return a.Status.PostDeploymentEvaluationStatus.IsSucceeded()
	}
	return a.Status.PostDeploymentEvaluationStatus.IsSucceeded() || a.Status.PostDeploymentEvaluationStatus.IsWarning()
}

func (a KeptnAppVersion) IsPostDeploymentEvaluationFailed() bool {
	return a.Status.PostDeploymentEvaluationStatus.IsFailed()
}

func (a KeptnAppVersion) IsPostDeploymentSucceeded(isBlocking bool) bool {
	if isBlocking {
		return a.Status.PostDeploymentStatus.IsSucceeded()
	}
	return a.Status.PostDeploymentStatus.IsSucceeded() || a.Status.PostDeploymentStatus.IsWarning()
}

func (a KeptnAppVersion) IsPromotionSucceeded() bool {
	return a.Status.PromotionStatus.IsSucceeded()
}

func (a KeptnAppVersion) AreWorkloadsCompleted() bool {
	return a.Status.WorkloadOverallStatus.IsCompleted()
}

func (a KeptnAppVersion) AreWorkloadsSucceeded() bool {
	return a.Status.WorkloadOverallStatus.IsSucceeded()
}

func (a KeptnAppVersion) AreWorkloadsFailed() bool {
	return a.Status.WorkloadOverallStatus.IsFailed()
}

func (a *KeptnAppVersion) SetStartTime() {
	if a.Status.StartTime.IsZero() {
		a.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (a *KeptnAppVersion) SetEndTime() {
	if a.Status.EndTime.IsZero() {
		a.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (a KeptnAppVersion) GetStartTime() time.Time {
	return a.Status.StartTime.Time
}

func (a KeptnAppVersion) GetEndTime() time.Time {
	return a.Status.EndTime.Time
}

func (a *KeptnAppVersion) IsStartTimeSet() bool {
	return !a.Status.StartTime.IsZero()
}

func (a *KeptnAppVersion) IsEndTimeSet() bool {
	return !a.Status.EndTime.IsZero()
}

func (a *KeptnAppVersion) Complete() {
	a.SetEndTime()
}

func (a KeptnAppVersion) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(a.Spec.AppName),
		common.AppVersion.String(a.Spec.Version),
		common.AppNamespace.String(a.Namespace),
	}
}

func (a KeptnAppVersion) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(a.Spec.AppName),
		common.AppVersion.String(a.Spec.Version),
		common.AppNamespace.String(a.Namespace),
		common.AppStatus.String(string(a.Status.Status)),
	}
}

func (a KeptnAppVersion) GetDurationMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(a.Spec.AppName),
		common.AppNamespace.String(a.Namespace),
		common.AppVersion.String(a.Spec.Version),
		common.AppPreviousVersion.String(a.Spec.PreviousVersion),
	}
}

func (a KeptnAppVersion) GetState() common.KeptnState {
	return a.Status.Status
}

func (a KeptnAppVersion) GetPreDeploymentTasks() []string {
	return a.Spec.PreDeploymentTasks
}

func (a KeptnAppVersion) GetPostDeploymentTasks() []string {
	return a.Spec.PostDeploymentTasks
}

func (a KeptnAppVersion) GetPromotionTasks() []string {
	return a.Spec.PromotionTasks
}

func (a KeptnAppVersion) GetPreDeploymentTaskStatus() []ItemStatus {
	return a.Status.PreDeploymentTaskStatus
}

func (a KeptnAppVersion) GetPostDeploymentTaskStatus() []ItemStatus {
	return a.Status.PostDeploymentTaskStatus
}

func (a KeptnAppVersion) GetPreDeploymentEvaluations() []string {
	return a.Spec.PreDeploymentEvaluations
}

func (a KeptnAppVersion) GetPostDeploymentEvaluations() []string {
	return a.Spec.PostDeploymentEvaluations
}

func (a KeptnAppVersion) GetPreDeploymentEvaluationTaskStatus() []ItemStatus {
	return a.Status.PreDeploymentEvaluationTaskStatus
}

func (a KeptnAppVersion) GetPostDeploymentEvaluationTaskStatus() []ItemStatus {
	return a.Status.PostDeploymentEvaluationTaskStatus
}

func (a KeptnAppVersion) GetPromotionTaskStatus() []ItemStatus {
	return a.Status.PromotionTaskStatus
}

func (a KeptnAppVersion) GetAppName() string {
	return a.Spec.AppName
}

func (a KeptnAppVersion) GetPreviousVersion() string {
	return a.Spec.PreviousVersion
}

func (a KeptnAppVersion) GetParentName() string {
	return a.Spec.AppName
}

func (a KeptnAppVersion) GetNamespace() string {
	return a.Namespace
}

func (a *KeptnAppVersion) SetState(state common.KeptnState) {
	a.Status.Status = state
}

func (a KeptnAppVersion) GetCurrentPhase() string {
	return a.Status.CurrentPhase
}

func (a *KeptnAppVersion) SetCurrentPhase(phase string) {
	a.Status.CurrentPhase = phase
}

func (a KeptnAppVersion) GetVersion() string {
	return a.Spec.Version
}

func (a KeptnAppVersion) GenerateTask(taskDefinition KeptnTaskDefinition, checkType common.CheckType) KeptnTask {
	return KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.GenerateTaskName(checkType, taskDefinition.Name),
			Namespace:   a.Namespace,
			Labels:      taskDefinition.Labels,
			Annotations: taskDefinition.Annotations,
		},
		Spec: KeptnTaskSpec{
			Context: TaskContext{
				AppName:    a.GetParentName(),
				AppVersion: a.GetVersion(),
				TaskType:   string(checkType),
				ObjectType: "App",
			},
			TaskDefinition:   taskDefinition.Name,
			Parameters:       TaskParameters{},
			SecureParameters: SecureParameters{},
			Type:             checkType,
			Retries:          taskDefinition.Spec.Retries,
			Timeout:          taskDefinition.Spec.Timeout,
		},
	}
}

func (a KeptnAppVersion) GenerateEvaluation(evaluationDefinition KeptnEvaluationDefinition, checkType common.CheckType) KeptnEvaluation {
	return KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.GenerateEvaluationName(checkType, evaluationDefinition.Name),
			Namespace: a.Namespace,
		},
		Spec: KeptnEvaluationSpec{
			AppVersion:           a.Spec.Version,
			AppName:              a.Spec.AppName,
			EvaluationDefinition: evaluationDefinition.Name,
			Type:                 checkType,
			FailureConditions: FailureConditions{
				RetryInterval: evaluationDefinition.Spec.FailureConditions.RetryInterval,
				Retries:       evaluationDefinition.Spec.FailureConditions.Retries,
			},
		},
	}
}

func (a KeptnAppVersion) GetSpanName(phase string) string {
	if phase == "" {
		return a.Name
	}
	return phase
}

func (a KeptnAppVersion) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(a.Spec.AppName),
		common.AppVersion.String(a.Spec.Version),
		common.AppNamespace.String(a.Namespace),
	}
}

func (a KeptnAppVersion) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(a.GetSpanAttributes()...)
}

func (a KeptnAppVersion) GetSpanKey(phase string) string {
	return fmt.Sprintf("%s.%s.%s.%s.%s", a.Spec.TraceId["traceparent"], a.Spec.AppName, a.ObjectMeta.Namespace, a.Spec.Version, phase)
}

func (v KeptnAppVersion) GetWorkloadNameOfApp(workloadName string) string {
	return fmt.Sprintf("%s-%s", v.Spec.AppName, workloadName)
}

//nolint:dupl
func (a *KeptnAppVersion) DeprecateRemainingPhases(phase common.KeptnPhaseType) {
	// no need to deprecate anything when promotion tasks fail
	if phase == common.PhasePromotion {
		return
	}
	// deprecate promotion tasks when post evaluation failed
	if phase == common.PhaseAppPostEvaluation {
		a.Status.PromotionStatus = common.StateDeprecated
	}
	// deprecate post evaluation when post tasks failed
	if phase == common.PhaseAppPostDeployment {
		a.Status.PostDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.PromotionStatus = common.StateDeprecated
	}
	// deprecate post evaluation and tasks when app deployment failed
	if phase == common.PhaseAppDeployment {
		a.Status.PostDeploymentStatus = common.StateDeprecated
		a.Status.PostDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.PromotionStatus = common.StateDeprecated
	}
	// deprecate app deployment, post tasks and evaluations if app pre-eval failed
	if phase == common.PhaseAppPreEvaluation {
		a.Status.PostDeploymentStatus = common.StateDeprecated
		a.Status.PostDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.WorkloadOverallStatus = common.StateDeprecated
		a.Status.PromotionStatus = common.StateDeprecated
	}
	// deprecate pre evaluations, app deployment and post tasks and evaluations when pre-tasks failed
	if phase == common.PhaseAppPreDeployment {
		a.Status.PostDeploymentStatus = common.StateDeprecated
		a.Status.PostDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.WorkloadOverallStatus = common.StateDeprecated
		a.Status.PreDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.PromotionStatus = common.StateDeprecated
	}
	// deprecate completely everything
	if phase == common.PhaseDeprecated {
		a.Status.PostDeploymentStatus = common.StateDeprecated
		a.Status.PostDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.WorkloadOverallStatus = common.StateDeprecated
		a.Status.PreDeploymentEvaluationStatus = common.StateDeprecated
		a.Status.PreDeploymentStatus = common.StateDeprecated
		a.Status.PromotionStatus = common.StateDeprecated
		a.Status.Status = common.StateDeprecated
		return
	}
	a.Status.Status = common.StateFailed
}

func (a *KeptnAppVersion) SetPhaseTraceID(phase string, carrier propagation.MapCarrier) {
	if a.Status.PhaseTraceIDs == nil {
		a.Status.PhaseTraceIDs = common.PhaseTraceID{}
	}
	a.Status.PhaseTraceIDs[common.GetShortPhaseName(phase)] = carrier
}

func (a KeptnAppVersion) GetEventAnnotations() map[string]string {
	return map[string]string{
		"appName":        a.Spec.AppName,
		"appVersion":     a.Spec.Version,
		"appVersionName": a.Name,
	}
}
