
<div id="htapm-spec-nodels-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>Title</th>
      <th>Fields</th>
      <th>Terms</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htapm-spec-nodels"></tbody>
</table>

<script id="htapm-spec-nodels-tpl" type="text/html">  
{[~it.nodeModels :v]}
<tr>
  <td class="htapm-font-fixspace">{[=v.meta.name]}</td>
  <td>{[=v.title]}</td>
  <td>{[=v._fieldsNum]}</td>
  <td>{[=v._termsNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="htapSpec.NodeSet('{[=it.meta.name]}', '{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
