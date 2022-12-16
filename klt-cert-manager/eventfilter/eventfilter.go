package eventfilter

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func ForObjectNameAndNamespace(name string, namespace string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		return isInNamespace(object, namespace) && hasName(object, name)
	})
}

func hasName(object client.Object, name string) bool {
	return object.GetName() == name
}

func isInNamespace(object client.Object, namespace string) bool {
	return object.GetNamespace() == namespace
}
