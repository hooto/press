{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="l5s-ctn-header">
    <h2>Article list</h2>
  </div>

  <div class="row">
    
    <div class="col-md-3">
        
        {{pagelet . .specid "search.tpl"}}

        {{pagelet . .specid "term/categories.tpl"}}

    </div>

    <div class="col-md-9">
    
    <ul class="l5s-list">
      {{range $v := .list.Items}}
      <li class="l5s-list-item">
        <h4 class="l5s-list-heading"><a href="{{$.baseuri}}/view/{{$v.ID}}">{{$v.Title}}</a></h4>
        <span class="l5s-list-info">{{TimeFormat $v.Created "atom" "Y-m-d H:i"}}</span>
        <p class="l5s-list-text">{{FieldSubHtml $v.Fields "content" 300}}</p>
      </li>
      {{end}}
    </ul>

    {{if .list_pager}}
    <ul class="pagination pagination-sm">
      {{if .list_pager.FirstPageNumber}}
      <li>
        <a href="{{$.baseuri}}/list?page={{.list_pager.FirstPageNumber}}">First</a>
      </li>
      {{end}}

      {{range $index, $page := .list_pager.RangePages}}
      <li {{if eq $page $.list_pager.CurrentPageNumber}}class="active"{{end}}>
        <a href="{{$.baseuri}}/list?page={{$page}}">{{$page}}</a>
      </li>
      {{end}}
      
      {{if .list_pager.LastPageNumber}}
      <li>
        <a href="{{$.baseuri}}/list?page={{.list_pager.LastPageNumber}}">Last</a>
      </li>
      {{end}}
    </ul>
    {{end}}

    </div>    
    
  </div>  

</div>

{{pagelet . "general" "footer.tpl"}}
</body>
</html>
