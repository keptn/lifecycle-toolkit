package e2e

import (
	"time"

	testv1alpha2 "github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha2"
	common2 "github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha2/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

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

var SchedulingError = errors.New("Pod is not scheduled nor existing, this tests works only on a real installation have you setup your kind env?")
var SchedulingInPending = errors.New("Pod is pending")

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
			pod              *apiv1.Pod
			workloadinstance *testv1alpha2.KeptnWorkloadInstance
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				err := k8sClient.Delete(ctx, pod)
				logErrorIfPresent(err)
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")

			pod = WithContainer(st.MakePod().
				Namespace("default").
				Name(name).
				Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "5"}).
				ZeroTerminationGracePeriod().
				Obj(), pause)
			pod.Annotations = annotations
			pod.Spec.SchedulerName = "keptn-scheduler"

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
				checkWorkload(workloadinstance, "myapp-myworkload-1.0.0", *pod, "Succeeded")
			})
		})
	})

	Describe("If NOT annotated or labeled for keptn-scheduler", func() {
		pause := imageutils.GetPauseImageName()
		var (
			pod *apiv1.Pod
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				err := k8sClient.Delete(ctx, pod)
				logErrorIfPresent(err)
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
					return podScheduled(pod.Namespace, pod.Name)
				}).WithTimeout(time.Second * 30).WithPolling(3 * time.Second).Should(Succeed())

			})
		})
	})

	Describe("If labeled for keptn-scheduler", func() {
		labels := map[string]string{
			WorkloadAnnotation: "myworkload",
			VersionAnnotation:  "1.0.1",
			AppAnnotation:      "mylabeledapp",
		}

		pause := imageutils.GetPauseImageName()
		var (
			pod              *apiv1.Pod
			workloadinstance *testv1alpha2.KeptnWorkloadInstance
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				err := k8sClient.Delete(ctx, pod)
				logErrorIfPresent(err)
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")

			pod = WithContainer(st.MakePod().
				Namespace("default").
				Name(name).
				Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "50"}).
				ZeroTerminationGracePeriod().
				Obj(), pause)
			pod.Labels = labels
			pod.Spec.SchedulerName = "keptn-scheduler"

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

				checkWorkload(workloadinstance, "mylabeledapp-myworkload-1.0.1", *pod, "Succeeded")
			})

		})
	})

	Describe("If labeled for keptn-scheduler", func() {
		labels := map[string]string{
			K8sRecommendedWorkloadAnnotations: "myworkload",
			K8sRecommendedVersionAnnotations:  "1.0.2",
			K8sRecommendedAppAnnotations:      "mylabeledapp",
		}

		pause := imageutils.GetPauseImageName()
		var (
			pod              *apiv1.Pod
			workloadinstance *testv1alpha2.KeptnWorkloadInstance
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				err := k8sClient.Delete(ctx, pod)
				logErrorIfPresent(err)
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")

			pod = WithContainer(st.MakePod().
				Namespace("default").
				Name(name).
				Req(map[apiv1.ResourceName]string{apiv1.ResourceMemory: "50"}).
				ZeroTerminationGracePeriod().
				Obj(), pause)
			pod.Labels = labels
			pod.Spec.SchedulerName = "keptn-scheduler"

			err := k8sClient.Create(ctx, pod)
			Expect(ignoreAlreadyExists(err)).NotTo(HaveOccurred(), "could not add pod")
		})

		Context("a new Pod", func() {

			It(" should NOT be scheduled when workload instance pre-evaluation checks fails", func() {
				checkWorkload(workloadinstance, "mylabeledapp-myworkload-1.0.2", *pod, "Failed")
			})
		})

	})
})

func checkWorkload(workloadinstance *testv1alpha2.KeptnWorkloadInstance, workloadname string, pod apiv1.Pod, status common2.KeptnState) {
	workloadinstance = initWorkloadInstance(workloadname)

	err := k8sClient.Create(ctx, workloadinstance)
	Expect(ignoreAlreadyExists(err)).To(BeNil())

	Eventually(func() error {
		err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: workloadname}, workloadinstance)
		return err
	}).Should(Succeed())
	workloadinstance.Status.PreDeploymentEvaluationStatus = status
	err = k8sClient.Status().Update(ctx, workloadinstance)

	Expect(err).To(BeNil())
	assertion := Eventually(func() error {
		return podScheduled(pod.Namespace, pod.Name)
	}).WithTimeout(time.Second * 20).WithPolling(3 * time.Second)

	if status == "Failed" {
		assertion.ShouldNot(Succeed())
	} else {
		assertion.Should(Succeed())
	}

	err = k8sClient.Delete(ctx, workloadinstance)
	Expect(err).NotTo(HaveOccurred(), "could not remove workloadinstance")
}

func initWorkloadInstance(name string) *testv1alpha2.KeptnWorkloadInstance {

	var fakeInstance = testv1alpha2.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnWorkloadInstance",
			APIVersion: "lifecycle.keptn.sh/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Spec: testv1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: testv1alpha2.KeptnWorkloadSpec{
				ResourceReference: testv1alpha2.ResourceReference{Name: "myfakeres"},
			},
		},
		Status: testv1alpha2.KeptnWorkloadInstanceStatus{},
	}

	return &fakeInstance
}

func podScheduled(namespace, name string) error {
	pod := apiv1.Pod{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &pod)
	if err != nil {
		// This could be a connection error we want to retry.
		GinkgoLogr.Error(err, "Failed to get", "pod", klog.KRef(namespace, name))
		return err
	}

	if pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodFailed || pod.Status.Phase == apiv1.PodRunning {
		return nil
	}

	for _, c := range pod.Status.Conditions {
		if c.Type == apiv1.PodScheduled {
			if c.Status == apiv1.ConditionTrue {
				return nil
			}
			return SchedulingInPending
		}
	}
	return SchedulingError
}

func WithContainer(pod *apiv1.Pod, image string) *apiv1.Pod {
	pod.Spec.Containers[0].Name = "web"
	pod.Spec.Containers[0].Image = image
	return pod
}
