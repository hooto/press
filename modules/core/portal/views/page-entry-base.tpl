<!DOCTYPE html>
<html lang="en">
{{pagelet . "core/general" "bs4/html-header.tpl"}}
<link rel="stylesheet" href="{{HttpSrvBasePath "hpress/+/static/pt/css/main.css"}}?v={{.sys_version_sign}}" type="text/css">
<link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/open-iconic/font/css/open-iconic-bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
<body>
{{pagelet . "core/general" "bs4/nav-header.tpl" "topnav"}}

<div>{{FieldHtml .page_entry.Fields "content"}}</div>

{{pagelet . "core/general" "footer.tpl"}}

{{pagelet . "core/general" "html-footer.tpl"}}
</body>
</html>

