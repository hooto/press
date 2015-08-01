<div id="l5s-topbar">
<div class="container">
<table width="100%">
  <tr>
    <td align="left" class="l5s-topbar-layout">
      <div class="l5s-topbar-sets">
        <div class="l5s-topbar-brand">Project Name</div>
      </div>

      <div class="l5s-topbar-sets l5s-topbar-nav" id="l5s-topbar-nav-main">
        {{range $v := .topnav.Items}}
        <a class="" href="{{FieldString $v.Fields "url"}}">{{$v.Title}}</a>
        {{end}}
      </div>
    </td>

    <td align="right" class="">
      <div class="l5s-topbar-sets l5s-topbar-nav">
        <a class="" href="#">Login</a>
      </div>
    </td>
  </tr>
</table>
</div>
</div>

<script type="text/javascript">
window.onload_hooks.push(function() {
    l5s.NavActive("l5s-topbar-nav-main", "{{.baseuri}}");
});
</script>