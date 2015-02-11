{{if .categories}}
<div class="l5s-sidebar-section">
  <div class="header">
    <h3>{{.categories.Model.Title}}</h3>
  </div>
  <div class="list-group">
    <!-- <a class="list-group-item active" href="{{$.baseuri}}/list">All</a> -->
    {{range $v := .categories.Items}}
    <a class="list-group-item" href="{{$.baseuri}}/list?term_{{$.categories.Model.Metadata.Name}}={{printf "%d" $v.ID}}">{{$v.Title}}</a>
    {{end}}
  </div>
</div>
{{end}}
