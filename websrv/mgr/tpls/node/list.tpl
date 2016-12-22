<table class="table table-hover" id="htapm-nodels"></table>

<div id="htapm-nodels-pager"></div>

<div id="htapm-node-list-opts" class="htapm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="htapNode.Set()">
      New Content
    </a>
  </li>
  <li>
    <form onsubmit="htapNode.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control htapm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>
</div>

<script id="htapm-nodels-tpl" type="text/html">  
  <thead>
    <tr>
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

<script id="htapm-nodels-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="htapm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="htapNode.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="htapNode.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="htapNode.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#htapm-nodels").on("click", ".node-item", function() {
    var id = $(this).attr("href").substr(1);
    htapNode.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

$("#htapm-nodels").on("click", ".node-item-del", function() {
    var id = $(this).attr("href").substr(1);
    htapNode.Del($(this).attr("modname"), $(this).attr("modelid"), id);
});
</script>
