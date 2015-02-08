{{pagelet . "general" "html-header.tpl"}}
<body>
<div class="container">
  
  {{pagelet . "general" "nav-header.tpl" "topnav"}}

  <h2>Article list</h2>

  <div class="row">
    <div class="col-md-9">
    <ul class="l5s-list">
    {{range $v := .list}}
    <li class="l5s-list-item">
        <h4 class="l5s-list-heading"><a href="#view/{{Field $v "id"}}">{{Field $v "title"}}</a></h4>
        <span class="l5s-list-info">{{Field $v "created"}}</span>
        <p class="l5s-list-text">{{Field $v "field_content"}}</p>
    </li>
    {{end}}
    </ul>
    </div>
    
    <div class="col-md-3">
        {{range $v := .categories}}

        {{end}}
    </div>
  </div>

  

  {{pagelet . "general" "footer.tpl"}}
</div>
</body>
</html>
