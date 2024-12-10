{{/*
Return the proper Image Registry Secret Names for lifecycle operator
*/}}
{{- define "imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.image) "context" $) }}
{{- end -}}
