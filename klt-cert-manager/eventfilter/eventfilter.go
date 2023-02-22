package eventfilter

import (
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func ForObjectNameAndNamespace(name string, namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return isInNamespace(object, namespace) && hasName(object, name)
	})
}

func ForLabelsAndNamespace(labels labels.Selector, namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return isInNamespace(object, namespace) && matchesLabels(object, labels)
	})
}

func matchesLabels(object client.Object, selector labels.Selector) bool {
	return selector.Matches(labels.Set(object.GetLabels()))
}

func hasName(object client.Object, name string) bool {
	return object.GetName() == name
}

func isInNamespace(object client.Object, namespace string) bool {
	return object.GetNamespace() == namespace
}
