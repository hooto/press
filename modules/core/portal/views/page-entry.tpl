<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "html-header.tpl"}}
<body>
{{pagelet . "core/general" "nav-header.tpl" "topnav"}}

<div class="container">
  
  <div class="hpress-ctn-header">
    <h2>{{.page_entry.Title}}</h2>
  </div>

  <div class="row">
    <div class="col-md-12">    
      <div class="hpress-nodev">
        <div class="content">{{FieldHtml .page_entry.Fields "content"}}</div>
      </div>      
    </div>
  </div>  

</div>

{{pagelet . "core/general" "footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>