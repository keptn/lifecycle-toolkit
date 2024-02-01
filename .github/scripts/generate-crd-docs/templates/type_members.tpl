{{- define "type_members" -}}
{{- $field := . -}}
{{- if and (eq $field.Name "metadata") (eq $field.Type.String "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta") -}}
Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects).
{{- else -}}
{{ markdownRenderFieldDoc $field.Doc }}
{{- end -}}
{{- end -}}
