
<div id="hpm-specls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Title</th>
      <th>Name</th>
      <th>Service Name</th>      
      <th>Version</th>
      <th>Nodes</th>
      <th>Actions</th>
      <th>Routes</th>
      <th>Status</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpm-specls"></tbody>
</table>

<script id="hpm-specls-tpl" type="text/html">  
{[~it.items :v]}
<tr>
  <td>{[=v.title]}</td>
  <td class="hpm-font-fixspace">{[=v.meta.name]}</td>
  <td class="hpm-font-fixspace">{[=v.srvname]}</td>
  <td>{[=v.meta.version]}</td>
  <td><button class="btn btn-default btn-sm" onclick="hpSpec.NodeList('{[=v.meta.name]}')">{[=v._nodeModelsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpSpec.ActionList('{[=v.meta.name]}')">{[=v._actionsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpSpec.RouteList('{[=v.meta.name]}')">{[=v._routesNum]}</button></td>
  <td>
    {[if (v.status) {]}
      <span class="label label-success">Enable</span>
    {[} else {]}
      <span class="label label-default">Disable</span>
    {[}]}
  </td>
  <td align="right">
    <button class="btn btn-default btn-sm" onclick="hpSpecEditor.Open('{[=v.meta.name]}')">Develop</button>
    <button class="btn btn-default btn-sm" onclick="hpSpec.InfoSet('{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>

