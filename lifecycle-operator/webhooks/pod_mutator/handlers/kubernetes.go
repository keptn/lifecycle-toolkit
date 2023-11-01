package handlers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/handler_mock.go . K8sHandler:MockHandler
type K8sHandler interface {
	Handle(ctx context.Context, pod *corev1.Pod, namespace string) error
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/decoder_mock.go . Decoder:MockDecoder
type Decoder interface {
	Decode(req admission.Request, into runtime.Object) error
	DecodeRaw(rawObj runtime.RawExtension, into runtime.Object) error
}
