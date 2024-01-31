{{- define "type_members" -}}
{{- $field := . -}}
{{- if and (eq $field.Name "metadata") (eq $field.Type.String "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta") -}}
Refer to Kubernetes API documentation for fields of `metadata`.
{{- else -}}
{{ markdownRenderFieldDoc $field.Doc }}
{{- end -}}
{{- end -}}
