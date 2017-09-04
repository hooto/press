
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
      <th>Views</th>
      <th>Routes</th>
      <th>Created</th>
      <th>Updated</th>
    </tr>
  </thead>
  <tbody id="hpressm-specls"></tbody>
</table>

<script id="hpressm-specls-tpl" type="text/html">  
{[~it.items :v]}
<tr>
  <td><a class="hpressm-specls-infoset-item" href="#{[=v.meta.name]}">{[=v.title]}</a></td>
  <td class="hpressm-font-fixspace">{[=v.meta.name]}</td>
  <td class="hpressm-font-fixspace">{[=v.srvname]}</td>
  <td>{[=v.meta.resourceVersion]}</td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.NodeList('{[=v.meta.name]}')">{[=v._nodeModelsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.ActionList('{[=v.meta.name]}')">{[=v._actionsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpecEditor.Open('{[=v.meta.name]}')">Editor</button></td>
  <td><button class="btn btn-default btn-sm" onclick="hpressSpec.RouteList('{[=v.meta.name]}')">{[=v._routesNum]}</button></td>
  <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
  <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
</tr>
{[~]}
</script>

<script type="text/javascript">

$("#hpressm-specls").on("click", ".hpressm-specls-infoset-item", function() {
    var name = $(this).attr("href").substr(1);
    hpressSpec.InfoSet(name);
});

</script>
