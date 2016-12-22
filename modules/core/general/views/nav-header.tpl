<div id="htap-topbar">
  <div class="container htap-topbar-collapse">
    <ul class="htap-nav">
      {{if SysConfig "frontend_header_site_logo_url"}}
      <li><img class="htap-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}"></li>
      {{end}}
      <li class="htap-topbar-brand">{{SysConfig "frontend_header_site_name"}}</li>
    </ul>

    <ul class="htap-nav htap-topbar-nav" id="htap-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <li><a class="" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a></li>
      {{end}}
    </ul>

    <ul class="htap-nav htap-nav-right" id="htap-topbar-userbar">
      <li class="iam-name"><a href="{{HttpSrvBasePath "auth/login"}}">Login</a></li>
    </ul>
  </div>
</div>

<script id="htap-topbar-user-signed-tpl" type="text/html">

<li class="iam-name">{[=it.name]}</li>
<li class="iam-photo" id="htap-topbar-user-signed"><img src="{[=it.photo_url]}"/></li>

<div id="htap-topbar-user-signed-modal" style="display:none;">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.name]}</div>
  <a class="btn btn-primary iam-btn" href="{[=htap.HttpSrvBasePath("mgr")]}">Content Manage</a>
  <a class="btn btn-default iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="btn btn-default iam-btn" href="{[=htap.HttpSrvBasePath("auth/sign-out")]}">Sign out</a>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    htap.NavActive("htap-topbar-nav-main", "{{.baseuri}}");
    htap.AuthSessionRefresh();
});
</script>
