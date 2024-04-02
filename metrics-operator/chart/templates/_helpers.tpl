{{/*
Return the proper Image Registry Secret Names
*/}}
{{- define "metricsOperator.imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.image) "context" $) }}
{{- end -}}
