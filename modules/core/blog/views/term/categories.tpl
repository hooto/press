{{if .categories}}
<div class="hpress-sidebar-section">
  <div class="header">
    <h3>{{.categories.Model.Title}}</h3>
  </div>
  <div class="list-group">
    <a class="list-group-item" href="{{$.baseuri}}/list">All</a>
    {{range $v := .categories.Items}}
    {{$id := printf "%d" $v.ID}}
    <a class="list-group-item{{if eq $.term_categories $id}} active{{end}}" 
        href="{{$.baseuri}}/list?term_{{$.categories.Model.Meta.Name}}={{$id}}">{{$v.Title}}</a>
    {{end}}
  </div>
</div>
{{end}}
