package common

import (
	"context"
	lfcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func AddApp(c client.Client, name string) error {
	app := &lfcv1alpha1.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: lfcv1alpha1.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lfcv1alpha1.KeptnAppStatus{},
	}
	return c.Create(context.TODO(), app)

}

func AddAppVersion(c client.Client, name string, status lfcv1alpha1.KeptnAppVersionStatus) error {
	app := &lfcv1alpha1.KeptnAppVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec:   lfcv1alpha1.KeptnAppVersionSpec{KeptnAppSpec: lfcv1alpha1.KeptnAppSpec{Version: "1.0.0"}},
		Status: status,
	}
	return c.Create(context.TODO(), app)

}
