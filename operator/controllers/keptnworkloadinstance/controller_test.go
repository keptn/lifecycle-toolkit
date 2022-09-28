package keptnworkloadinstance

import (
	"context"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestKeptnWorkloadInstanceReconciler_IsPodRunning(t *testing.T) {
	p1 := makeNominatedPod("pod1", "node1", v1.PodRunning)
	p2 := makeNominatedPod("pod2", "node1", v1.PodPending)
	podList := &v1.PodList{Items: []v1.Pod{p1, p2}}
	podList2 := &v1.PodList{Items: []v1.Pod{p2}}
	r := &KeptnWorkloadInstanceReconciler{
		Client: fake.NewClientBuilder().WithLists(podList).Build(),
	}
	if !r.IsPodRunning(context.TODO(), v1alpha1.ResourceReference{UID: types.UID("pod1")}, "node1") {
		t.Errorf("Wrong!")
	}

	r2 := &KeptnWorkloadInstanceReconciler{
		Client: fake.NewClientBuilder().WithLists(podList2).Build(),
	}
	if r2.IsPodRunning(context.TODO(), v1alpha1.ResourceReference{UID: types.UID("pod1")}, "node1") {
		t.Errorf("Wrong!")
	}

}

func makeNominatedPod(podName string, nodeName string, phase v1.PodPhase) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nodeName,
			Name:      podName,
			UID:       types.UID(podName),
		},
		Status: v1.PodStatus{
			Phase:             phase,
			NominatedNodeName: nodeName,
		},
	}
}
