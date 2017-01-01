<div id="htpm-node-navbar">
  <ul id="htpm-node-optools" class="htpm-node-nav"></ul>
  <ul id="htpm-node-tmodels" class="htpm-node-nav htpm-nav-right"></ul>
  <ul id="htpm-node-nmodels" class="htpm-node-nav htpm-nav-right"></ul>
</div>

<div id="htpm-node-workspace">
  <div id="htpm-node-alert"></div>
  <div id="work-content"></div>
</div>

<script id="htpm-node-nmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="node-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script id="htpm-node-tmodels-tpl" type="text/html">
{[~it.items :v]}
  <li class="term-item pure-button {[if (it.active == v.meta.name) {]}active{[}]}" tgname="{[=v.meta.name]}">
    <a href="#{[=v.meta.name]}">{[=v.title]}</a>
  </li>
{[~]}
</script>

<script type="text/javascript">

$("#htpm-node-nmodels").on("click", ".node-item", function() {
    
    $("#htpm-node-nmodels").find(".active").removeClass("active");
    $("#htpm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("htpm_nmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("htpm_nodels_page");
    l4iStorage.Del("htpm_termls_page");

    htpNode.List(null, $(this).attr("tgname"));
});

$("#htpm-node-tmodels").on("click", ".term-item", function() {

    $("#htpm-node-nmodels").find(".active").removeClass("active");
    $("#htpm-node-tmodels").find(".active").removeClass("active");
    $(this).addClass("active");

    l4iStorage.Set("htpm_tmodel_active", $(this).attr("tgname"));
    l4iStorage.Del("htpm_nodels_page");
    l4iStorage.Del("htpm_termls_page");

    htpTerm.List(null, $(this).attr("tgname"));
});

</script>
