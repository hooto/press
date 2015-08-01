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
      <th>Title</th>
      <th>Status</th>
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-nodels"></tbody>
</table>

<div id="l5smgr-nodels-pager"></div>


<script id="l5smgr-nodels-tpl" type="text/html">  
  {[~it.items :v]}
    <tr>
      <td>
        <a class="node-item" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">{[=v.title]}</a>
      </td>
      <td>{[=v.status]}</td>
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      <td align="right">
        <a class="node-item" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Edit</a>
      </td>
    </tr>
  {[~]}
</script>

<script id="l5smgr-nodels-pager-tpl" type="text/html">
<ul class="pagination pagination-sm">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="l5sNode.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="l5sNode.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="l5sNode.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
</script>

<script type="text/javascript">

$("#l5smgr-nodels").on("click", ".node-item", function() {
    var id = $(this).attr("href").substr(1);
    l5sNode.Set($(this).attr("modname"), $(this).attr("modelid"), id);
    l4iStorage.Set("l5smgr_nodels_page", 0);
});

</script>
