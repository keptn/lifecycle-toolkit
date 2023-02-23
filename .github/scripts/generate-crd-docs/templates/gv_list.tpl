{{- define "gvList" -}}
{{- $groupVersions := . -}}
{{- $groupVersion := index $groupVersions 0 -}}

---
title: {{ $groupVersion.Version }}
description: Reference information for {{ $groupVersion.GroupVersionString }}
---
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
