<div id="htapm-node-navbar">
  <ul id="htapm-node-optools" class="htapm-node-nav"></ul>
  <ul id="htapm-node-tmodels" class="htapm-node-nav htapm-nav-right"></ul>
  <ul id="htapm-node-nmodels" class="htapm-node-nav htapm-nav-right"></ul>
</div>

<div id="htapm-node-workspace">
  <div id="htapm-node-alert"></div>
  <div id="work-content"></div>
</div>

<script id="htapm-node-nmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="node-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script id="htapm-node-tmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="term-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script type="text/javascript">

$("#htapm-node-nmodels").on("click", ".node-item", function() {
    
    $("#htapm-node-nmodels").find(".active").removeClass("active");
    $("#htapm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("htapm_nmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("htapm_nodels_page");
    l4iStorage.Del("htapm_termls_page");

    htapNode.List(null, $(this).attr("tgname"));
});

$("#htapm-node-tmodels").on("click", ".term-item", function() {

    $("#htapm-node-nmodels").find(".active").removeClass("active");
    $("#htapm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("htapm_tmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("htapm_nodels_page");
    l4iStorage.Del("htapm_termls_page");

    htapTerm.List(null, $(this).attr("tgname"));
});

</script>
