package common

import ctrl "sigs.k8s.io/controller-runtime"

// GetRequestInfo extracts name and namespace from a controller request.
func GetRequestInfo(req ctrl.Request) map[string]string {
	return map[string]string{
		"name":      req.Name,
		"namespace": req.Namespace,
	}
}
