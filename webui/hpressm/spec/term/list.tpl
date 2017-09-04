
<div id="hpressm-spec-termls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Title</th>
      <th>Name</th>
      <th>Type</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpressm-spec-termls"></tbody>
</table>

<script id="hpressm-spec-termls-tpl" type="text/html">  
{[~it.termModels :v]}
<tr>
  <td>{[=v.title]}</td>
  <td class="hpressm-font-fixspace">{[=v.meta.name]}</td>
  <td>{[=v.type]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="hpressSpec.TermSet('{[=it.meta.name]}', '{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
