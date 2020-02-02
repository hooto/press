<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "bs4/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body>
{{pagelet . "core/general" "bs4/nav-header.tpl" "topnav"}}

<div class="hp-gdoc-index-frame-blue hp-gdoc-node-content hp-gdoc-bgimg-hexagons">
<div class="container" style="padding: 10px;">
  <div class="row">
  <div class="col">
    <ol class="hp-gdoc-nav-ol">
      <li>
        <span class="oi oi-document" title="icon name" aria-hidden="true"></span>
        <a href="{{.baseuri}}/">
          {{T .LANG "Documents"}}
        </a>
      </li>
      <li class="active">
        <a href="{{.baseuri}}/view/{{.doc_entry.ID}}/">
          {{FieldStringPrint .doc_entry "title" .LANG}}
        </a>
      </li>
    </ol>
  </div>
  <div class="col font-dark" style="text-align: right;">
    <a href="{{FieldStringPrint .doc_entry "repo_url" .LANG}}">
      Git Version {{FieldStringPrint .doc_entry "repo_version" .LANG}}
    </a>
  </div>
  </div>
</div>
</div>

<div class="hp-gdoc-index-frame-dark-light">
<div class="container hp-gdoc-node-content">
  <div class="row" style="padding: 20px 0;">
    <div class="col" style="min-width:200px;max-width:300px">
      <div class="hp-gdoc-entry-summary hp-content" style="">
        {{FieldHtmlPrint .doc_entry "content" .LANG}}
      </div>
    </div>
    <div class="col" style="min-width:400px;">
      <div class="hp-gdoc-entry-content hp-gdoc-page-content hp-content">{{FieldHtmlPrint .page_entry "content" .LANG}}</div>
    </div>
  </div>
</div>
</div>
    
{{pagelet . "core/general" "footer.tpl"}}

<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.CodeRender();
});
</script>

{{pagelet . "core/general" "html-footer.tpl"}}

</body>
</html>
