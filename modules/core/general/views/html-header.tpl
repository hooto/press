<head>
  <meta charset="utf-8">
  <title>{{if .__html_head_title__}}{{.__html_head_title__}} | {{end}}{{SysConfig "frontend_html_head_subtitle"}}</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/bs/3.3/css/bootstrap.css"}}?v={{.sys_version}}" type="text/css">
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/hpress/css/main.css"}}?v={{.sys_version}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "hpress/~/hpress/img/ap.ico"}}?v={{.sys_version}}">
  <meta name="keywords" content="{{SysConfig "frontend_html_head_meta_keywords"}}">
  <meta name="description" content="{{SysConfig "frontend_html_head_meta_description"}}">
  <script src="{{HttpSrvBasePath "hpress/~/lessui/js/sea.js"}}?v={{.sys_version}}"></script>
  <script src="{{HttpSrvBasePath "hpress/~/hpress/js/main.js"}}?v={{.sys_version}}"></script>
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath "hpress/"}};
    window._sys_version = {{.sys_version}};
    window.onload_hooks = [];
  </script>
</head>
