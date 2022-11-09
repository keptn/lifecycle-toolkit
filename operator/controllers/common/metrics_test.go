package common

import (
	"context"
	"testing"
	"time"

	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMetrics_GetDeploymentDuration(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
		result        []apicommon.GaugeFloatValue
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha1.KeptnAppList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppList{},
			err:           ErrCannotWrapToListItem,
			result:        nil,
		},
		{
			name: "no endtime set",
			list: &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						Status: lifecyclev1alpha1.KeptnAppVersionStatus{},
					},
				},
			},
			err:    nil,
			result: []apicommon.GaugeFloatValue{},
		},
		{
			name: "endtime set",
			list: &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha1.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
			result: []apicommon.GaugeFloatValue{
				{
					Value: 5 * time.Second.Seconds(),
					Attributes: []attribute.KeyValue{
						common.AppName.String("appName"),
						common.AppVersion.String("version"),
						common.AppPreviousVersion.String("previousVersion"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lifecyclev1alpha1.AddToScheme(scheme.Scheme)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			res, err := GetDeploymentDuration(context.TODO(), client, tt.list)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.result, res)
		})

	}
}

func TestMetrics_GetActiveInstances(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
		result        []apicommon.GaugeValue
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha1.KeptnAppList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppList{},
			err:           ErrCannotWrapToListItem,
			result:        nil,
		},
		{
			name: "no endtime set - active instances",
			list: &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha1.KeptnAppVersionStatus{},
					},
				},
			},
			err: nil,
			result: []apicommon.GaugeValue{
				{
					Value: 1,
					Attributes: []attribute.KeyValue{
						common.AppName.String("appName"),
						common.AppVersion.String("version"),
						common.AppNamespace.String("namespace"),
					},
				},
			},
		},
		{
			name: "endtime set - no active instances",
			list: &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha1.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
			result: []apicommon.GaugeValue{
				{
					Value: 0,
					Attributes: []attribute.KeyValue{
						common.AppName.String("appName"),
						common.AppVersion.String("version"),
						common.AppNamespace.String("namespace"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lifecyclev1alpha1.AddToScheme(scheme.Scheme)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			res, err := GetActiveInstances(context.TODO(), client, tt.list)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.result, res)
		})

	}
}

func TestMetrics_GetDeploymentInterval(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		clientObject  client.Object
		list          client.ObjectList
		previous      client.Object
		err           error
		result        []apicommon.GaugeFloatValue
	}{
		{
			name:          "failed to create wrapper",
			list:          &lifecyclev1alpha1.KeptnAppList{},
			clientObjects: &lifecyclev1alpha1.KeptnAppList{},
			clientObject:  &lifecyclev1alpha1.KeptnApp{},
			err:           ErrCannotWrapToListItem,
			result:        nil,
		},
		{
			name:         "no previous version",
			list:         &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObject: &lifecyclev1alpha1.KeptnApp{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
					},
				},
			},
			err:    nil,
			result: []apicommon.GaugeFloatValue{},
		},
		{
			name:         "previous version - no previous object",
			list:         &lifecyclev1alpha1.KeptnAppVersionList{},
			clientObject: &lifecyclev1alpha1.KeptnApp{},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
				},
			},
			err:    nil,
			result: nil,
		},
		{
			name:     "previous version - object found but cannot unwrap",
			list:     &lifecyclev1alpha1.KeptnAppVersionList{},
			previous: &lifecyclev1alpha1.KeptnApp{},
			clientObject: &lifecyclev1alpha1.KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "appName-previousVersion",
					Namespace: "namespace",
				},
			},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
				},
			},
			err:    ErrCannotWrapToPhaseItem,
			result: nil,
		},
		{
			name:     "previous version - object found but no endtime",
			list:     &lifecyclev1alpha1.KeptnAppVersionList{},
			previous: &lifecyclev1alpha1.KeptnAppVersion{},
			clientObject: &lifecyclev1alpha1.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "appName-previousVersion",
					Namespace: "namespace",
				},
				Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
						Version: "previousVersion",
					},
					AppName:         "appName",
					PreviousVersion: "",
				},
			},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
				},
			},
			err:    nil,
			result: []apicommon.GaugeFloatValue{},
		},
		{
			name:     "previous version - object found with endtime",
			list:     &lifecyclev1alpha1.KeptnAppVersionList{},
			previous: &lifecyclev1alpha1.KeptnAppVersion{},
			clientObject: &lifecyclev1alpha1.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "appName-previousVersion",
					Namespace: "namespace",
				},
				Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
						Version: "previousVersion",
					},
					AppName:         "appName",
					PreviousVersion: "",
				},
				Status: lifecyclev1alpha1.KeptnAppVersionStatus{
					EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
					StartTime: metav1.Time{Time: metav1.Now().Time},
				},
			},
			clientObjects: &lifecyclev1alpha1.KeptnAppVersionList{
				Items: []lifecyclev1alpha1.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: lifecyclev1alpha1.KeptnAppVersionSpec{
							KeptnAppSpec: lifecyclev1alpha1.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: lifecyclev1alpha1.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
			result: []apicommon.GaugeFloatValue{
				{
					Value: 5,
					Attributes: []attribute.KeyValue{
						common.AppName.String("appName"),
						common.AppVersion.String("version"),
						common.AppPreviousVersion.String("previousVersion"),
					},
				},
				{
					Value: 5,
					Attributes: []attribute.KeyValue{
						common.AppName.String("appName"),
						common.AppVersion.String("version"),
						common.AppPreviousVersion.String("previousVersion"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lifecyclev1alpha1.AddToScheme(scheme.Scheme)
			client := fake.NewClientBuilder().WithObjects(tt.clientObject).WithLists(tt.clientObjects).Build()
			res, err := GetDeploymentInterval(context.TODO(), client, tt.list, tt.previous)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.result, res)
		})

	}
}
