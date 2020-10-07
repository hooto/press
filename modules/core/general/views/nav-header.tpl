<div id="hp-topbar">
  <div class="container hp-topbar-collapse">
    <ul class="hp-nav">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="hp-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}" height="30"></li>
      {{end}}
      <li class="hp-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>

    <ul class="hp-nav hp-topbar-nav" id="hp-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <li class="nav-item"><a class="nav-link" href="{{FieldString $v.Fields "url"}}">{{FieldStringPrint $v "title" $.LANG}}</a></li>
      {{end}}
    </ul>

    <ul class="hp-nav hp-nav-right" id="hp-topbar-userbar"></ul>
  </div>
</div>

<script id="hp-topbar-user-unsigned-tpl" type="text/html">
<li class="iam-name"><a href="{{HttpSrvBasePath "hp/auth/login"}}">Login</a></li>
</script>

<script id="hp-topbar-user-signed-tpl" type="text/html">

<li class="iam-name">{[=it.display_name]}</li>
<li class="iam-photo" id="hp-topbar-user-signed"><img src="{[=it.photo_url]}"/></li>

<div id="hp-topbar-user-signed-modal" style="display:none;">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.display_name]}</div>
  {[? it.instance_owner]}<a class="btn btn-primary iam-btn" href="{[=hp.HttpSrvBasePath('mgr')]}" target="_blank">Content Manage</a>{[?]}
  <a class="btn btn-primary iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="btn btn-primary iam-btn" href="{[=hp.HttpSrvBasePath('auth/sign-out')]}">Sign out</a>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.NavActive("hp-topbar-nav-main", "{{.http_request_path}}");
    hp.AuthSessionRefresh();
});
</script>
