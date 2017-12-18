
<div id="hpm-termls"></div>

<div id="hpm-termls-pager"></div>

<div id="hpm-node-term-opts" class="hpm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="hpTerm.Set()" id="hpm-term-list-new-title">
      New Term
    </a>
  </li>
  <li>
    <form onsubmit="hpTerm.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control hpm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>  
</div>


<script id="hpm-termls-tpl" type="text/html">  
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
      <td class="hpm-font-fixspace">{[=v.id]}</td>
      <td>{[=v.title]}</td>
      {[ if (it.model.type == "taxonomy") { ]}
      <td>{[=v.weight]}</td>
      {[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
	  <td align="right">
        <button class="pure-button button-xsmall" onclick="hpTerm.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v.id]}')">Edit</button>
	  </td>
    </tr>
    {[? v._subs]}
    {[~v._subs :v2]}
    <tr>
      <td class="hpm-font-fixspace">{[=v2.id]}</td>
      <td>{[=hpTerm.Sprint(v2._dp)]}{[=v2.title]}</td>
      <td>{[=hpTerm.Sprint(v2._dp)]}{[=v2.weight]}</td>
      <td>{[=v2.created]}</td>
      <td>{[=v2.updated]}</td>
      <td align="right">
        <button class="pure-button button-xsmall" onclick="hpTerm.Set('{[=it.modname]}', '{[=it.modelid]}', '{[=v2.id]}')">Edit</button>
	  </td>
    </tr>
    {[~]}
    {[?]}
  {[ } ]}
  {[~]}
  </tbody>
</table>
</script>

<script id="hpm-termls-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="hpm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="hpTerm.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="hpTerm.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="hpTerm.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#hpm-termls").on("click", ".term-item", function() {
    var id = $(this).attr("href").substr(1);
    hpTerm.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

</script>
