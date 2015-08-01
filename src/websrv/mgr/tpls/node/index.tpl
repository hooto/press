
<div id="l5smgr-node-specls" style="padding-bottom:10px;border-bottom:1px solid #ccc"></div>

<script id="l5smgr-node-specls-tpl" type="text/html">
<ul class="nav nav-pills">
  {[~it.items :v]}
  <li class="spec-item {[if (it.active == v.meta.name) {]}active{[}]}" tgspec="{[=v.meta.name]}"><a href="#{[=v.meta.name]}">{[=v.title]}</a></li>
  {[~]}
</ul>
</script>

<table width="100%" style="margin-top:10px">
<tr>
  <td width="180px" valign="top">
    <div id="l5smgr-node-nmodels"></div>
    <!-- <div style="border-bottom: 1px solid #ccc"></div> -->
    <div id="l5smgr-node-tmodels"></div>
  </td>
  <td width="20px"></td>
  <td valign="top">
    <div id="l5smgr-node-alert"></div>
    <div id="work-content"></div>
  </td>
</tr>
</table>

<script id="l5smgr-node-nmodels-tpl" type="text/html">
<ul class="nav nav-pills nav-stacked">
  {[~it.items :v]}
  <li class="node-item {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
  {[~]}
</ul>
</script>

<script id="l5smgr-node-tmodels-tpl" type="text/html">
<ul class="nav nav-pills nav-stacked">
  {[~it.items :v]}
  <li class="term-item {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
  {[~]}
</ul>
</script>

<script type="text/javascript">

$("#l5smgr-node-specls").on("click", ".spec-item", function() {

    $("#l5smgr-node-specls").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("l5smgr_spec_active", $(this).attr("tgspec"));
    l4iStorage.Del("l5smgr_nmodel_active");
    l4iStorage.Del("l5smgr_tmodel_active");

    l5sNode.Index();
});

$("#l5smgr-node-nmodels").on("click", ".node-item", function() {
    
    $("#l5smgr-node-nmodels").find(".active").removeClass("active");
    $("#l5smgr-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("l5smgr_nmodel_active", $(this).attr("tgname"));

    l5sNode.List(null, $(this).attr("tgname"));
});

$("#l5smgr-node-tmodels").on("click", ".term-item", function() {

    $("#l5smgr-node-nmodels").find(".active").removeClass("active");
    $("#l5smgr-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("l5smgr_tmodel_active", $(this).attr("tgname"));

    l5sTerm.List(null, $(this).attr("tgname"));
});

</script>
