{{/*
Return the proper Image Registry Secret Names
*/}}
{{- define "certManager.imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.image) "context" $) }}
{{- end -}}
