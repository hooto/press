{{pagelet . "general" "html-header.tpl"}}
<body>
{{pagelet . "general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="l5s-ctn-header">
    <h2>Article Entry</h2>
  </div>

  <div class="row">
    <div class="col-md-12">
    
      <div class="node-entry-view">
        
        <div class="header">
          <h2>{{.entry.Title}}</h2>
          <div class="hinfo">{{TimeFormat .entry.Created "atom" "Y-m-d H:i"}}</div>
        </div>

        <div>{{FieldHtml .entry.Fields "content"}}</div>
      </div>
    
    </div>
    
  </div>  

</div>

{{pagelet . "general" "footer.tpl"}}

</body>
</html>
