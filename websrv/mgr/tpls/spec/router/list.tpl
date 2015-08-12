
<div id="l5smgr-spec-routels-alert"></div>

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
  <tbody id="l5smgr-spec-routels"></tbody>
</table>

<script id="l5smgr-spec-routels-tpl" type="text/html">  
{[~it.router.routes :v]}
<tr>
  <td class="l5smgr-font-fixspace">{[=v.path]}</td>
  <td>{[=v.dataAction]}</td>
  <td>{[=v.template]}</td>
  <td>{[=v._paramsNum]}</td>
  <td align="right">
    <button class="btn btn-default" onclick="l5sSpec.RouteSet('{[=it._modname]}', '{[=v.path]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
