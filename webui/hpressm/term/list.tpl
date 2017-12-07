
<div id="hpressm-termls"></div>

<div id="hpressm-termls-pager"></div>

<div id="hpressm-node-term-opts" class="hpressm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="hpressTerm.Set()" id="hpressm-term-list-new-title">
      New Term
    </a>
  </li>
  <li>
    <form onsubmit="hpressTerm.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control hpressm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>  
</div>


<script id="hpressm-termls-tpl" type="text/html">  
<table class="table table-hover">
  <thead>
    <tr>
      <th width="80px">ID</th>
      <th>Title</th>
      {[ if (it.model.type == "taxonomy") { ]}
      <th>Weight</th>
      {[ } ]}
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody>
  {[~it.items :v]}
  {[ if (v.pid == 0) { ]}
    <tr>
      <td class="hpressm-font-fixspace">{[=v.id]}</td>
      <td>{[=v.title]}</td>
      {[ if (it.model.type == "taxonomy") { ]}
      <td>{[=v.weight]}</td>
      {[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
	  <td align="right">
        <button class="pure-button button-xsmall" onclick="hpressTerm.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v.id]}')">Edit</button>
	  </td>
    </tr>
    {[? v._subs]}
    {[~v._subs :v2]}
    <tr>
      <td class="hpressm-font-fixspace">{[=v2.id]}</td>
      <td>{[=hpressTerm.Sprint(v2._dp)]}{[=v2.title]}</td>
      <td>{[=hpressTerm.Sprint(v2._dp)]}{[=v2.weight]}</td>
      <td>{[=v2.created]}</td>
      <td>{[=v2.updated]}</td>
      <td align="right">
        <button class="pure-button button-xsmall" onclick="hpressTerm.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v2.id]}')">Edit</button>
	  </td>
    </tr>
    {[~]}
    {[?]}
  {[ } ]}
  {[~]}
  </tbody>
</table>
</script>

<script id="hpressm-termls-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="hpressm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="hpressTerm.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="hpressTerm.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="hpressTerm.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#hpressm-termls").on("click", ".term-item", function() {
    var id = $(this).attr("href").substr(1);
    hpressTerm.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

</script>
