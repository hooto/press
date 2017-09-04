<div id="hpressm-node-navbar">
  <ul id="hpressm-node-optools" class="hpressm-node-nav"></ul>
  <ul id="hpressm-node-tmodels" class="hpressm-node-nav hpressm-nav-right"></ul>
  <ul id="hpressm-node-nmodels" class="hpressm-node-nav hpressm-nav-right"></ul>
</div>

<div id="hpressm-node-workspace">
  <div id="hpressm-node-alert"></div>
  <div id="work-content"></div>
</div>

<script id="hpressm-node-nmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="node-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script id="hpressm-node-tmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="term-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script type="text/javascript">

$("#hpressm-node-nmodels").on("click", ".node-item", function() {
    
    $("#hpressm-node-nmodels").find(".active").removeClass("active");
    $("#hpressm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("hpressm_nmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("hpressm_nodels_page");
    l4iStorage.Del("hpressm_termls_page");

    hpressNode.List(null, $(this).attr("tgname"));
});

$("#hpressm-node-tmodels").on("click", ".term-item", function() {

    $("#hpressm-node-nmodels").find(".active").removeClass("active");
    $("#hpressm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("hpressm_tmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("hpressm_nodels_page");
    l4iStorage.Del("hpressm_termls_page");

    hpressTerm.List(null, $(this).attr("tgname"));
});

</script>
