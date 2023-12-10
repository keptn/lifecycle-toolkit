package v1beta1

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	v1alpha3common "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1beta1.KeptnTask to the hub version (v1alpha3.KeptnTask)
//
//nolint:gocyclo
func (src *KeptnTask) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha3.KeptnTask)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnAppVersion)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TaskDefinition = src.Spec.TaskDefinition
	dst.Spec.Parameters.Inline = src.Spec.Parameters.Inline
	dst.Spec.Context.WorkloadName = src.Spec.Context.WorkloadName
	dst.Spec.Context.AppName = src.Spec.Context.AppName
	dst.Spec.Context.AppVersion = src.Spec.Context.AppVersion
	dst.Spec.Context.WorkloadVersion = src.Spec.Context.WorkloadVersion
	dst.Spec.Context.TaskType = src.Spec.Context.TaskType
	dst.Spec.Context.ObjectType = src.Spec.Context.ObjectType

	dst.Spec.SecureParameters.Secret = src.Spec.SecureParameters.Secret
	dst.Spec.Retries = src.Spec.Retries
	dst.Spec.Timeout = src.Spec.Timeout
	dst.Spec.Type = v1alpha3common.CheckType(src.Spec.Type)

	dst.Status.JobName = src.Status.JobName
	dst.Status.Message = src.Status.Message
	dst.Status.StartTime = src.Status.StartTime
	dst.Status.EndTime = src.Status.EndTime
	dst.Status.Reason = src.Status.Reason
	dst.Status.Status = v1alpha3common.KeptnState(src.Status.Status)

	return nil
}

// ConvertFrom converts from the hub version (v1alpha3.KeptnTask) to this version (v1beta1.KeptnTask)
//
//nolint:gocyclo
func (dst *KeptnTask) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha3.KeptnTask)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.ErrCannotCastKeptnAppVersion)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TaskDefinition = src.Spec.TaskDefinition
	dst.Spec.Context.AppVersion = src.Spec.Context.AppVersion
	dst.Spec.Context.WorkloadName = src.Spec.Context.WorkloadName
	dst.Spec.Context.AppName = src.Spec.Context.AppName
	dst.Spec.Context.WorkloadVersion = src.Spec.Context.WorkloadVersion
	dst.Spec.Context.ObjectType = src.Spec.Context.ObjectType
	dst.Spec.Parameters.Inline = src.Spec.Parameters.Inline
	dst.Spec.Context.TaskType = src.Spec.Context.TaskType
	dst.Spec.SecureParameters.Secret = src.Spec.SecureParameters.Secret
	dst.Spec.Type = common.CheckType(src.Spec.Type)
	dst.Spec.Retries = src.Spec.Retries
	dst.Spec.Timeout = src.Spec.Timeout

	dst.Status.JobName = src.Status.JobName
	dst.Status.Message = src.Status.Message
	dst.Status.EndTime = src.Status.EndTime
	dst.Status.Status = common.KeptnState(src.Status.Status)
	dst.Status.StartTime = src.Status.StartTime
	dst.Status.Reason = src.Status.Reason

	return nil
}
