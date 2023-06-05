package taskdefinition_test

import (
	"context"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
)

var _ = Describe("Taskdefinition", Ordered, func() {
	var (
		taskDefinitionName string
		namespace          string
	)

	BeforeEach(func() { // list var here they will be copied for every spec
		taskDefinitionName = names.SimpleNameGenerator.GenerateName("my-taskdefinition-")
		namespace = "default" // namespaces are not deleted in the api so be careful
	})

	Describe("Creation of a TaskDefinition", func() {
		var (
			taskDefinition *klcv1alpha3.KeptnTaskDefinition
			configmap      *v1.ConfigMap
		)
		Context("Reconcile TaskDefinition", func() {
			BeforeEach(func() {
			})

			It("create ConfigMap from inline function", func() {
				By("Create TaskDefinition")
				taskDefinition = &klcv1alpha3.KeptnTaskDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name:      taskDefinitionName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
						Function: &klcv1alpha3.RuntimeSpec{
							Inline: klcv1alpha3.Inline{
								Code: "console.log(Hello);",
							},
						},
					},
				}

				err := k8sClient.Create(context.TODO(), taskDefinition)
				Expect(err).To(BeNil())

				By("Check if ConfigMap was created")

				configmap = &v1.ConfigMap{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      "keptnfn-" + taskDefinitionName,
					}, configmap)
					g.Expect(err).To(BeNil())
					g.Expect(configmap.Data["code"]).To(Equal(taskDefinition.Spec.Function.Inline.Code))

				}, "30s").Should(Succeed())

				By("Check if TaskDefinition was updated")

				taskDefinition2 := &klcv1alpha3.KeptnTaskDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      taskDefinition.Name,
					}, taskDefinition2)
					g.Expect(err).To(BeNil())
					g.Expect(taskDefinition2.Status.Function.ConfigMap).To(Equal(configmap.Name))

				}, "30s").Should(Succeed())

				err = k8sClient.Delete(context.TODO(), configmap)
				common.LogErrorIfPresent(err)
			})

			It("TaskDefinition referencing existing Configmap defaulting to Deno", func() {
				By("Create TaskDefinition")
				taskDefinition = &klcv1alpha3.KeptnTaskDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name:      taskDefinitionName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
						Function: &klcv1alpha3.RuntimeSpec{
							ConfigMapReference: klcv1alpha3.ConfigMapReference{
								Name: "my-configmap",
							},
						},
					},
				}

				err := k8sClient.Create(context.TODO(), taskDefinition)
				Expect(err).To(BeNil())

				By("Create ConfigMap")

				configmap = &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-configmap",
						Namespace: namespace,
					},
					Data: map[string]string{
						"code": "console.log(Hello);",
					},
				}

				err = k8sClient.Create(context.TODO(), configmap)
				Expect(err).To(BeNil())

				By("Check if TaskDefinition was updated")

				taskDefinition2 := &klcv1alpha3.KeptnTaskDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      taskDefinition.Name,
					}, taskDefinition2)
					g.Expect(err).To(BeNil())
					g.Expect(taskDefinition2.Status.Function.ConfigMap).To(Equal(configmap.Name))

				}, "40s").Should(Succeed())

				err = k8sClient.Delete(context.TODO(), configmap)
				common.LogErrorIfPresent(err)
			})

			It("TaskDefinition referencing existing Configmap for Python runtime", func() {
				By("Create TaskDefinition")
				taskDefinition = &klcv1alpha3.KeptnTaskDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name:      taskDefinitionName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
						Python: &klcv1alpha3.RuntimeSpec{
							ConfigMapReference: klcv1alpha3.ConfigMapReference{
								Name: "my-configmap",
							},
						},
					},
				}

				err := k8sClient.Create(context.TODO(), taskDefinition)
				Expect(err).To(BeNil())

				By("Create ConfigMap")

				configmap = &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-configmap",
						Namespace: namespace,
					},
					Data: map[string]string{
						"code": "print(\"Hello, World!\")",
					},
				}

				err = k8sClient.Create(context.TODO(), configmap)
				Expect(err).To(BeNil())

				By("Check if TaskDefinition was updated")

				taskDefinition2 := &klcv1alpha3.KeptnTaskDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      taskDefinition.Name,
					}, taskDefinition2)
					g.Expect(err).To(BeNil())
					g.Expect(taskDefinition2.Status.Function.ConfigMap).To(Equal(configmap.Name))

				}, "30s").Should(Succeed())

				err = k8sClient.Delete(context.TODO(), configmap)
				common.LogErrorIfPresent(err)
			})

			It("TaskDefinition referencing non-existing Configmap", func() {
				By("Create TaskDefinition")
				taskDefinition = &klcv1alpha3.KeptnTaskDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name:      taskDefinitionName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
						Function: &klcv1alpha3.RuntimeSpec{
							ConfigMapReference: klcv1alpha3.ConfigMapReference{
								Name: "my-configmap-non-existing",
							},
						},
					},
				}

				err := k8sClient.Create(context.TODO(), taskDefinition)
				Expect(err).To(BeNil())

				By("Check that ConfigMap does not exists")

				configmap = &v1.ConfigMap{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      "my-configmap-non-existing",
					}, configmap)
					g.Expect(err).NotTo(BeNil())
				}, "30s").Should(Succeed())

				By("Check that TaskDefinition Status was not updated")

				taskDefinition2 := &klcv1alpha3.KeptnTaskDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      taskDefinition.Name,
					}, taskDefinition2)
					g.Expect(err).To(BeNil())
					g.Expect(taskDefinition2.Status.Function.ConfigMap).NotTo(Equal("my-configmap-non-existing"))

				}, "30s").Should(Succeed())

			})

			AfterEach(func() {
				err := k8sClient.Delete(context.TODO(), taskDefinition)
				common.LogErrorIfPresent(err)

			})
		})
	})
})
