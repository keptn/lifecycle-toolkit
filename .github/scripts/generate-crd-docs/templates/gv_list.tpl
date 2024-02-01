{{- define "gvList" -}}
{{- $groupVersions := . -}}
{{- $groupVersion := index $groupVersions 0 -}}

# {{ $groupVersion.Version }}

Reference information for {{ $groupVersion.GroupVersionString }}

<!-- markdownlint-disable -->

## Packages
{{- range $groupVersions }}
- {{ markdownRenderGVLink . }}
{{- end }}

{{ range $groupVersions }}
{{ template "gvDetails" . }}
{{ end }}

{{- end -}}
<!-- markdownlint-enable -->
