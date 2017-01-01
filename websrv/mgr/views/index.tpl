<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{SysConfig "frontend_html_head_subtitle"}}</title>
  <script src="{{HttpSrvBasePath "mgr/~/lessui/js/sea.js"}}"></script>
  <script src="{{HttpSrvBasePath "mgr/-/js/main.js"}}"></script>
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "mgr/~/htp/img/ap.ico"}}">
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath ""}};
  </script>
</head>
<body id="body-content">
<div id="htpm-topbar" style="display:none">
  <div class="htpm-topbar-collapse">
    <ul class="htpm-nav" id="htpm-topbar-siteinfo">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="htpm-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}"></li>
      {{end}}
      <li class="htpm-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>
    <ul class="htpm-nav htpm-topbar-nav" id="htpm-topbar-nav-node-specls">
    </ul>
    <ul class="htpm-nav htpm-nav-right" id="htpm-topbar-userbar">
      <li><a href="{{HttpSrvBasePath "auth/sign-out"}}">Sign Out</a></li>
    </ul>
    <ul class="htpm-nav htpm-nav-right">
      <li><a class="l4i-nav-item" href="#s2/index">Storage</a></li>
      <li><a class="l4i-nav-item" href="#spec/index">Modules</a></li>
      <li><a class="l4i-nav-item" href="#sys/index">System</a></li>
    </ul>
  </div>
</div>
<div id="com-content" class="">loading</div>
<script id="htpm-topbar-nav-node-specls-tpl" type="text/html">
{[~it.items :v]}
  <li><a class="l4i-nav-item" href="#node/index/{[=v.meta.name]}">{[=v.title]}</a></li>
{[~]}
</script>
</body>
<script type="text/javascript">
window.onload = htpMgr.Boot();
</script>
</html>