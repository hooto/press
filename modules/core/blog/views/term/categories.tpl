{{if .categories}}
<div class="hp-sidebar-section">
  <div class="header">
    <h3>{{.categories.Model.Title}}</h3>
  </div>
  <div class="list-group term-taxonomy-group">
    <div class="list-group-item">
      <a class="term-taxonomy-item" href="{{$.baseuri}}/list">All</a>
    </div>
    {{range $v := .categories.Items}}
    {{$id := printf "%d" $v.ID}}
    {{$pid := printf "%d" $v.PID}}
    {{if eq $pid "0"}}
    <div class="list-group-item">
      <a class="term-taxonomy-item {{if eq $.term_categories $id}} active{{end}}" 
        href="{{$.baseuri}}/list?term_{{$.categories.Model.Meta.Name}}={{$id}}">{{$v.Title}}</a>
      {{range $v2 := $.categories.Items}}
      {{$id2 := printf "%d" $v2.ID}}
      {{$pid2 := printf "%d" $v2.PID}}
      {{if eq $pid2 $id}}
      <a class="term-taxonomy-subitem {{if eq $.term_categories $id2}} active{{end}}"
        href="{{$.baseuri}}/list?term_{{$.categories.Model.Meta.Name}}={{$v2.ID}}">{{$v2.Title}}</a>
      {{end}}
      {{end}}
    </div>
    {{end}}
    {{end}}
  </div>
</div>
{{end}}
