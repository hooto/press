<div id="hpress-topbar">
  <div class="container hpress-topbar-collapse">
    <ul class="hpress-nav">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="hpress-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}"></li>
      {{end}}
      <li class="hpress-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>

    <ul class="hpress-nav hpress-topbar-nav" id="hpress-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <li><a class="" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a></li>
      {{end}}
    </ul>

    <ul class="hpress-nav hpress-nav-right" id="hpress-topbar-userbar"></ul>
  </div>
</div>

<script id="hpress-topbar-user-unsigned-tpl" type="text/html">
<li class="iam-name"><a href="{{HttpSrvBasePath "hpress/auth/login"}}">Login</a></li>
</script>

<script id="hpress-topbar-user-signed-tpl" type="text/html">

<li class="iam-name">{[=it.display_name]}</li>
<li class="iam-photo" id="hpress-topbar-user-signed"><img src="{[=it.photo_url]}"/></li>

<div id="hpress-topbar-user-signed-modal" style="display:none;">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.display_name]}</div>
  {[? it.instance_owner]}<a class="btn btn-primary iam-btn" href="{[=hpress.HttpSrvBasePath('mgr')]}" target="_blank">Content Manage</a>{[?]}
  <a class="btn btn-default iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="btn btn-default iam-btn" href="{[=hpress.HttpSrvBasePath('auth/sign-out')]}">Sign out</a>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    hpress.NavActive("hpress-topbar-nav-main", "{{.http_request_path}}");
    hpress.AuthSessionRefresh();
});
</script>
