<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "v2/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/-/static/pt/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/fa/v5/css/fas.css"}}?v={{.sys_version_sign}}" type="text/css">
<body id="hp-body">
{{pagelet . "core/general" "v2/nav-header.tpl" "topnav"}}

{{FieldHtmlPrint .page_entry "content" .LANG}}

{{pagelet . "core/general" "v2/footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>

