<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v2/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/gdoc/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/fa/v5/css/fas.css"}}?v={{.sys_version_sign}}" type="text/css">
<body id="hp-body">
{{pagelet . "core/general" "v2/nav-header.tpl" "topnav"}}

<div class="hp-container-full hp-gdoc-index-frame-dark hp-gdoc-node-content hp-gdoc-bgimg-hexagons">
<div class="container" style="padding: 10px;">
  <div class="columns">
    <div class="column">
      <ol class="hp-gdoc-nav-ol">
        <li>
          <a href="{{.baseuri}}/">
            <span class="icon"><i class="fas fa-file-alt"></i></span>
            <span>{{T .LANG "Documents"}}</span>
          </a>
        </li>
        <li class="active">
          <a href="{{.baseuri}}/view/{{.doc_entry.ExtPermalinkName}}/">
            {{FieldStringPrint .doc_entry "title" .LANG}}
          </a>
        </li>
      </ol>
    </div>
    <div class="column font-dark hp-is-desktop" style="text-align: right;">
      <a href="{{FieldStringPrint .doc_entry "repo_url" .LANG}}">
        Git Version {{FieldStringPrint .doc_entry "repo_version" .LANG}}
      </a>
    </div>
  </div>
</div>
</div>

<div class="hp-container-full hp-gdoc-index-frame-dark-light">
<div class="container hp-gdoc-node-content">
  <div class="columns hp-is-mobile">
    <div class="column">
      <button class="button" onclick="hp.NavbarMenuToggle('hp-gdoc-entry-summary')">
	    <span class="icon">
	      <i class="fas fa-list"></i>
        </span>
	    <span>
		  Menu
        </span>
      </button>
    </div>
  </div>
  <div class="columns">
    <div id="hp-gdoc-entry-summary" class="column is-3 hp-is-desktop">
      <div class="hp-gdoc-entry-summary hp-content" style="">
        {{FieldHtmlPrint .doc_entry "content" .LANG}}
      </div>
    </div>
    <div class="column is-9">
      <div class="hp-gdoc-entry-content hp-gdoc-page-content content hp-content">{{FieldHtmlPrint .doc_entry "preface" .LANG}}</div>
    </div>
  </div>
</div>
</div>

{{pagelet . "core/general" "v2/footer.tpl"}}


<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.CodeRender();
});
</script>

{{pagelet . "core/general" "html-footer.tpl"}}

</body>
</html>
