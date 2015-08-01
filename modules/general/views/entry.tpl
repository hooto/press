{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="l5s-ctn-header">
    <h2>{{.page.Title}}</h2>
  </div>

  <div class="row">
    <div class="col-md-12">    
      <div class="l5s-nodev">
        <div class="content">{{FieldHtml .page.Fields "content"}}</div>
      </div>      
    </div>
  </div>  

</div>

{{pagelet . "general" "footer.tpl"}}

</body>
</html>
