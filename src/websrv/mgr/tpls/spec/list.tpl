
<div id="l5smgr-specls-alert">ff</div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>ID</th>
      <th>Name</th>
      <th>Version</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-specls"></tbody>
</table>

<script id="l5smgr-specls-tpl" type="text/html">  
  {[~it.items :v]}
    <tr>
      <td>{[=v.metadata.id]}</td>
      <td>{[=v.title]}</td>
      <td>{[=v.metadata.resourceVersion]}</td>
      <td>{[=v.metadata.updated]}</td>
      <td><a class="xl37hg" href="#{[=v.metadata.id]}">Setting</a>
    </tr>
  {[~]}
</script>

<script type="text/javascript">

$("#l5smgr-specls").on("click", ".xl37hg", function() {
    var id = $(this).attr("href").substr(1);
    // l5yLps.ChannelSet(id);
});

</script>
