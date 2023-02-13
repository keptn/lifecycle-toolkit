{{/*
Return the proper image name
*/}}
{{- define "lifecycletoolkit.operator.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.imageLifecycle "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name
*/}}
{{- define "lifecycletoolkit.scheduler.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.imageScheduler "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name
*/}}
{{- define "lifecycletoolkit.certmanager.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.imageCertmanager "global" .Values.global) }}
{{- end -}}