<table class="table table-hover" id="hpm-nodels"></table>

<div id="hpm-nodels-pager"></div>

<div id="hpm-node-list-opts" class="hpm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary" id="hpm-node-list-refer-back" style="display:none">
    <a href="#" onclick="hpNode.ReferBack()">
      Back
    </a>
  </li>
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="hpNode.Set()" id="hpm-node-list-new-title">
      New Content
    </a>
  </li>
  <li>
    <form onsubmit="hpNode.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control hpm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>
  <li id="hpm-nodels-batch-select-todo-btn" 
    class="pure-button btapm-btn btapm-btn-primary" style="display:none">
    <a href="#" onclick="hpNode.ListBatchSelectTodo()">
      Select Contents todo ...
    </a>
  </li>
</div>

<script id="hpm-nodels-tpl" type="text/html">  
  <thead>
    <tr>
      <th width="20">
        <input class="row-checkbox hpm-nodels-chk-all" type="checkbox" onclick="hpNode.ListBatchSelectAll()">
      </th>
      <th>Title</th>
      <th>Status</th>
      {[if (it.model.extensions.access_counter) { ]}<th>Access</th>{[ } ]}
      <th>Created</th>
      <th>Updated</th>
      {[if (it.model.extensions.node_sub_refer) {]}
      <th></th>
	  {[}]}
      <th></th>
    </tr>
  </thead>
  <tbody>
  {[~it.items :v]}
    <tr>
      <td>
        <input class="row-checkbox hpm-nodels-chk-item" type="checkbox" value="{[=v.id]}"
          onclick="hpNode.ListBatchSelectTodoBtnRefresh()">
      </td>
      <td>
        <a class="node-item" onclick="hpNode.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v.id]}')" href="#{[=v.id]}">{[=v.title]}</a>
      </td>
      <td>
      {[~it._status_def :sv]}
        {[if (sv.type == v.status) { ]}{[=sv.name]}{[ } ]}
      {[~]}
      </td>
      {[if (it.model.extensions.access_counter) { ]}<td>{[=v.ext_access_counter]}</td>{[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      {[if (it.model.extensions.node_sub_refer) {]}
      <td>
        <!--<button class="pure-button button-xsmall" onclick="hpNode.Set('{[=it.modname]}', '{[=it.model.extensions.node_sub_refer]}', null, '{[=v.id]}')">New Sub Content</button>-->
        <button class="pure-button button-xsmall" onclick="hpNode.List('{[=it.modname]}', '{[=it.model.extensions.node_sub_refer]}', '{[=v.id]}')">Sub Contents</button>
      </td>
      {[}]}
      <td align="right">
        <button class="pure-button button-xsmall" onclick="hpNode.Del('{[=it.modname]}', '{[=it.modelid]}', '{[=v.id]}')">Delete</button>
        <button class="pure-button button-xsmall" onclick="hpNode.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v.id]}')">Edit</button>
      </td>
    </tr>
  {[~]}
  </tbody>
</script>

<script id="hpm-nodels-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="hpm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="hpNode.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="hpNode.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="hpNode.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

