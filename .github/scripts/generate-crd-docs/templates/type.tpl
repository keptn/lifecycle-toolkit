{{- define "type" -}}
{{- $type := . -}}
{{- if markdownShouldRenderType $type -}}

#### {{ $type.Name }}

{{ if $type.IsAlias }}_Underlying type:_ _{{ markdownRenderTypeLink $type.UnderlyingType  }}_{{ end }}

{{ $type.Doc }}

{{ if $type.Validation -}}
_Validation:_
{{- range $type.Validation }}
- {{ . }}
{{- end }}
{{- end }}

{{ if $type.References -}}
_Appears in:_
{{- range $type.SortedReferences }}
- {{ markdownRenderTypeLink . }}
{{- end }}
{{- end }}

{{ if $type.Members -}}
| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
{{ if $type.GVK -}}
| `apiVersion` _string_ | `{{ $type.GVK.Group }}/{{ $type.GVK.Version }}` | | | |
| `kind` _string_ | `{{ $type.GVK.Kind }}` | | | |
{{ end -}}

{{ range $type.Members -}}
| `{{ .Name }}` {{ if .Type.IsAlias }}_{{  markdownRenderTypeLink .Type.UnderlyingType  }}_{{else}}_{{ markdownRenderType .Type }}_{{ end }} | {{ template "type_members" . }} |
{{- if index .Markers "kubebuilder:default" -}}
{{- with index (index .Markers "kubebuilder:default") 0 -}}
 {{ .Value -}}
{{ end -}}
{{ end -}}
| {{ if index .Markers "optional" }}✓{{ else }}x{{ end }} | {{ range .Validation -}} {{ . }} <br />{{ end }} |
{{ end }}

{{- end -}}
{{- end -}}

{{- end -}}
