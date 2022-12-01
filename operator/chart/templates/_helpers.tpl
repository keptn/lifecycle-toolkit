{{/*
Define the chart.namespace variable. If no namespace is defined via the command line then the default value 
provided in .Values.defaultNamespace is used
*/}}
{{- define "chart.namespace" -}}
{{- if eq .Release.Namespace "default" -}}
{{- .Values.defaultNamespace -}}
{{- else -}}
{{- .Release.Namespace -}}
{{- end -}}
{{- end -}}
