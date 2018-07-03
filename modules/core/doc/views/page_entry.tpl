<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "bs4/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/doc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body>
{{pagelet . "core/general" "bs4/nav-header.tpl" "topnav"}}

<div class="hpdoc_index_frame_blue hpdoc_node_content hpdoc_bgimg_hexagons">
<div class="container" style="padding: 20px 10px; text-align: center;">
  <div class="hpdoc_index_frame_title">{{FieldStringPrint .doc_entry "title" .LANG}}</div>
</div>
</div>

<div class="hpdoc_index_frame_dark_light">
<div class="container hpdoc_node_content">
  <div class="row" style="padding: 20px 0;">
    <div class="col-3">
      <div class="hpdoc_entry_summary markdown-body hp-content" style="">
        <ul><li><a href=".">Preface</a></li></ul>
        {{FieldHtmlPrint .doc_entry "content" .LANG}}
      </div>
    </div>
    <div class="col-9">
      <div class="hpdoc_entry_content hpdoc_page_content markdown-body hp-content">{{FieldHtmlPrint .page_entry "content" .LANG}}</div>
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
