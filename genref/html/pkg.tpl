{{ define "packages" }}
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.8.2/css/font-awesome.min.css">
      <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css">
      <style type="text/css">
        td p {
          margin-bottom: 0
        }
        code {
          color: #802060;
          display: inline-block;
        }
      </style>
    </head>
    <body>
      <div class="container">
        <!--
        <H2>Packages:</H2>
        <ul>
          {{ range .packages }}
            {{/* Skip package which doesn't have a group name defined */}}
            {{ if ne .GroupName ""}}
              <li>
                <a href="#{{- .Anchor -}}">{{ .DisplayName }}</a>
              </li>
            {{ end }}
          {{ end }}
        </ul>
        -->
        {{ range .packages }}
          {{/* Only display package that has a group name */}}
          {{ if ne .GroupName "" }}
            <H2 id="{{- .Anchor -}}">Package: <span style="font-family: monospace">{{- .DisplayName -}}</span></H2>
            <p>{{ .GetComment }}</p>

            {{/* TODO: Make the following line conditional */}}
            <H3>Resource Types:</H3>
            <ul>
              {{- range .VisibleTypes -}}
                {{ if .IsExported -}}
                  <li>
                    <a href="{{ .Link }}">{{ .DisplayName }}</a>
                  </li>
                {{- end }}
              {{- end -}}
            </ul>

            {{/* For package with a group name, list all type definitions in it. */}}
            {{ range .VisibleTypes }}
              {{ template "type" .  }}
            {{ end }}
          {{ else }}
            {{/* For package without a group name, list only type definitions that are referenced. */}}
            {{ range .VisibleTypes }}
              {{ if .Referenced }}
                {{ template "type" .  }}
              {{ end }}
            {{ end }}
          {{ end }}

          <HR />
        {{ end }}
      </div>

      <div class="container">
        <p><em>Generated with <code>gendoc</code>{{ with .gitCommit }} on git commit <code>{{ . }}</code>{{end}}</em></p>
      </div>
    </body>
  </html>
{{ end }}

