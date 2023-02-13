package common

import (
	"context"
	"fmt"
	"strings"

	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ObserveDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge asyncfloat64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.IsEndTimeSet() {
			duration := reconcileObject.GetEndTime().Sub(reconcileObject.GetStartTime())
			gauge.Observe(ctx, duration.Seconds(), reconcileObject.GetDurationMetricsAttributes()...)
		}
	}

	return nil
}

func ObserveDeploymentInterval(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge asyncfloat64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	items := piWrapper.GetItems()
	for index := 0; index < len(items); index++ {
		ro := items[index]
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.GetPreviousVersion() != "" {
			if !reconcileObject.IsEndTimeSet() {
				continue
			}
			predecessor := getPredecessor(reconcileObject, items)
			if predecessor == nil {
				continue
			}

			previousInterval := reconcileObject.GetEndTime().Sub(predecessor.GetStartTime())
			gauge.Observe(ctx, previousInterval.Seconds(), reconcileObject.GetDurationMetricsAttributes()...)
		}
	}

	return nil
}

func getPredecessor(successor *interfaces.MetricsObjectWrapper, items []client.Object) interfaces.MetricsObject {
	var predecessor interfaces.MetricsObject
	for i := 0; i < len(items); i++ {
		// to calculate the interval, we take the earliest revision of the previous version as a reference
		if strings.HasPrefix(items[i].GetName(), fmt.Sprintf("%s-%s", successor.GetParentName(), successor.GetPreviousVersion())) {
			predecessorCandidate, err := interfaces.NewMetricsObjectWrapperFromClientObject(items[i])
			if err != nil {
				// continue with the other items
				continue
			}
			if predecessor == nil || predecessorCandidate.GetStartTime().Before(predecessor.GetStartTime()) {
				predecessor = predecessorCandidate
			}
		}
	}
	return predecessor
}

func ObserveActiveInstances(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge asyncint64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	for _, ro := range piWrapper.GetItems() {
		activeMetricsObject, _ := interfaces.NewActiveMetricsObjectWrapperFromClientObject(ro)
		gaugeValue := int64(0)
		if !activeMetricsObject.IsEndTimeSet() {
			gaugeValue = int64(1)
		}

		gauge.Observe(ctx, gaugeValue, activeMetricsObject.GetActiveMetricsAttributes()...)
	}

	return nil
}
