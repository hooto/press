
<table width="100%">
  <tr>
    <td>
      <form id="lps-infols-qry" action="#" class="form-inlines">
        <input id="lps-infols-qry-text" type="text"
          class="form-control l5smgr-query-input" 
          placeholder="Press Enter to Search" 
          value="">
      </form>
    </td>
    <td align="right">
      <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="l5sNode.Set()">
        New Content
      </button>
    </td>
  </tr>
</table>

<table class="table table-hover">
  <thead>
    <tr>
      <th>ID</th>
      <th>Title</th>
      <th>State</th>
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-nodels"></tbody>
</table>

<script id="l5smgr-nodels-tpl" type="text/html">  
  {[~it.items :v]}
    <tr>
      <td class="l5smgr-font-wfix"><a class="node-item" specid="{[=it.specid]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">{[=v.id]}</a></td>
      <td>{[=v.title]}</td>
      <td>{[=v.state]}</td>
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      <td align="right">
        <a class="node-item" specid="{[=it.specid]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Edit</a>
      </td>
    </tr>
  {[~]}
</script>

<script type="text/javascript">

$("#l5smgr-nodels").on("click", ".node-item", function() {
    var id = $(this).attr("href").substr(1);
    l5sNode.Set($(this).attr("specid"), $(this).attr("modelid"), id);
});

</script>
