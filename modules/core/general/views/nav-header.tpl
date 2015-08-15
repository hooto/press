<div id="l5s-topbar">
<div class="container">
<table width="100%">
  <tr>
    <td align="left" class="l5s-topbar-layout">
      <div class="l5s-topbar-sets">
        {{if SysConfig "frontend_header_site_logo_url"}}
        <img class="l5s-topbar-logo" src="{{SysConfig "frontend_header_site_logo_url"}}">
        {{end}}
        <span class="l5s-topbar-brand">{{SysConfig "frontend_header_site_name"}}</span>
      </div>

      <div class="l5s-topbar-sets l5s-topbar-nav" id="l5s-topbar-nav-main">
        {{range $v := .topnav.Items}}
        <a class="" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a>
        {{end}}
      </div>
    </td>

    <td align="right">
      <div id="l5s-topvar-user-box" classs="l5s-topbar-sets l5s-topbar-nav">
        <a href="/auth/login">Login</a>
      </div>
    </td>
  </tr>
</table>
</div>
</div>

<script id="l5s-topvar-user-box-tpl" type="text/html">  
  <div id="l5s-topvar-user-box">
    <span class="lunb-name">{[=it.name]}</span>
    <span><img class="lnub-photo" src="{[=it.photo_url]}" /></span>
  </div>
  <div id="l5s-topvar-user-modal" style="display:none;">
    <img class="lnum-photo" src="{[=it.photo_url]}">
    <div class="lnum-name">{[=it.name]}</div>
    <a class="btn btn-primary lnum-btn" href="/mgr">Content Manage</a>
    <a class="btn btn-default lnum-btn" href="{[=it.ids_url]}" target="_blank">Account Center</a>
    <a class="btn btn-default lnum-btn" href="/auth/sign-out">Sign out</a>
  </div>
</script>

<script type="text/javascript">
window.onload_hooks.push(function() {
    l5s.NavActive("l5s-topbar-nav-main", "{{.baseuri}}");
    l5s.AuthSessionRefresh();
});
</script>