
<div id="hpressm-spec-viewls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Template</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpressm-spec-viewls">
{[~it.items :v]}
<tr>
  <td><strong>{[=v.path]}</strong></td>
  <td align="right">
    <button class="btn btn-default" onclick="hpressSpec.RouteSetTemplateSelectOne('{[=v.path]}')">Select</button>
  </td>
</tr>
{[~]}
  </tbody>
</table>

<script id="hpressm-spec-viewls-tpl" type="text/html">  

</script>
