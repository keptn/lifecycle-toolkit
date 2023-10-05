package handlers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/handler_mock.go . K8sHandler:MockHandler
type K8sHandler interface {
	Handle(ctx context.Context, pod *corev1.Pod, namespace string) error
}
