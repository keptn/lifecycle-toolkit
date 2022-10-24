package keptnapp

import (
	"context"
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/oteltest"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllertest"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"testing"
)

//simple UNIT TESTS

func TestKeptnAppReconciler_createAppVersion(t *testing.T) {

	tests := []struct {
		name    string
		app     *klcv1alpha1.KeptnApp
		want    *klcv1alpha1.KeptnAppVersion
		wantErr bool
	}{
		{
			name: "pass all field ",
			app: &klcv1alpha1.KeptnApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "keptnapps",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myapp",
					Namespace:   "default",
					UID:         "123",
					Labels:      nil,
					Annotations: nil,
				},
				Spec:   klcv1alpha1.KeptnAppSpec{},
				Status: klcv1alpha1.KeptnAppStatus{},
			},
			want: &klcv1alpha1.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myapp",
					Namespace: "default",
				},
				Spec: klcv1alpha1.KeptnAppVersionSpec{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &KeptnAppReconciler{
				Log:    GinkgoLogr,
				Tracer: oteltest.DefaultTracer(),
			}
			got, err := r.createAppVersion(context.TODO(), tt.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAppVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createAppVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//component test using fake api, or integration test if it uses real api
var _ = Describe("KeptnApp controller", func() {
	var reconciled chan reconcile.Request
	ctx := context.Background()

	BeforeEach(func() {
		reconciled = make(chan reconcile.Request)
		Expect(cfg).NotTo(BeNil())
	})

	Describe("controller", func() {
		// TODO(directxman12): write a whole suite of controller-client interaction tests

		It("should reconcile", func() {
			By("Creating the Manager")
			cm, err := manager.New(cfg, manager.Options{})
			Expect(err).NotTo(HaveOccurred())

			By("Creating the Controller")
			instance, err := controller.New("foo-controller", cm, controller.Options{
				Reconciler: reconcile.Func(
					func(_ context.Context, request reconcile.Request) (reconcile.Result, error) {
						reconciled <- request
						return reconcile.Result{}, nil
					}),
			})
			Expect(err).NotTo(HaveOccurred())

			By("Watching Resources")
			err = instance.Watch(&source.Kind{Type: &appsv1.ReplicaSet{}}, &handler.EnqueueRequestForOwner{
				OwnerType: &appsv1.Deployment{},
			})
			Expect(err).NotTo(HaveOccurred())

			err = instance.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{})
			Expect(err).NotTo(HaveOccurred())

			err = cm.GetClient().Get(ctx, types.NamespacedName{Name: "foo"}, &corev1.Namespace{})
			Expect(err).To(Equal(&cache.ErrCacheNotStarted{}))
			err = cm.GetClient().List(ctx, &corev1.NamespaceList{})
			Expect(err).To(Equal(&cache.ErrCacheNotStarted{}))

			By("Starting the Manager")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				defer GinkgoRecover()
				Expect(cm.Start(ctx)).NotTo(HaveOccurred())
			}()

			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "deployment-name"},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"foo": "bar"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"foo": "bar"}},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx",
									SecurityContext: &corev1.SecurityContext{
										Privileged: truePtr(),
									},
								},
							},
						},
					},
				},
			}
			expectedReconcileRequest := reconcile.Request{NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "deployment-name",
			}}

			By("Invoking Reconciling for Create")
			//deployment, err = k8sClient.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
			//Expect(err).NotTo(HaveOccurred())
			Expect(<-reconciled).To(Equal(expectedReconcileRequest))

			By("Invoking Reconciling for Update")
			newDeployment := deployment.DeepCopy()
			newDeployment.Labels = map[string]string{"foo": "bar"}
			//_, err = k8sClient.AppsV1().Deployments("default").Update(ctx, newDeployment, metav1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(<-reconciled).To(Equal(expectedReconcileRequest))

			By("Invoking Reconciling for an OwnedObject when it is created")
			//replicaset := &appsv1.ReplicaSet{
			//	ObjectMeta: metav1.ObjectMeta{
			//		Name: "rs-name",
			//		OwnerReferences: []metav1.OwnerReference{
			//			*metav1.NewControllerRef(deployment, schema.GroupVersionKind{
			//				Group:   "apps",
			//				Version: "v1",
			//				Kind:    "Deployment",
			//			}),
			//		},
			//	},
			//	Spec: appsv1.ReplicaSetSpec{
			//		Selector: &metav1.LabelSelector{
			//			MatchLabels: map[string]string{"foo": "bar"},
			//		},
			//		Template: deployment.Spec.Template,
			//	},
			//}

			//replicaset, err = k8sClient.Create(ctx, replicaset, metav1.CreateOptions{})
			//Expect(err).NotTo(HaveOccurred())
			//Expect(<-reconciled).To(Equal(expectedReconcileRequest))
			//
			//By("Invoking Reconciling for an OwnedObject when it is updated")
			//newReplicaset := replicaset.DeepCopy()
			//newReplicaset.Labels = map[string]string{"foo": "bar"}
			//_, err = k8sClient.AppsV1().ReplicaSets("default").Update(ctx, newReplicaset, metav1.UpdateOptions{})
			//Expect(err).NotTo(HaveOccurred())
			//Expect(<-reconciled).To(Equal(expectedReconcileRequest))
			//
			//By("Invoking Reconciling for an OwnedObject when it is deleted")
			//err = k8sClient.AppsV1().ReplicaSets("default").Delete(ctx, replicaset.Name, metav1.DeleteOptions{})
			//Expect(err).NotTo(HaveOccurred())
			//Expect(<-reconciled).To(Equal(expectedReconcileRequest))
			//
			//By("Invoking Reconciling for Delete")
			//err = k8sClient.AppsV1().Deployments("default").
			//	Delete(ctx, "deployment-name", metav1.DeleteOptions{})
			//Expect(err).NotTo(HaveOccurred())
			//Expect(<-reconciled).To(Equal(expectedReconcileRequest))

			By("Listing a type with a slice of pointers as items field")
			err = cm.GetClient().
				List(context.Background(), &controllertest.UnconventionalListTypeList{})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func truePtr() *bool {
	t := true
	return &t
}
