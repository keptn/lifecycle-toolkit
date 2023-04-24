package options

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-logr/logr"
	optionsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/options/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestKeptnConfigReconciler_Reconcile(t *testing.T) {
	reconciler := setupReconciler()

	// set up logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	type args struct {
		ctx context.Context
		req ctrl.Request
	}
	tests := []struct {
		name              string
		args              args
		lastAppliedConfig *optionsv1alpha1.KeptnConfigSpec
		want              ctrl.Result
		wantErr           bool
	}{
		{
			name: "test 1",
			args: args{
				ctx: context.TODO(),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "keptn-lifecycle-toolkit-system",
						Name:      "empty-config",
					},
				},
			},
			lastAppliedConfig: &optionsv1alpha1.KeptnConfigSpec{},
			want:              ctrl.Result{},
			wantErr:           false,
		},
		{
			name: "test 2",
			args: args{
				ctx: context.TODO(),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "keptn-lifecycle-toolkit-system",
						Name:      "empty-config",
					},
				},
			},
			want:    ctrl.Result{},
			wantErr: false,
		},
		{
			name: "test 3",
			args: args{
				ctx: context.TODO(),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "keptn-lifecycle-toolkit-system",
						Name:      "not-found-config",
					},
				},
			},
			want:    ctrl.Result{},
			wantErr: false,
		},
		{
			name: "test 4",
			args: args{
				ctx: context.TODO(),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: "keptn-lifecycle-toolkit-system",
						Name:      "config1",
					},
				},
			},
			lastAppliedConfig: &optionsv1alpha1.KeptnConfigSpec{
				OTelCollectorUrl: "some-url",
			},
			want:    ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reconciler.LastAppliedSpec = tt.lastAppliedConfig
			got, err := reconciler.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptnConfigReconciler_initConfig(t *testing.T) {
	type fields struct {
		Client          client.Client
		Scheme          *runtime.Scheme
		Log             logr.Logger
		LastAppliedSpec *optionsv1alpha1.KeptnConfigSpec
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test basic initialization",
			fields: fields{
				LastAppliedSpec: &optionsv1alpha1.KeptnConfigSpec{
					OTelCollectorUrl: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &KeptnConfigReconciler{
				Client:          tt.fields.Client,
				Scheme:          tt.fields.Scheme,
				Log:             tt.fields.Log,
				LastAppliedSpec: tt.fields.LastAppliedSpec,
			}
			r.initConfig()
		})
	}
}

func TestKeptnConfigReconciler_reconcileOtelCollectorUrl(t *testing.T) {
	// set up logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	type fields struct {
		Client          client.Client
		Scheme          *runtime.Scheme
		Log             logr.Logger
		LastAppliedSpec *optionsv1alpha1.KeptnConfigSpec
	}
	type args struct {
		config *optionsv1alpha1.KeptnConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ctrl.Result
		wantErr bool
	}{
		{
			name: "Test garbage URL",
			fields: fields{
				Client: nil,
				Scheme: nil,
				Log:    ctrl.Log.WithName("test-keptn-config-controller"),
				LastAppliedSpec: &optionsv1alpha1.KeptnConfigSpec{
					OTelCollectorUrl: "",
				},
			},
			args: args{
				config: &optionsv1alpha1.KeptnConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-config",
					},
					Spec: optionsv1alpha1.KeptnConfigSpec{
						OTelCollectorUrl: "some-url",
					},
				},
			},
			want:    ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: true,
		},
		{
			name: "Test with no URL",
			fields: fields{
				Client: nil,
				Scheme: nil,
				Log:    ctrl.Log.WithName("test-keptn-config-controller"),
			},
			args: args{
				config: &optionsv1alpha1.KeptnConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-config",
					},
					Spec: optionsv1alpha1.KeptnConfigSpec{
						OTelCollectorUrl: "",
					},
				},
			},
			want:    ctrl.Result{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &KeptnConfigReconciler{
				Client:          tt.fields.Client,
				Scheme:          tt.fields.Scheme,
				Log:             tt.fields.Log,
				LastAppliedSpec: tt.fields.LastAppliedSpec,
			}
			got, err := r.reconcileOtelCollectorUrl(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("reconcileOtelCollectorUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconcileOtelCollectorUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupReconciler() *KeptnConfigReconciler {
	emptyConfig := &optionsv1alpha1.KeptnConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "empty-config",
			Namespace: "keptn-lifecycle-toolkit-system",
		},
		Spec: optionsv1alpha1.KeptnConfigSpec{
			OTelCollectorUrl: "",
		},
	}
	config1 := &optionsv1alpha1.KeptnConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config1",
			Namespace: "keptn-lifecycle-toolkit-system",
		},
		Spec: optionsv1alpha1.KeptnConfigSpec{
			OTelCollectorUrl: "url1",
		},
	}
	config2 := &optionsv1alpha1.KeptnConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config2",
			Namespace: "keptn-lifecycle-toolkit-system",
		},
		Spec: optionsv1alpha1.KeptnConfigSpec{
			OTelCollectorUrl: "url2",
		},
	}

	//setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	fakeClient := fake.NewClient(emptyConfig, config1, config2)

	r := &KeptnConfigReconciler{
		Client: fakeClient,
		Scheme: fakeClient.Scheme(),
		Log:    ctrl.Log.WithName("test-keptnconfig-controller"),
	}
	return r
}
