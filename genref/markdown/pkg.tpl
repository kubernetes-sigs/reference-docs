{{ define "packages" -}}

{{- range $idx, $val := .packages -}}
{{/* Only display package that has a group name */}}
{{- if and (ne .GroupName "") (eq $idx 0) -}}
---
title: {{ .Title }}
content_type: tool-reference
package: {{ .DisplayName }}
auto_generated: true
---
{{ .GetComment }}
{{- end -}}
{{- end }}

## Resource Types 

{{ range .packages }}
  {{ if ne .GroupName "" -}}
    {{- range .VisibleTypes -}}
      {{ if .IsExported }}
- [{{ .DisplayName }}]({{ .Link }})
      {{- end -}}
    {{- end -}}
  {{ end -}}
{{- end -}}

{{ range .packages }}
  {{ if ne .GroupName "" -}}
     
    {{/* For package with a group name, list all type definitions in it. */}}
    {{ range .VisibleTypes }}
{{ template "type" . }}
    {{ end }}
  {{ else }}
    {{/* For package w/o group name, list only types referenced. */}}
    {{- range .VisibleTypes -}}
      {{- if .Referenced -}}
{{ template "type" . }}
      {{- end -}}
    {{- end }}
  {{- end }}
{{- end }}
{{ end }}
