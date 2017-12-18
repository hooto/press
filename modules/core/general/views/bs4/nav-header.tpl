<nav class="navbar navbar-expand-lg navbar-light bg-light" style="sbackground-color: #1890FF;" id="hp-topbar">
<div class="container">
  {{if SysConfig "frontend_header_site_logo_url"}}
  <span class="navbar-brand">
    <img src="{{SysConfig "frontend_header_site_logo_url"}}" height="30" alt="">
  </span>
  {{end}}
  <span class="navbar-brand mb-0 h1">{{SysConfig "frontend_header_site_name"}}</span>

  <div class="collapse navbar-collapse">
    <ul class="navbar-nav mr-auto mt-2 mt-lg-0" id="hp-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <li class="nav-item"><a class="nav-link" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a></li>
      {{end}}
    </ul>
    <div class="form-inline my-2 my-lg-0" id="hp-topbar-userbar"></div>
  </div>
</div>
</nav>

<script id="hp-topbar-user-unsigned-tpl" type="text/html">
<span class="iam-name"><a href="{{HttpSrvBasePath "hp/auth/login"}}">Sign In</a></span>
</script>

<script id="hp-topbar-user-signed-tpl" type="text/html">
<span class="iam-name">{[=it.display_name]}</span>
<span class="iam-photo" id="hp-topbar-user-signed"><img src="{[=it.photo_url]}"/></span>
<div id="hp-topbar-user-signed-modal" style="display:none;">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.display_name]}</div>
  {[? it.instance_owner]}<a class="btn btn-outline-primary iam-btn" href="{[=hp.HttpSrvBasePath('mgr')]}" target="_blank">Content Manage</a>{[?]}
  <a class="btn btn-outline-secondary iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="btn btn-outline-secondary iam-btn" href="{[=hp.HttpSrvBasePath('auth/sign-out')]}">Sign out</a>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.NavActive("hp-topbar-nav-main", "{{.http_request_path}}");
    hp.AuthSessionRefresh();
});
</script>
