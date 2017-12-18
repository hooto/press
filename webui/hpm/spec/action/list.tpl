
<div id="hpm-spec-actionls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>Datax</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpm-spec-actionls"></tbody>
</table>

<script id="hpm-spec-actionls-tpl" type="text/html">  
{[~it.actions :v]}
<tr>
  <td class="hpm-font-fixspace">{[=v.name]}</td>
  <td>{[=v._dataxNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="hpSpec.ActionSet('{[=it._modname]}', '{[=v.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
