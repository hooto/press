<table class="table table-hover" id="htpm-nodels"></table>

<div id="htpm-nodels-pager"></div>

<div id="htpm-node-list-opts" class="htpm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="htpNode.Set()">
      New Content
    </a>
  </li>
  <li>
    <form onsubmit="htpNode.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control htpm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>
</div>

<script id="htpm-nodels-tpl" type="text/html">  
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

<script id="htpm-nodels-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="htpm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="htpNode.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="htpNode.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="htpNode.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#htpm-nodels").on("click", ".node-item", function() {
    var id = $(this).attr("href").substr(1);
    htpNode.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

$("#htpm-nodels").on("click", ".node-item-del", function() {
    var id = $(this).attr("href").substr(1);
    htpNode.Del($(this).attr("modname"), $(this).attr("modelid"), id);
});
</script>
