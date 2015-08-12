
<div id="l5smgr-spec-actionls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>Datax</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-spec-actionls"></tbody>
</table>

<script id="l5smgr-spec-actionls-tpl" type="text/html">  
{[~it.actions :v]}
<tr>
  <td class="l5smgr-font-fixspace">{[=v.name]}</td>
  <td>{[=v._dataxNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="l5sSpec.ActionSet('{[=it._modname]}', '{[=v.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
