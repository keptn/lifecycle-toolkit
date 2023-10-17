{{/*
Expand the name of the chart.
*/}}
{{- define "common.name" -}}
{{- default .context.Chart.Name .context.Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "common.fullname" -}}
{{- if .context.Values.fullnameOverride }}
{{- .context.Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .context.Chart.Name .context.Values.nameOverride }}
{{- if contains $name .context.Release.Name }}
{{- .context.Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .context.Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "common.chart" -}}
{{- printf "%s-%s" (kebabcase .context.Chart.Name) .context.Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}


{{- define "common.labels.standard" -}}
{{- $default := fromYaml (include "common.labels" (dict "context" .context) ) -}}
{{ template "common.tplvalues.merge" (dict "values" (list .context.Values.global.commonLabels $default) "context" .context ) }}
{{- end -}}


{{/*
Common labels
*/}}
{{- define "common.labels" -}}
helm.sh/chart: {{ include "common.chart" (dict "context" .context) }}
{{ include "common.selectorLabels" (dict "context" .context) }}
{{- if .context.Chart.AppVersion }}
app.kubernetes.io/version: {{ .context.Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .context.Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "common.selectorLabels" -}}
app.kubernetes.io/name: {{kebabcase (include  "common.name" (dict "context" .context) ) }}
app.kubernetes.io/instance: {{ .context.Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "common.serviceAccountName" -}}
{{- if .context.Values.serviceAccount.create }}
{{- default (include "common.fullname" (dict "context" .context )) .context.Values.serviceAccount.name }}
{{- else }}
{{- default "default" .context.Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "common.annotations" -}}
  {{- if or .context.Values.annotations .context.Values.global.commonAnnotations }}
  {{- $annotations := include "common.tplvalues.merge" ( dict "values" ( list .context.Values.annotations .context.Values.global.commonAnnotations ) "context" .context ) }}
  annotations: {{- include "common.tplvalues.render" ( dict "value" $annotations "context" .context) | nindent 4 }}
  {{- end }}
{{- end }}