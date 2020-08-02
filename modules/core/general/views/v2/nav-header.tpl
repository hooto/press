<nav class="navbar is-light" id="hp-topbar">
<div class="container">
  <div class="navbar-brand">
   {{if SysConfig "frontend_header_site_logo_url"}}
    <a class="navbar-item" href="#">
      <img src="{{SysConfig "frontend_header_site_logo_url"}}"  height="30" alt="">
      {{SysConfig "frontend_header_site_name"}}
    </a>
    {{end}}
    <a role="button" class="navbar-burger" onclick="hp.NavbarMenuToggle('hpex-topbar-nav-main')">
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
    </a>
  </div>
 
  <div class="navbar-menu" id="hpex-topbar-nav-main">
    <div class="navbar-start" id="hp-topbar-nav-main">
      {{range $v := .topnav.Items}}
      <a class="navbar-item is-tab" href="{{FieldString $v.Fields "url"}}">{{FieldStringPrint $v "title" $.LANG}}</a>
      {{end}}
    </div>
    <div class="navbar-end">
      {{if $.frontend_langs}}
      <div class="navbar-item hp-footer-powerby-item">
        <select onclick="hp.LangChange(this)" class="hp-footer-langs">
        {{range $v := $.frontend_langs}}
        <option value="{{$v.Id}}" {{if eq $v.Id $.LANG}}selected{{end}}>{{$v.Name}}</option>
        {{end}}
        </select>
      </div>
      {{end}}
      <div class="navbar-item navbar-link" id="hp-topbar-userbar" onclick="hp.NavbarMenuUserToggle()"></div>
    </div>
  </div>
</div>
</nav>

<script id="hp-topbar-user-unsigned-tpl" type="text/html">
<a class="navbar-item" href="{{HttpSrvBasePath "hp/auth/login"}}">Sign In</a>
</script>

<script id="hp-topbar-user-signed-tpl" type="text/html">

<span class="status-name">{[=it.display_name]}</span>
<img id="hp-topbar-user-signed" class="status-photo" src="{[=it.photo_url]}" width="28" height="28"/>

<div id="hp-topbar-user-signed-modal" style="display:none;">
<div class="hp-topbar-user-signed-modal" onclick="hp.NavbarMenuUserClose()">
  <img class="iam-photo" src="{[=it.photo_url]}">
  <div class="iam-name">{[=it.display_name]}</div>
  {[? it.instance_owner]}
  <a class="button is-fullwidth iam-btn" href="{[=hp.HttpSrvBasePath('mgr')]}" target="_blank">Content Manage</a>
  {[?]}
  <a class="button is-fullwidth iam-btn" href="{[=it.iam_url]}" target="_blank">Account Center</a>
  <a class="button is-fullwidth" href="{[=hp.HttpSrvBasePath('auth/sign-out')]}">Sign Out</a>
</div>
</div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    hp.NavActive("hp-topbar-nav-main", "{{.http_request_path}}");
    hp.AuthSessionRefresh();
});
</script>
