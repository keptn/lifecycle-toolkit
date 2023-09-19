package v1alpha3

type ObjectReference struct {
	// Name defines the name of the referenced object
	Name string `json:"name"`
	// Namespace defines the namespace of the referenced object
	Namespace string `json:"namespace,omitempty"`
}

func (o *ObjectReference) IsNamespaceSet() bool {
	return o.Namespace != ""
}

func (o *ObjectReference) GetNamespace(defaultNamespace string) string {
	if o.IsNamespaceSet() {
		return o.Namespace
	}

	return defaultNamespace
}
