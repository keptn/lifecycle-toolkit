package test

import (
	"fmt"
	testv1alpha1 "github.com/keptn/lifecycle-toolkit/scheduler/test/fake/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	st "k8s.io/kubernetes/pkg/scheduler/testing"
	imageutils "k8s.io/kubernetes/test/utils/image"
	"time"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const K8sRecommendedWorkloadAnnotations = "app.kubernetes.io/name"
const K8sRecommendedVersionAnnotations = "app.kubernetes.io/version"
const K8sRecommendedAppAnnotations = "app.kubernetes.io/part-of"

// clean example of E2E test/ integration test --

var _ = Describe("KeptnScheduler", Ordered, func() {

	//create a test Pod

	annotations := map[string]string{
		WorkloadAnnotation: "myworkload",
		VersionAnnotation:  "1.0.0",
		AppAnnotation:      "myapp",
	}

	pause := imageutils.GetPauseImageName()

	pod := WithContainer(st.MakePod().
		Namespace("default").
		Name("mypod").
		Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "50"}).
		ZeroTerminationGracePeriod().
		Obj(), pause)
	pod.Annotations = annotations

	// Create Deployment
	Describe("Creation of a new Deployment annotated for keptn-scheduler", func() {

		Context("a new Deployment", func() {

			BeforeEach(func() {
				err := k8sClient.Create(ctx, pod)
				Expect(err).NotTo(HaveOccurred(), "could not add deployment")

				//example of creating a node
				//nodeName := "fake-node"
				//node = st.MakeNode().Name("fake-node").Label("node", nodeName).Obj()
				//node.Status.Allocatable = apiv1.ResourceList{
				//	apiv1.ResourcePods:   *resource.NewQuantity(32, resource.DecimalSI),
				//	apiv1.ResourceMemory: *resource.NewQuantity(300, resource.DecimalSI),
				//}
				//node.Status.Capacity = apiv1.ResourceList{
				//	apiv1.ResourcePods:   *resource.NewQuantity(32, resource.DecimalSI),
				//	apiv1.ResourceMemory: *resource.NewQuantity(300, resource.DecimalSI),
				//}
				//node, err = testCtx.ClientSet.CoreV1().Nodes().Create(testCtx.Ctx, node, metav1.CreateOptions{})
				//Expect(err).NotTo(HaveOccurred(), "Could not add node")

			})
			It(" should stay pending until workload instance is done", func() {
				newPod := &apiv1.Pod{}
				var err error
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, newPod)
					return err
				}).Should(Succeed())

				Expect(newPod.Status.Phase).To(Equal(apiv1.PodPending))
				workload := initWorkloadInstance()

				Expect(err).To(BeNil())
				err = k8sClient.Create(ctx, workload)
				Expect(err).To(BeNil())

				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: "myapp-myworkload-1.0.0"}, workload)
					return err
				}).Should(Succeed())
				workload.Status.PreDeploymentEvaluationStatus = "Succeeded"
				err = k8sClient.Status().Update(ctx, workload)

				Expect(err).To(BeNil())
				newPod = &apiv1.Pod{}
				err = wait.Poll(1*time.Second, 120*time.Second, func() (bool, error) {

					if !deploymentRunning(pod.Namespace, pod.Name) {
						return false, nil
					}
					return true, nil
				})

				Expect(err).NotTo(HaveOccurred())

			})
		})
	})
})

func initWorkloadInstance() *testv1alpha1.KeptnWorkloadInstance {

	var fakeInstance = testv1alpha1.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnWorkloadInstance",
			APIVersion: "lifecycle.keptn.sh/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "myapp-myworkload-1.0.0", Namespace: "default"},
		Status:     testv1alpha1.KeptnWorkloadInstanceStatus{PreDeploymentEvaluationStatus: "Succeeded"},
	}

	//	fmt.Println("known", target, ok)
	return &fakeInstance
}

func deploymentRunning(namespace, name string) bool {
	pod := apiv1.Pod{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &pod)
	if err != nil {
		// This could be a connection error so we want to retry.
		GinkgoLogr.Error(err, "Failed to get", "pod", klog.KRef(namespace, name))
		return false
	}
	fmt.Println(fmt.Sprintf("depl %+v", pod.Status.Phase))

	for _, c := range pod.Status.Conditions {
		if c.Type == "PodScheduled" {
			return c.Status == "True"
		}
	}
	return false
}

func WithContainer(pod *apiv1.Pod, image string) *apiv1.Pod {
	pod.Spec.Containers[0].Name = "web"
	pod.Spec.Containers[0].Image = "nginx:1.12"
	return pod
}
