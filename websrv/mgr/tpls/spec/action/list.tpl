
<div id="htapm-spec-actionls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>Datax</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htapm-spec-actionls"></tbody>
</table>

<script id="htapm-spec-actionls-tpl" type="text/html">  
{[~it.actions :v]}
<tr>
  <td class="htapm-font-fixspace">{[=v.name]}</td>
  <td>{[=v._dataxNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="htapSpec.ActionSet('{[=it._modname]}', '{[=v.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
