{{/* vim: set filetype=mustache: */}}
{{/*
Return the proper image name
{{ include "common.images.image" ( dict "imageRoot" .Values.path.to.the.image "global" .Values.global ) }}
*/}}
{{- define "common.images.image" -}}
{{- $registryName :=  .global.imageRegistry -}}
{{- $repositoryName :=  .imageRoot -}}
{{- $separator := ":" -}}
{{- $termination := .global.imageTag | toString -}}
{{- printf "%s/%s%s%s" $registryName $repositoryName $separator $termination -}}
{{- end -}}
