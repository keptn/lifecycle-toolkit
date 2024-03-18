{{/*
Return the proper Image Registry Secret Names for lifecycle operator
*/}}
{{- define "lifecycleOperator.imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.lifecycleOperator.image) "context" $) }}
{{- end -}}

{{/*
Return the proper Image Registry Secret Names for scheduler
*/}}
{{- define "scheduler.imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.scheduler.image) "context" $) }}
{{- end -}}
