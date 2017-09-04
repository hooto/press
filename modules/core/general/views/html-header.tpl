<head>
  <meta charset="utf-8">
  <title>{{if .__html_head_title__}}{{.__html_head_title__}} | {{end}}{{SysConfig "frontend_html_head_subtitle"}}</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/bs/3.3/css/bootstrap.css"}}" type="text/css">
  <link rel="stylesheet" href="{{HttpSrvBasePath "hpress/~/hpress/css/main.css"}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "hpress/~/hpress/img/ap.ico"}}">
  <meta name="keywords" content="{{SysConfig "frontend_html_head_meta_keywords"}}">
  <meta name="description" content="{{SysConfig "frontend_html_head_meta_description"}}">
  <script src="{{HttpSrvBasePath "hpress/~/lessui/js/sea.js"}}"></script>
  <script src="{{HttpSrvBasePath "hpress/~/hpress/js/main.js"}}"></script>
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath "hpress/"}};
    window.onload_hooks = [];
  </script>
</head>
