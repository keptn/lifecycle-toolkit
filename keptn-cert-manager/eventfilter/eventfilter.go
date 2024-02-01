package eventfilter

import (
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func ForLabelsAndNamespace(labels labels.Selector, namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return isInNamespace(object, namespace) && matchesLabels(object, labels)
	})
}

func ForNamesAndNamespace(names []string, namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return isInNamespace(object, namespace) && matchesName(object, names)
	})
}

func matchesLabels(object client.Object, selector labels.Selector) bool {
	return selector.Matches(labels.Set(object.GetLabels()))
}

func isInNamespace(object client.Object, namespace string) bool {
	return object.GetNamespace() == namespace
}

func matchesName(object client.Object, names []string) bool {
	for _, name := range names {
		if object.GetName() == name {
			return true
		}
	}
	return false
}
