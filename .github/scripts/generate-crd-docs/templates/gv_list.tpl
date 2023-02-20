{{- define "gvList" -}}
{{- $groupVersions := . -}}
{{- $groupVersion := index $groupVersions 0 -}}

---
title: {{ $groupVersion.Version }}
description: Reference information about the KLT CRDs
weight: 100
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
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
