<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{if .__html_head_title__}}{{.__html_head_title__}}{{else}}CMS{{end}}</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "~/bootstrap/3.3/css/bootstrap.min.css"}}" type="text/css">
  <link rel="stylesheet" href="{{HttpSrvBasePath "~/l5s/css/main.css"}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "~/l5s/img/ap.ico"}}">
  <script src="{{HttpSrvBasePath "/~/lessui/js/sea.js"}}"></script>
  <script src="{{HttpSrvBasePath "/~/l5s/js/main.js"}}"></script>
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath ""}};
    window.onload_hooks = [];
    window.onload = l5s.Boot();
  </script>
</head>

