<style type="text/css">
#l5s-s2-nav li a {
  padding: 5px 10px;
  margin-bottom: 10px;
}
</style>

<div id="work-content" class="">
<div id="l5smgr-s2-objls-dirnav"></div>
<table class="table table-hover">
  <thead>
    <tr>
      <th width="80px"></th>
      <th>Name</th>
      <th>Dir</th>
      <th>Size</th>
      <th></th>
      <th></th>
    </tr>
  </thead>
  <tbody id="l5smgr-s2-objls"></tbody>
</table>
<div>
  <button type="button" 
    class="btn btn-primary btn-sm" 
    onclick="l5sS2.ObjNew('file')">
    Upload New File
  </button>
</div>
</div>

<script id="l5smgr-s2-objls-dirnav-tpl" type="text/html">
<ol class="breadcrumb" style="background-color:#fff;">
  {[~it.items :v]}
  <li><a href="#{[=v.path]}" onclick="l5sS2.ObjList('{[=v.path]}')">{[=v.name]}</a></li>
  {[~]}
</ol>
</script>

<script id="l5smgr-s2-objls-tpl" type="text/html">  
  {[~it.items :v]}
    <tr id="obj{[=v._id]}">
      <td>
      {[ if (v.isdir) { ]}
        <span class="glyphicon glyphicon-folder-open" aria-hidden="true"></span>
      {[ } else if (v._isimg) { ]}
        <a href="{[=v.self_link]}" target="_blank"><img src="{[=v.self_link]}?ipn=i6040"></a>
      {[ } ]}
      </td>
      <td class="ts3-fontmono">
      {[ if (v.isdir) { ]}
        <a class="obj-item-dir" href="#objs" path="{[=v._abspath]}">{[=v.name]}</a>
      {[ } else { ]}
        <a class="obj-item-file" href="{[=v.self_link]}" target="_blank">{[=v.name]}</a>
      {[ } ]}
      </td>
      <td>
      {[?v.isdir]}YES{[?]}
      </td>
      <td>
      {[?!v.isdir]}
        {[=l5sS2.UtilResourceSizeFormat(v.size)]}</td>
      {[?]}
      <td align="right">{[=l4i.TimeParseFormat(v.modtime, "Y-m-d H:i:s")]}</td>
      <td align="right">
      {[ if (!v.isdir) { ]}
        <a class="obj-item-del btn btn-default btn-xs" href="#obj-del" obj="{[=v._abspath]}">
          <span class="glyphicon glyphicon-cog" aria-hidden="true"></span> Delete
        </a>
      {[ } ]}
      </td>
    </tr>
  {[~]}
</script>

<!-- TPL : File New -->
<script id="l5smgr-s2-objnew-tpl" type="text/html"> 
<form id="{[=it.formid]}" action="#" onsubmit="l5sS2.ObjNewSave('{[=it.formid]}');return false;">
<input type="hidden" name="type" value="{[=it.type]}">
<div class="form-group">
  <label>Folder Path</label>
  <input type="text" name="path" class="form-control" placeholder="Folder Path" value="{[=it.path]}">
</div>
<div class="form-group">
  <label>Select File</label>
  <input id="l5smgr-s2-objnew-files" type="file" name="file" class="form-control" placeholder="File Path" value="">
</div>
</form>
<div id="{[=it.formid]}-alert" class="alert alert-success" style="display:none"></div>
</script>

<!-- TPL : File Rename -->
<script id="l5smgr-s2-objrename-tpl" type="text/html"> 
<form id="{[=it.formid]}" action="#" onsubmit="l5sS2.ObjRenameSave('{[=it.formid]}');return false;">
  <div class="input-prepend" style="margin-left:2px">
    <span class="add-on">
        <img src="{[=l5sMgr.base]}-/img/folder_edit.png" class="h5c_icon">
    </span>
    <input type="text" name="pathset" value="{[=it.path]}" style="width:500px;">
    <input type="hidden" name="path" value="{[=it.path]}">
  </div>
</form>
</script>

<script type="text/javascript">
$("#l5smgr-s2-objls").on("click", ".obj-item-dir", function() {
    l5sS2.ObjList($(this).attr("path"));
});
$("#l5smgr-s2-objls").on("click", ".obj-item-del", function() {
    var r = confirm("This file will be deleted, Confirm?");
    if (r == true) {
      l5sS2.ObjDel($(this).attr("obj"));
    }
});
$("#l5sS2-object-dirnav").on("click", ".obj-item-dir", function() {
    l5sS2.ObjList($(this).attr("path"));
});
</script>