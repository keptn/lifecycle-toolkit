package e2e

import (
	testv1alpha1 "github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/klog/v2"
	st "k8s.io/kubernetes/pkg/scheduler/testing"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const K8sRecommendedWorkloadAnnotations = "app.kubernetes.io/name"
const K8sRecommendedVersionAnnotations = "app.kubernetes.io/version"
const K8sRecommendedAppAnnotations = "app.kubernetes.io/part-of"

// clean example of E2E test/ integration test --

var _ = Describe("[E2E] KeptnScheduler", Ordered, func() {
	BeforeAll(func() {
		wg.Add(1) //this tells the suite that all test have finished
	})
	AfterAll(func() {
		wg.Done() //this tells the suite that all test have finished
	})
	Describe("If annotated for keptn-scheduler", func() {
		annotations := map[string]string{
			WorkloadAnnotation: "myworkload",
			VersionAnnotation:  "1.0.0",
			AppAnnotation:      "myapp",
		}

		pause := imageutils.GetPauseImageName()
		var (
			workloadinstance *testv1alpha1.KeptnWorkloadInstance
			pod              *apiv1.Pod
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				k8sClient.Delete(ctx, pod)
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")

			pod = WithContainer(st.MakePod().
				Namespace("default").
				Name(name).
				Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "50"}).
				ZeroTerminationGracePeriod().
				Obj(), pause)
			pod.Annotations = annotations

			err := k8sClient.Create(ctx, pod)
			Expect(ignoreAlreadyExists(err)).NotTo(HaveOccurred(), "could not add pod")
		})

		Context("a new Pod", func() {

			It(" should stay pending if no workload instance is available", func() {

				newPod := &apiv1.Pod{}
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, newPod)
					return err
				}).Should(Succeed())

				Expect(newPod.Status.Phase).To(Equal(apiv1.PodPending))

			})

			It(" should be scheduled when workload instance pre-evaluation checks are done", func() {

				workloadinstance = initWorkloadInstance()

				err := k8sClient.Create(ctx, workloadinstance)
				Expect(ignoreAlreadyExists(err)).To(BeNil())

				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: "myapp-myworkload-1.0.0"}, workloadinstance)
					return err
				}).Should(Succeed())
				workloadinstance.Status.PreDeploymentEvaluationStatus = "Succeeded"
				err = k8sClient.Status().Update(ctx, workloadinstance)

				Expect(err).To(BeNil())
				Eventually(func() error {
					return podRunning(pod.Namespace, pod.Name)
				}).Should(Succeed())

				err = k8sClient.Delete(ctx, workloadinstance)
				Expect(err).NotTo(HaveOccurred(), "could not remove workloadinstance")

			})
		})
	})

	Describe("If NOT annotated for keptn-scheduler", func() {
		pause := imageutils.GetPauseImageName()
		var (
			pod *apiv1.Pod
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				k8sClient.Delete(ctx, pod)
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")

			pod = WithContainer(st.MakePod().
				Namespace("default").
				Name(name).
				Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "50"}).
				ZeroTerminationGracePeriod().
				Obj(), pause)

			err := k8sClient.Create(ctx, pod)
			Expect(ignoreAlreadyExists(err)).NotTo(HaveOccurred(), "could not add pod")
		})

		Context("a new Pod", func() {

			It(" should be immediately scheduled", func() {

				Eventually(func() error {
					return podRunning(pod.Namespace, pod.Name)
				}).Should(Succeed())

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

	return &fakeInstance
}

func podRunning(namespace, name string) error {
	pod := apiv1.Pod{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &pod)
	if err != nil {
		// This could be a connection error so we want to retry.
		GinkgoLogr.Error(err, "Failed to get", "pod", klog.KRef(namespace, name))
		return err
	}

	for _, c := range pod.Status.Conditions {
		if c.Type == "PodScheduled" {
			return nil
		}
	}
	return err
}

func WithContainer(pod *apiv1.Pod, image string) *apiv1.Pod {
	pod.Spec.Containers[0].Name = "web"
	pod.Spec.Containers[0].Image = image
	return pod
}
