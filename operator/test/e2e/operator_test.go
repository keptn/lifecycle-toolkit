package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"

var _ = Describe("[E2E] KeptnOperator", Ordered, func() {
	BeforeAll(func() {
		wg.Add(1) //this tells the suite that all test have finished
	})
	AfterAll(func() {
		wg.Done() //this tells the suite that all test have finished
	})
	Describe("If annotated for keptn", func() {
		annotations := map[string]string{
			WorkloadAnnotation: "myworkload",
			VersionAnnotation:  "1.0.0",
			AppAnnotation:      "myapp",
		}

		var (
			pod    *apiv1.Pod
			newPod *apiv1.Pod
		)
		BeforeEach(func() {
			DeferCleanup(func() {
				err := k8sClient.Delete(ctx, pod)
				Expect(err).NotTo(HaveOccurred(), "could not remove pod")
			})

			//create a test Pod
			name := names.SimpleNameGenerator.GenerateName("my-testpod-")
			pod = &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Annotations: annotations,
					Namespace:   "default",
				},
				Spec: apiv1.PodSpec{
					SchedulerName: "",
					Containers: []apiv1.Container{
						{
							Name:  "mybusy",
							Image: "busybox:1.32.1",
						},
					},
				},
			}

			err := k8sClient.Create(ctx, pod)
			Expect(ignoreAlreadyExists(err)).NotTo(HaveOccurred(), "could not add pod")
		})

		Context("a new Pod", func() {

			It(" should stay pending", func() {

				newPod = &apiv1.Pod{}
				Eventually(func() error {
					err := k8sClient.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, newPod)
					return err
				}).Should(Succeed())

				Expect(newPod.Status.Phase).To(Equal(apiv1.PodPending))

			})

			It(" should be assigned to keptn scheduler", func() {
				Expect(newPod.Spec.SchedulerName == "keptn-scheduler")
			})
		})
	})
})
