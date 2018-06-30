<div id="hpm-node-navbar">
  <ul id="hpm-node-optools" class="hpm-node-nav"></ul>
  <ul id="hpm-node-tmodels" class="hpm-node-nav hpm-nav-right"></ul>
  <ul id="hpm-node-nmodels" class="hpm-node-nav hpm-nav-right"></ul>
</div>

<div id="hpm-node-workspace">
  <div id="hpm-node-alert"></div>
  <div id="work-content"></div>
</div>

<script id="hpm-node-nmodels-tpl" type="text/html">
{[~it.items :v]}
{[if (!v.extensions.node_refer) {]}
  <li class="node-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[}]}
{[~]}
</script>

<script id="hpm-node-tmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="term-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script type="text/javascript">

$("#hpm-node-nmodels").on("click", ".node-item", function() {
    
    $("#hpm-node-nmodels").find(".active").removeClass("active");
    $("#hpm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

	hpNode.SpecNodeModelActive($(this).attr("tgname"));
    l4iStorage.Del("hpm_nodels_page");
    l4iStorage.Del("hpm_termls_page");

    hpNode.List(null, $(this).attr("tgname"));
});

$("#hpm-node-tmodels").on("click", ".term-item", function() {

    $("#hpm-node-nmodels").find(".active").removeClass("active");
    $("#hpm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

	hpTerm.SpecTermModelActive($(this).attr("tgname"));
    l4iStorage.Del("hpm_nodels_page");
    l4iStorage.Del("hpm_termls_page");

    hpTerm.List(null, $(this).attr("tgname"));
});

</script>
