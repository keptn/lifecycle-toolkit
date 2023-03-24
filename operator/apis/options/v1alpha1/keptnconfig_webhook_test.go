package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeClient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestValidateCreate(t *testing.T) {
	tests := []struct {
		name    string
		conf    KeptnConfig
		args    []KeptnConfig
		wantErr bool
	}{
		{
			name:    "empty",
			conf:    KeptnConfig{},
			args:    []KeptnConfig{},
			wantErr: false,
		},
		{
			name: "same name different namespace",
			conf: KeptnConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "keptnconf",
					Namespace: "klt",
				},
			},
			args: []KeptnConfig{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "keptnconf",
						Namespace: "klt2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "same namespace different name",
			conf: KeptnConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "keptnconf",
					Namespace: "klt",
				},
			},
			args: []KeptnConfig{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "keptnconf2",
						Namespace: "klt",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple",
			conf: KeptnConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "keptnconf",
					Namespace: "klt",
				},
			},
			args: []KeptnConfig{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "keptnconf",
						Namespace: "klt",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "keptnconf-fail",
						Namespace: "klt",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ns = tt.conf.Namespace
			_ = AddToScheme(scheme.Scheme)
			objs := []client.Object{}
			for _, arg := range tt.args {
				c := arg.DeepCopy()
				objs = append(objs, c)
			}
			fc := newFakeReadClient(objs...)
			_client = fc
			err := tt.conf.ValidateCreate()
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func newFakeReadClient(objs ...client.Object) client.Client {
	_ = AddToScheme(scheme.Scheme)
	return fakeClient.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(objs...).Build()
}

func TestValidateUpdateAndDelete(t *testing.T) {
	c := KeptnConfig{}
	require.Nil(t, c.ValidateUpdate(nil))
	require.Nil(t, c.ValidateDelete())
}
