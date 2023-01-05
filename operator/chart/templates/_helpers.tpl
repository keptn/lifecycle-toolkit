{{/*
Define the chart.namespace variable. Use the namespace in .Values.defaultNamespace if provided
*/}}
{{- define "chart.namespace" -}}
{{- if .Values.defaultNamespace -}}
{{- .Values.defaultNamespace -}}
{{- else -}}
default
{{- end -}}
{{- end -}}
