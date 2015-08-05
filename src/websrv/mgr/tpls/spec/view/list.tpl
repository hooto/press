
<div id="l5smgr-spec-viewls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Template</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-spec-viewls">
{[~it.items :v]}
<tr>
  <td><strong>{[=v.path]}</strong></td>
  <td align="right">
    <button class="btn btn-default" onclick="l5sSpec.RouteSetTemplateSelectOne('{[=v.path]}')">Select</button>
  </td>
</tr>
{[~]}
  </tbody>
</table>

<script id="l5smgr-spec-viewls-tpl" type="text/html">  

</script>
