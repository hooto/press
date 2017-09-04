<table class="table table-hover" id="hpressm-nodels"></table>

<div id="hpressm-nodels-pager"></div>

<div id="hpressm-node-list-opts" class="hpressm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="hpressNode.Set()">
      New Content
    </a>
  </li>
  <li>
    <form onsubmit="hpressNode.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control hpressm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>
  <li id="hpressm-nodels-batch-select-todo-btn" 
    class="pure-button btapm-btn btapm-btn-primary" style="display:none">
    <a href="#" onclick="hpressNode.ListBatchSelectTodo()">
      Select Contents todo ...
    </a>
  </li>
</div>

<script id="hpressm-nodels-tpl" type="text/html">  
  <thead>
    <tr>
      <th width="20">
        <input class="row-checkbox hpressm-nodels-chk-all" type="checkbox" onclick="hpressNode.ListBatchSelectAll()">
      </th>
      <th>Title</th>
      <th>Status</th>
      {[if (it.model.extensions.access_counter) { ]}<th>Access</th>{[ } ]}
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody>
  {[~it.items :v]}
    <tr>
      <td>
        <input class="row-checkbox hpressm-nodels-chk-item" type="checkbox" value="{[=v.id]}"
          onclick="hpressNode.ListBatchSelectTodoBtnRefresh()">
      </td>
      <td>
        <a class="node-item" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">{[=v.title]}</a>
      </td>
      <td>
      {[~it._status_def :sv]}
        {[if (sv.type == v.status) { ]}{[=sv.name]}{[ } ]}
      {[~]}
      </td>
      {[if (it.model.extensions.access_counter) { ]}<td>{[=v.ext_access_counter]}</td>{[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      <td align="right">
        <a class="node-item-del btn btn-default btn-xs" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Del</a>
        <a class="node-item btn btn-default btn-xs" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Edit</a>
      </td>
    </tr>
  {[~]}
  </tbody>
</script>

<script id="hpressm-nodels-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="hpressm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="hpressNode.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="hpressNode.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="hpressNode.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#hpressm-nodels").on("click", ".node-item", function() {
    var id = $(this).attr("href").substr(1);
    hpressNode.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

$("#hpressm-nodels").on("click", ".node-item-del", function() {
    var id = $(this).attr("href").substr(1);
    hpressNode.Del($(this).attr("modname"), $(this).attr("modelid"), id);
});
</script>
