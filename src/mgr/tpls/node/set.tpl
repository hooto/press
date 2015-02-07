
<div id="l5smgr-nodeset" style="box-sizing: border-box;">loading</div>

<div>
  <button class="btn btn-default btn-primary" onclick="l5sNode.SetCommit()">Save</button>
  <button class="btn btn-default" onclick="l5sNode.List()">Cancel</button>
</div>

<script id="l5smgr-nodeset-tpl" type="text/html">
  <input type="hidden" name="id" value="{[=it.id]}">
  <input type="hidden" name="state" value="{[=it.state]}">
  <div class="l4i-form-group">
    <label>Title</label>
    <p><input name="title" type="text" value="{[=it.title]}" class="l4i-form-control"></p>
  </div>
</script>

<script id="l5smgr-nodeset-tpltext" type="text/html">
  <div class="l4i-form-group">
    <label id="field_{[=it.name]}_tools">
      {[=it.title]}
      <a class="btn btn-default btn-xs preview_open" href="#preview-open" onclick="l5sEditor.PreviewOpen('{[=it.name]}')" style="display:none">Open real-time preview</a>
      <a class="btn btn-default btn-xs preview_close" href="#preview-close" onclick="l5sEditor.PreviewClose('{[=it.name]}')" style="display:none">Close real-time preview</a>
    </label>
    <input type="hidden" name="field_{[=it.name]}_attr_format" value="{[=it.attr_format]}">
    
    <table width="100%" id="field_{[=it.name]}_layout">
    <tr>
      <td id="field_{[=it.name]}_editor" valign="top">
        <textarea class="l4i-form-control " id="field_{[=it.name]}" name="field_{[=it.name]}" rows="{[if (it.attr_ui_rows) {]}{[=it.attr_ui_rows]}{[} else {]}6{[}]}">{[=it.value]}</textarea>
      </td>
      <td id="field_{[=it.name]}_colspace" style="display:none"></td>
      <td id="field_{[=it.name]}_colpreview" valign="top" style="display:none" classs="l5s-scroll">
        <div class="markdown-body l5s-scroll" 
          id="field_{[=it.name]}_preview" style="float:right;padding:5px 5px 5px 0;"></div>
      </td>
    </tr>
    </table>
  </div>
</script>

<script id="l5smgr-nodeset-tplint" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="l5smgr-nodeset-tplstring" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="l5smgr-nodeset-tplterm_tag" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="term_{[=it.metadata.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="l5smgr-nodeset-tplterm_taxonomy" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.item.title]}</label>
    <div>
    <select class="form-control" name="term_{[=it.item.metadata.name]}">
      {[~it.items :v]}
      <option value="{[=v.id]}" {[if (it.item.value == v.id) { ]}selected{[ } ]}>{[=v.title]}</option>
      {[~]}
    </select>
    </div>
  </div>
</script>
