
<div id="l5smgr-specls-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Name</th>
      <th>SrvName</th>
      <th>Title</th>
      <th>Version</th>
      <th>Nodes</th>
      <th>Terms</th>
      <th>Actions</th>
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-specls"></tbody>
</table>

<script id="l5smgr-specls-tpl" type="text/html">  
{[~it.items :v]}
<tr>
  <td class="l5smgr-font-fixspace">{[=v.meta.name]}</td>
  <td class="l5smgr-font-fixspace">{[=v.srvname]}</td>
  <td>{[=v.title]}</td>
  <td>{[=v.meta.resourceVersion]}</td>
  <td>{[=v._nodeModelsNum]}</td>
  <td>{[=v._termModelsNum]}</td>
  <td>{[=v._actionsNum]}</td>
  <td>{[=l4i.TimeParseFormat(v.meta.created, "Y-m-d")]}</td>
  <td>{[=l4i.TimeParseFormat(v.meta.updated, "Y-m-d H:i")]}</td>
  <td align="right">
    <a class="xl37hg btn btn-default btn-xs" href="#{[=v.meta.id]}">Setting</a>
  </td>
</tr>
{[~]}
</script>

<script type="text/javascript">

$("#l5smgr-specls").on("click", ".xl37hg", function() {
    var id = $(this).attr("href").substr(1);
    l5sSpec.InfoSet(id);
});

</script>
