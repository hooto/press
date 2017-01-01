
<div id="htpm-spec-viewls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Template</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htpm-spec-viewls">
{[~it.items :v]}
<tr>
  <td><strong>{[=v.path]}</strong></td>
  <td align="right">
    <button class="btn btn-default" onclick="htpSpec.RouteSetTemplateSelectOne('{[=v.path]}')">Select</button>
  </td>
</tr>
{[~]}
  </tbody>
</table>

<script id="htpm-spec-viewls-tpl" type="text/html">  

</script>
