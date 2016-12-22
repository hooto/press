
<div id="htapm-spec-routels-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>Data Action</th>
      <th>Template</th>
      <th>Params</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htapm-spec-routels"></tbody>
</table>

<script id="htapm-spec-routels-tpl" type="text/html">  
{[~it.router.routes :v]}
<tr>
  <td class="htapm-font-fixspace">{[=v.path]}</td>
  <td>{[=v.dataAction]}</td>
  <td>{[=v.template]}</td>
  <td>{[=v._paramsNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="htapSpec.RouteSet('{[=it._modname]}', '{[=v.path]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
