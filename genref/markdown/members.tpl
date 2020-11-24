{{ define "members" }}

  {{/* . is a apiType */}}
  {{- range .GetMembers -}}
    {{/* . is a apiMember */}}
    {{- if not .Hidden }}
<tr><td><samp>{{ .FieldName }}</samp>
      {{- if not .IsOptional }} <B>[Required]</B>{{- end -}}
<br/>
{{/* Link for type reference */}}
      {{- with .GetType -}}
        {{- if .Link -}}
<a href="{{ .Link }}"><samp>{{ .DisplayName }}</samp></a>
        {{- else -}}
<samp>{{ .DisplayName }}</samp>
        {{- end -}}
      {{- end }}
</td>
<td>
   {{- if .IsInline -}}
(Members of <samp>{{ .FieldName }}</samp> are embedded into this type.)
   {{- end }}
   {{ if .GetComment -}}
   {{ .GetComment }}
   {{- else -}}
   <span class="text-muted">No description provided.</span>
   {{ end }}
   {{- if and (eq (.GetType.Name.Name) "ObjectMeta") -}}
Refer to the Kubernetes API documentation for the fields of the <samp>metadata</samp> field.
   {{- end -}}
</td>
</tr>
    {{ end }}
  {{ end }}
{{ end }}
