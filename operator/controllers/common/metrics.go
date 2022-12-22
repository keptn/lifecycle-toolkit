package common

import (
	"context"
	"fmt"

	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList) ([]apicommon.GaugeFloatValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []apicommon.GaugeFloatValue{}

	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
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
		return nil, fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []apicommon.GaugeFloatValue{}
	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.GetPreviousVersion() != "" {
			err := client.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", reconcileObject.GetParentName(), reconcileObject.GetPreviousVersion()), Namespace: reconcileObject.GetNamespace()}, previousObject)
			if err != nil {
				return nil, nil
			}
			piWrapper2, err := interfaces.NewMetricsObjectWrapperFromClientObject(previousObject)
			if err != nil {
				return nil, err
			}
			if reconcileObject.IsEndTimeSet() {
				previousInterval := reconcileObject.GetEndTime().Sub(piWrapper2.GetStartTime())
				res = append(res, apicommon.GaugeFloatValue{
					Value:      previousInterval.Seconds(),
					Attributes: reconcileObject.GetDurationMetricsAttributes(),
				})
			}
		}
	}
	return res, nil
}

func GetActiveInstances(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList) ([]apicommon.GaugeValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []apicommon.GaugeValue{}
	for _, ro := range piWrapper.GetItems() {
		activeMetricsObject, _ := interfaces.NewActiveMetricsObjectWrapperFromClientObject(ro)
		gaugeValue := int64(0)
		if !activeMetricsObject.IsEndTimeSet() {
			gaugeValue = int64(1)
		}
		res = append(res, apicommon.GaugeValue{
			Value:      gaugeValue,
			Attributes: activeMetricsObject.GetActiveMetricsAttributes(),
		})
	}

	return res, nil
}
