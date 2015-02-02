
<div id="l5smgr-node-specls" style="padding-bottom:10px;border-bottom:1px solid #ccc"></div>

<script id="l5smgr-node-specls-tpl" type="text/html">
<ul class="nav nav-pills">
  {[~it.items :v]}
  <li class="{[if (it.active == v.metadata.id) {]}active{[}]}"><a href="#{[=v.metadata.id]}">{[=v.title]}</a></li>
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
  <li class="{[if (it.active == v.metadata.name) {]}active{[}]}"><a href="#{[=v.metadata.id]}">{[=v.title]}</a></li>
  {[~]}
</ul>
</script>

<script id="l5smgr-node-tmodels-tpl" type="text/html">
<ul class="nav nav-pills nav-stacked">
  {[~it.items :v]}
  <li class="{[if (it.active == v.metadata.name) {]}active{[}]}"><a href="#{[=v.metadata.id]}">{[=v.title]}</a></li>
  {[~]}
</ul>
</script>