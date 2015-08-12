
<div id="l5smgr-spec-termls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Title</th>
      <th>Name</th>
      <th>Type</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-spec-termls"></tbody>
</table>

<script id="l5smgr-spec-termls-tpl" type="text/html">  
{[~it.termModels :v]}
<tr>
  <td>{[=v.title]}</td>
  <td class="l5smgr-font-fixspace">{[=v.meta.name]}</td>
  <td>{[=v.type]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="l5sSpec.TermSet('{[=it.meta.name]}', '{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
