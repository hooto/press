
<div id="htpm-termls"></div>

<div id="htpm-termls-pager"></div>

<div id="htpm-node-term-opts" class="htpm-hide">
  <li class="pure-button btapm-btn btapm-btn-primary">
    <a href="#" onclick="htpTerm.Set()">
      New Term
    </a>
  </li>
  <li>
    <form onsubmit="htpTerm.List(); return false;" action="#" class="form-inlines">
      <input id="qry_text" type="text"
        class="form-control htpm-query-input" 
        placeholder="Press Enter to Search" 
        value="">
    </form>
  </li>  
</div>


<script id="htpm-termls-tpl" type="text/html">  
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
      <td class="htpm-font-fixspace">{[=v.id]}</td>
      <td>{[=v.title]}</td>
      {[ if (it.model.type == "taxonomy") { ]}
      <td>{[=v.weight]}</td>
      {[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      <td align="right"><a class="term-item btn btn-default btn-xs" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Edit</a>
    </tr>
    {[? v._subs]}
    {[~v._subs :v2]}
    <tr>
      <td class="htpm-font-fixspace">{[=v2.id]}</td>
      <td>{[=htpTerm.Sprint(v2._dp)]}{[=v2.title]}</td>
      <td>{[=htpTerm.Sprint(v2._dp)]}{[=v2.weight]}</td>
      <td>{[=v2.created]}</td>
      <td>{[=v2.updated]}</td>
      <td align="right"><a class="term-item btn btn-default btn-xs" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v2.id]}">Edit</a>
    </tr>
    {[~]}
    {[?]}
  {[ } ]}
  {[~]}
  </tbody>
</table>
</script>

<script id="htpm-termls-pager-tpl" type="text/html">
{[ if (it.RangePages.length > 1) { ]}
<ul class="htpm-pager">
  {[ if (it.FirstPageNumber > 0) { ]}
  <li><a href="#{[=it.FirstPageNumber]}" onclick="htpTerm.ListPage({[=it.FirstPageNumber]})">First</a></li>
  {[ } ]}
  {[~it.RangePages :v]}
  <li {[ if (v == it.CurrentPageNumber) { ]}class="active"{[ } ]}><a href="#{[=v]}" onclick="htpTerm.ListPage({[=v]})">{[=v]}</a></li>
  {[~]}
  {[ if (it.LastPageNumber > 0) { ]}
  <li><a href="#{[=it.LastPageNumber]}" onclick="htpTerm.ListPage({[=it.LastPageNumber]})">Last</a></li>
  {[ } ]}
</ul>
{[ } ]}
</script>

<script type="text/javascript">

$("#htpm-termls").on("click", ".term-item", function() {
    var id = $(this).attr("href").substr(1);
    htpTerm.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

</script>
