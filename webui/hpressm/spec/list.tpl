
<div id="hpressm-specls-alert"></div>

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
  <tbody id="hpressm-specls"></tbody>
</table>

<script id="hpressm-specls-tpl" type="text/html">  
{[~it.items :v]}
<tr>
  <td>{[=v.title]}</td>
  <td class="hpressm-font-fixspace">{[=v.meta.name]}</td>
  <td class="hpressm-font-fixspace">{[=v.srvname]}</td>
  <td>{[=v.meta.resourceVersion]}</td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.NodeList('{[=v.meta.name]}')">{[=v._nodeModelsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.ActionList('{[=v.meta.name]}')">{[=v._actionsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.RouteList('{[=v.meta.name]}')">{[=v._routesNum]}</button></td>
  <td>
    {[if (v.status) {]}
      <span class="label label-success">Enable</span>
    {[} else {]}
      <span class="label label-default">Disable</span>
    {[}]}
  </td>
  <td align="right">
    <button class="btn btn-default btn-sm" onclick="hpressSpecEditor.Open('{[=v.meta.name]}')">Develop</button>
    <button class="btn btn-default btn-sm" onclick="hpressSpec.InfoSet('{[=v.meta.name]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>

