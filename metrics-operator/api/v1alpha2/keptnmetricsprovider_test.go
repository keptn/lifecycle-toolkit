/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestHasSecretDefined(t *testing.T) {
	tests := []struct {
		name     string
		provider KeptnMetricsProvider
		result   bool
	}{
		{
			name: "Correct Definition",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "mysecret",
						},
						Key: "mykey",
					},
				},
			},
			result: true,
		},
		{
			name: "Missing key",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "mysecret",
						},
						Key: "",
					},
				},
			},
			result: false,
		},
		{
			name: "Missing name",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "",
						},
						Key: "mykey",
					},
				},
			},
			result: false,
		},
		{
			name: "Key made by spaces",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "mysecret",
						},
						Key: "    ",
					},
				},
			},
			result: false,
		},
		{
			name: "Name made by spaces",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "    ",
						},
						Key: "mykey",
					},
				},
			},
			result: false,
		},
		{
			name: "Empty secret struct",
			provider: KeptnMetricsProvider{
				Spec: KeptnMetricsProviderSpec{
					TargetServer: "",
					SecretKeyRef: corev1.SecretKeySelector{},
				},
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.result, tt.provider.HasSecretDefined())
		})

	}
}
