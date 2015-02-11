categories
{{if .categories}}
<ul>
    {{range $v := .categories}}
    <li>{{$v.Title}}</li>
    {{end}}
</ul>
{{end}}
