
<div id="hpm-spec-routels-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Path</th>
      <th>Data Action</th>
      <th>Template</th>
      <th>Params</th>
      <th>Default</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpm-spec-routels"></tbody>
</table>

<script id="hpm-spec-routels-tpl" type="text/html">  
{[~it.router.routes :v]}
<tr>
  <td class="hpm-font-fixspace">{[=v.path]}</td>
  <td>{[=v.dataAction]}</td>
  <td>{[=v.template]}</td>
  <td>{[=v._paramsNum]}</td>
  <td>
    {[if (v.default) {]}Yes{[} else {]}No{[}]}
  </td>
  <td align="right">
    <button class="btn btn-default" onclick="hpSpec.RouteSet('{[=it._modname]}', '{[=v.path]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
