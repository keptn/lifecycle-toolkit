package common

import (
	"context"
	"testing"
	"time"

	lifecyclev1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/stretchr/testify/require"
	noop "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMetrics_ObserveDeploymentDuration(t *testing.T) {

	gauge, err := noop.NewNoopMeter().AsyncFloat64().Gauge("mine")
	require.Nil(t, err)

	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
		gauge         asyncfloat64.Gauge
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha2.KeptnAppList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
			gauge:         nil,
		},
		{
			name: "no endtime set",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{},
					},
				},
			},
			err:   nil,
			gauge: gauge,
		},
		{
			name: "endtime set",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err:   nil,
			gauge: gauge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecyclev1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			err = ObserveDeploymentDuration(context.TODO(), client, tt.list, gauge)
			require.ErrorIs(t, err, tt.err)
		})

	}
}

func TestMetrics_ObserveActiveInstances(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha2.KeptnAppList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
		},
		{
			name: "no endtime set - active instances",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{},
					},
				},
			},
			err: nil,
		},
		{
			name: "endtime set - no active instances",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecyclev1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			gauge, err := noop.NewNoopMeter().AsyncInt64().Gauge("mine")
			require.Nil(t, err)
			err = ObserveActiveInstances(context.TODO(), client, tt.list, gauge)
			require.ErrorIs(t, err, tt.err)

		})

	}
}

func TestMetrics_ObserveDeploymentInterval(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		previous      client.Object
		err           error
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha2.KeptnAppList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
		},
		{
			name: "no previous version",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "previous version - no previous object",
			list: &lifecyclev1alpha2.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "previous version - object found but no endtime",
			list:     &lifecyclev1alpha2.KeptnAppVersionList{},
			previous: &lifecyclev1alpha2.KeptnAppVersion{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-previousVersion",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "previous version - object found with endtime",
			list:     &lifecyclev1alpha2.KeptnAppVersionList{},
			previous: &lifecyclev1alpha2.KeptnAppVersion{},
			clientObjects: &lifecyclev1alpha2.KeptnAppVersionList{
				Items: []lifecyclev1alpha2.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-previousVersion",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
						Status: lifecyclev1alpha2.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
		},
	}

	gauge, err := noop.NewNoopMeter().AsyncFloat64().Gauge("mine")
	require.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecyclev1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			fakeClient := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			err = ObserveDeploymentInterval(context.TODO(), fakeClient, tt.list, gauge)
			require.ErrorIs(t, err, tt.err)
		})

	}
}

func TestGetPredecessor(t *testing.T) {
	now := time.Now()
	appVersions := &lifecyclev1alpha2.KeptnAppVersionList{
		Items: []lifecyclev1alpha2.KeptnAppVersion{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-app-1.0.0-1",
				},
				Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
						Version:  "1.0.0",
						Revision: 0,
					},
					AppName: "my-app",
				},
				Status: lifecyclev1alpha2.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-app-1.0.0-2",
				},
				Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
						Version:  "1.0.0",
						Revision: 0,
					},
					AppName: "my-app",
				},
				Status: lifecyclev1alpha2.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now.Add(1 * time.Second)),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-app-1.1.0-1",
				},
				Spec: lifecyclev1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: lifecyclev1alpha2.KeptnAppSpec{
						Version:  "1.0.0",
						Revision: 0,
					},
					AppName: "my-app",
				},
				Status: lifecyclev1alpha2.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
		},
	}

	appVersionsWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(appVersions)
	require.Nil(t, err)

	latestAppVersion, err := interfaces.NewMetricsObjectWrapperFromClientObject(appVersionsWrapper.GetItems()[2])

	require.Nil(t, err)
	predecessor := getPredecessor(latestAppVersion, appVersionsWrapper.GetItems())

	expectedPredecessor, err := interfaces.NewMetricsObjectWrapperFromClientObject(appVersionsWrapper.GetItems()[0])
	require.Nil(t, err)

	require.Equal(t, expectedPredecessor, predecessor)
}
