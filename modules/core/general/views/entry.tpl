{{pagelet . "core/general" "html-header.tpl"}}

{{pagelet . "core/general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="htp-ctn-header">
    <h2>{{.page.Title}}</h2>
  </div>

  <div class="row">
    <div class="col-md-12">    
      <div class="htp-nodev">
        <div class="content">{{FieldHtml .page.Fields "content"}}</div>
      </div>      
    </div>
  </div>  

</div>

{{pagelet . "core/general" "footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
