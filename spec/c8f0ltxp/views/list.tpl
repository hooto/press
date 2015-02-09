{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}

<div class="container">
  <h2>Article list</h2>
  
  <div class="row">
    <div class="col-md-9">
    <ul class="l5s-list">
    {{range $v := .list.Items}}
    <li class="l5s-list-item">
        <h4 class="l5s-list-heading"><a href="{{$.baseuri}}/view/{{$v.ID}}">{{$v.Title}}</a></h4>
        <span class="l5s-list-info">{{TimeFormat $v.Created "atom" "Y-m-d H:i"}}</span>
        <p class="l5s-list-text">{{FieldSubHtml $v.Fields "content" 400}}</p>
    </li>
    {{end}}
    </ul>
    </div>
    
    <div class="col-md-3">
      a
        {{range $v := .categories}}

        {{end}}
    </div>
  </div>  

  {{pagelet . "general" "footer.tpl"}}
</div>
</body>
</html>
