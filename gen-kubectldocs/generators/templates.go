/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

var CategoryTemplate = `
# <strong>{{.Name}}</strong>
`

var CommandTemplate = `{{"---\ntitle: "}}{{.MainCommand.Name}}
{{"content_template: templates/tool-reference\n---"}}

### Overview
{{.MainCommand.Description}}

### Usage

` + "`" + `$ {{.MainCommand.Usage}}` + "`" + `

{{if .MainCommand.Example}}
### Example

{{.MainCommand.Example}}
{{end}}

{{if .MainCommand.Options}}
### Flags

<div class="table-responsive"><table class="table table-bordered">
<thead class="thead-light">
<tr>
            <th>Name</th>
            <th>Shorthand</th>
            <th>Default</th>
            <th>Usage</th>
        </tr>
    </thead>
    <tbody>
    {{range $option := .MainCommand.Options}}
    <tr>
    <td>{{$option.Name}}</td><td>{{$option.Shorthand}}</td><td>{{$option.DefaultValue}}</td><td>{{$option.Usage}}</td>
    </tr>{{end}}
</tbody>
</table></div>
{{end}}
{{range $sub := .SubCommands}}

<hr>

## {{$sub.Path}}


### Overview
{{$sub.Description}}

### Usage

` + "`" + `$ {{$sub.Usage}}` + "`" + `

{{if $sub.Example}}
### Example
{{$sub.Example}}
{{end}}

{{if $sub.Options}}
### Flags

<div class="table-responsive"><table class="table table-bordered">
<thead class="thead-light">
<tr>
            <th>Name</th>
            <th>Shorthand</th>
            <th>Default</th>
            <th>Usage</th>
        </tr>
    </thead>
    <tbody>
    {{range $option := $sub.Options}}
    <tr>
    <td>{{$option.Name}}</td><td>{{$option.Shorthand}}</td><td>{{$option.DefaultValue}}</td><td>{{$option.Usage}}</td>
    </tr>{{end}}
</tbody>
</table></div>
{{end}}
{{end}}

`
