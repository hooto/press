<div id="htp-topbar">
  <div class="container htp-topbar-collapse">
    <ul class="htp-nav">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="htp-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}"></li>
      {{end}}
      <li class="htp-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>

    <ul class="htp-nav htp-topbar-nav" id="htp-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <li><a class="" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a></li>
      {{end}}
    </ul>

    <ul class="htp-nav htp-nav-right" id="htp-topbar-userbar"></ul>
  </div>
</div>

<script id="htp-topbar-user-unsigned-tpl" type="text/html">
<li class="iam-name"><a href="{{HttpSrvBasePath "auth/login"}}">Login</a></li>
</script>

<script id="htp-topbar-user-signed-tpl" type="text/html">

<li class="iam-name">{[=it.name]}</li>
<li class="iam-photo" id="htp-topbar-user-signed"><img src="{[=it.photo_url]}"/></li>

<div id="htp-topbar-user-signed-modal" style="display:none;">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.name]}</div>
  {[? it.instance_owner]}<a class="btn btn-primary iam-btn" href="{[=htp.HttpSrvBasePath('mgr')]}">Content Manage</a>{[?]}
  <a class="btn btn-default iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="btn btn-default iam-btn" href="{[=htp.HttpSrvBasePath('auth/sign-out')]}">Sign out</a>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    htp.NavActive("htp-topbar-nav-main", "{{.baseuri}}");
    htp.AuthSessionRefresh();
});
</script>
