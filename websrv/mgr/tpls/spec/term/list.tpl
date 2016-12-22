
<div id="htapm-spec-termls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Title</th>
      <th>Name</th>
      <th>Type</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htapm-spec-termls"></tbody>
</table>

<script id="htapm-spec-termls-tpl" type="text/html">  
{[~it.termModels :v]}
<tr>
  <td>{[=v.title]}</td>
  <td class="htapm-font-fixspace">{[=v.meta.name]}</td>
  <td>{[=v.type]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="htapSpec.TermSet('{[=it.meta.name]}', '{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
