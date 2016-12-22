
<div id="htapm-specls-alert"></div>

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
  <tbody id="htapm-specls"></tbody>
</table>

<script id="htapm-specls-tpl" type="text/html">  
{[~it.items :v]}
<tr>
  <td><a class="htapm-specls-infoset-item" href="#{[=v.meta.name]}">{[=v.title]}</a></td>
  <td class="htapm-font-fixspace">{[=v.meta.name]}</td>
  <td class="htapm-font-fixspace">{[=v.srvname]}</td>
  <td>{[=v.meta.resourceVersion]}</td>
  <td><button class="btn btn-default btn-sm" onclick="htapSpec.NodeList('{[=v.meta.name]}')">{[=v._nodeModelsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="htapSpec.ActionList('{[=v.meta.name]}')">{[=v._actionsNum]}</button></td>
  <td><button class="btn btn-default btn-sm" onclick="htapSpecEditor.Open('{[=v.meta.name]}')">Editor</button></td>
  <td><button class="btn btn-default btn-sm" onclick="htapSpec.RouteList('{[=v.meta.name]}')">{[=v._routesNum]}</button></td>
  <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
  <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d")]}</td>
</tr>
{[~]}
</script>

<script type="text/javascript">

$("#htapm-specls").on("click", ".htapm-specls-infoset-item", function() {
    var name = $(this).attr("href").substr(1);
    htapSpec.InfoSet(name);
});

</script>
