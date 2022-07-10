<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v3/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/fa/v5/css/fas.css"}}?v={{.sys_version_sign}}" type="text/css">
<script src="{{HttpSrvBasePath "hp/-/static/gdoc/js/gdoc.js"}}?v={{.sys_version_sign}}"></script>
<body id="hp-body">
{{pagelet . "core/general" "v3/nav-header.tpl" "topnav" "topbar_class=navbar-light"}}

<div class="hp-container-full hp-gdoc-index-frame-dark-light">
<nav class="container hp-block-p-std">
  <ol class="breadcrumb " style="margin:0">
        <li class="breadcrumb-item">
          <span class="icon"><i class="fas fa-file-alt"></i></span>
          <a href="{{.baseuri}}/" style="margin-left: 10px">
            <span>{{T .LANG "Documents"}}</span>
          </a>
        </li>
        <li class="breadcrumb-item active">
          <a href="{{.baseuri}}/view/{{.doc_entry.ExtPermalinkName}}/">
            {{FieldStringPrint .doc_entry "title" .LANG}}
          </a>
        </li>
  </ol>
</nav>
</div>

<div class="hp-container-full hp-gdoc-index-frame-dark-light">
<div class="container hp-gdoc-node-content">
  <div class="row hp-is-mobile">
    <div class="col-auto">
      <button class="btn btn-dark" onclick="hp.NavbarMenuToggle('hp-gdoc-entry-summary')">
	    <span class="icon">
	      <i class="fas fa-list"></i>
        </span>
	    <span>
		  Menu
        </span>
      </button>
    </div>
  </div>
  <div class="row">
    <div class="col-md-2 hp-block-px-std hp-block-pb-std">
      <div id="hp-gdoc-entry-summary" class=" hp-gdoc-entry-summary hp-content hp-is-desktop hp-scroll">
      {{FieldHtmlPrint .doc_entry "content" .LANG}}
      </div>
    </div>
    <div class="col-md-8 hp-block-px-std hp-block-pb-std">
      <div id="hp-gdoc-page-entry-content" class="hp-gdoc-entry-content hp-gdoc-page-content content hp-content">
      {{FieldHtmlPrint .doc_entry "preface" .LANG}}
      </div>
    </div>
    <div id="hp-gdoc-page-entry-toc" class="col-md-2 hp-is-desktop hp-gdoc-entry-toc hp-scroll hp-block-px-std hp-block-pb-std">
    </div>
  </div>
</div>
</div>

{{pagelet . "core/general" "v3/footer.tpl"}}


<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.CodeRender({"theme": "monokai"});
	gdoc.PageEntryRender({
        "doc_base_path": "{{.baseuri}}/view/{{.doc_entry.ExtPermalinkName}}/",
	});
});
</script>

{{pagelet . "core/general" "html-footer.tpl"}}

</body>
</html>
