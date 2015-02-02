{{template "html-header.tpl"}}
<body>
{{template "nav-header.tpl"}}
<div class="container">
{{with $ls := Query . "lower:HEllo"}}

  {{range $v := $ls}}
    {{Field $v "id" "string"}}
    {{$v.Field "created"}}
  {{end}}

{{end}}
</div>
</body>
</html>
