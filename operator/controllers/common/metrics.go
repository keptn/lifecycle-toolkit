package common

import (
	"context"
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList) ([]apicommon.GaugeFloatValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve instances: %w", err)
	}

	piWrapper, err := NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []apicommon.GaugeFloatValue{}

	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := NewPhaseItemWrapperFromClientObject(ro)
		if reconcileObject.IsEndTimeSet() {
			duration := reconcileObject.GetEndTime().Sub(reconcileObject.GetStartTime())
			res = append(res, apicommon.GaugeFloatValue{
				Value:      duration.Seconds(),
				Attributes: reconcileObject.GetDurationMetricsAttributes(),
			})
		}
	}

	return res, nil
}

func GetDeploymentInterval(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, previousObject client.Object) ([]apicommon.GaugeFloatValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve instances: %w", err)
	}

	piWrapper, err := NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []common.GaugeFloatValue{}
	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := NewPhaseItemWrapperFromClientObject(ro)
		if reconcileObject.GetPreviousVersion() != "" {
			err := client.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", reconcileObject.GetParentName(), reconcileObject.GetPreviousVersion()), Namespace: reconcileObject.GetNamespace()}, previousObject)
			if err != nil {
				return nil, nil
			}
			piWrapper2, err := NewPhaseItemWrapperFromClientObject(previousObject)
			if err != nil {
				return nil, err
			}
			if reconcileObject.IsEndTimeSet() {
				previousInterval := reconcileObject.GetEndTime().Sub(piWrapper2.GetStartTime())
				res = append(res, common.GaugeFloatValue{
					Value:      previousInterval.Seconds(),
					Attributes: reconcileObject.GetDurationMetricsAttributes(),
				})
			}
		}
	}
	return res, nil
}

func GetActiveInstances(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList) ([]common.GaugeValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve instances: %w", err)
	}

	piWrapper, err := NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []common.GaugeValue{}
	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := NewPhaseItemWrapperFromClientObject(ro)
		gaugeValue := int64(0)
		if !reconcileObject.IsEndTimeSet() {
			gaugeValue = int64(1)
		}
		res = append(res, common.GaugeValue{
			Value:      gaugeValue,
			Attributes: reconcileObject.GetActiveMetricsAttributes(),
		})
	}

	return res, nil
}
