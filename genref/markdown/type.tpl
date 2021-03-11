{{ define "type" }}

## `{{ .Name.Name }}`     {#{{ .Anchor }}}
    
{{ if eq .Kind "Alias" -}}
(Alias of `{{ .Underlying }}`)
{{- end }}

{{ with .References }}
**Appears in:**
{{ range . }}
- [{{ .DisplayName }}]({{ .Link }})
{{ end }}
{{- end }}

{{ if .GetComment -}}
{{ .GetComment }}
{{- end }}

{{ if .GetMembers -}}
<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    {{/* . is a apiType */}}
    {{- if .IsExported -}}
{{/* Add apiVersion and kind rows if deemed necessary */}}
<tr><td><samp>apiVersion</samp><br/>string</td><td><samp>{{- .APIGroup -}}</samp></td></tr>
<tr><td><samp>kind</samp><br/>string</td><td><samp>{{- .Name.Name -}}</samp></td></tr>
    {{ end -}}

{{/* The actual list of members is in the following template */}}
{{- template "members" . -}}
</tbody>
</table>
{{- end -}}
{{- end -}}
