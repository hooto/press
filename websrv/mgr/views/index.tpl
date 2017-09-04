<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{SysConfig "frontend_html_head_subtitle"}}</title>
  <script src="{{HttpSrvBasePath "hpress/~/lessui/js/sea.js"}}?v={{.sys_version}}"></script>
  <script src="{{HttpSrvBasePath "hpress/~/hpressm/js/main.js"}}?v={{.sys_version}}"></script>
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "hpress/~/hpress/img/ap.ico"}}?v={{.sys_version}}">
  <script type="text/javascript">
    window._basepath = {{HttpSrvBasePath "/hpress"}};
  </script>
</head>
<body id="body-content">
<div id="hpressm-topbar" style="display:none">
  <div class="hpressm-topbar-collapse">
    <ul class="hpressm-nav" id="hpressm-topbar-siteinfo">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="hpressm-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}"></li>
      {{end}}
      <li class="hpressm-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>
    <ul class="hpressm-nav hpressm-topbar-nav" id="hpressm-topbar-nav-node-specls">
    </ul>
    <ul class="hpressm-nav hpressm-nav-right" id="hpressm-topbar-userbar">
      <li><a href="{{HttpSrvBasePath "hpress/auth/sign-out"}}">Sign Out</a></li>
    </ul>
    <ul class="hpressm-nav hpressm-nav-right">
      <li><a class="l4i-nav-item" href="#s2/index">Storage</a></li>
      <li><a class="l4i-nav-item" href="#spec/index">Modules</a></li>
      <li><a class="l4i-nav-item" href="#sys/index">System</a></li>
    </ul>
  </div>
</div>
<div id="com-content" class="">loading</div>
<script id="hpressm-topbar-nav-node-specls-tpl" type="text/html">
{[~it.items :v]}
  <li><a class="l4i-nav-item" href="#node/index/{[=v.meta.name]}">{[=v.title]}</a></li>
{[~]}
</script>
</body>
<script type="text/javascript">
window._sys_version = {{.sys_version}};
window.onload = hpressMgr.Boot();
</script>
</html>
