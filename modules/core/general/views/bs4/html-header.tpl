<head>
  <meta charset="utf-8">
  <title>{{if .__html_head_title__}}{{.__html_head_title__}} | {{end}}{{SysConfig "frontend_html_head_subtitle"}}</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/bs/4.0/css/bootstrap.css"}}?v={{.sys_version_sign}}" type="text/css">
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/hpress/css/base.css"}}?v={{.sys_version_sign}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "hpress/~/hpress/img/ap.ico"}}?v={{.sys_version_sign}}">
  <meta name="keywords" content="{{SysConfig "frontend_html_head_meta_keywords"}}">
  <meta name="description" content="{{SysConfig "frontend_html_head_meta_description"}}">
  <script src="{{HttpSrvBasePath "hpress/~/lessui/js/sea.js"}}?v={{.sys_version_sign}}"></script>
  <script src="{{HttpSrvBasePath "hpress/~/hpress/js/main.js"}}?v={{.sys_version_sign}}"></script>
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath "hpress/"}};
    window._sys_version_sign = {{.sys_version_sign}};
    window.onload_hooks = [];
  </script>
</head>
