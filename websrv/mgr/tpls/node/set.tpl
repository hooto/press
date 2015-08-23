
<div id="l5smgr-nodeset" style="box-sizing: border-box;">loading</div>

<div style="margin-bottom:20px;">
  <button class="btn btn-default btn-primary" onclick="l5sNode.SetCommit()">Save</button>
  <button class="btn btn-default" onclick="l5sNode.List()">Cancel</button>
</div>

<script id="l5smgr-nodeset-tpl" type="text/html">
  <input type="hidden" name="id" value="{[=it.id]}">
  <div class="l4i-form-group">
    <label>Title</label>
    <p><input name="title" type="text" value="{[=it.title]}" class="l4i-form-control"></p>
  </div>
  <div id="l5smgr-nodeset-fields"></div>
  <div class="l4i-form-group">
    <label>Status</label>
    <p>
      <select name="status" class="l4i-form-control">
      {[~it._status_def :sv]}
        <option value="{[=sv.type]}" {[if (sv.type == it.status) { ]}selected{[ } ]}>{[=sv.name]}</option>
      {[~]}
      </select>
    </p>
  </div>
</script>

<script id="l5smgr-nodeset-tpltext" type="text/html">
  <div class="l4i-form-group l5smgr-nodeset-tpltext">
    <label id="field_{[=it.name]}_tools">
      <span>{[=it.title]}</span>
      <span id="field_{[=it.name]}_editor_nav" class="editor-nav">
        <a class="tpltext-editor-item editor-nav-text" href="#" 
          onclick="l5sEditor.Open('{[=it.name]}', 'text')">Text</a>
        <a class="tpltext-editor-item editor-nav-md" href="#"
          onclick="l5sEditor.Open('{[=it.name]}', 'md')">Markdown</a>
      </span>
      <span id="field_{[=it.name]}_editor_mdr" class="editor_mdr" style="display:none">
        <a class="tpltext-editor-item preview_open" href="#preview-open" onclick="l5sEditor.PreviewOpen('{[=it.name]}')" style="display:none">Open Markdown Preview</a>
        <a class="tpltext-editor-item preview_close" href="#preview-close" onclick="l5sEditor.PreviewClose('{[=it.name]}')" style="display:none">Close Markdown Preview</a>
      </span>
    </label>
    <input type="hidden" id="field_{[=it.name]}_attr_format" name="field_{[=it.name]}_attr_format" value="{[=it.attr_format]}">
    
    <table width="100%" id="field_{[=it.name]}_layout" class="editor-fra">
    <tr>
      <td id="field_{[=it.name]}_editor" valign="top">
        <textarea class="l4i-form-control " id="field_{[=it.name]}" name="field_{[=it.name]}" rows="{[if (it.attr_ui_rows) {]}{[=it.attr_ui_rows]}{[} else {]}6{[}]}">{[=it.value]}</textarea>
      </td>
      <td id="field_{[=it.name]}_colspace" style="display:none"></td>
      <td id="field_{[=it.name]}_colpreview" valign="top" style="display:none" classs="l5s-scroll">
        <div class="markdown-body l5s-scroll" 
          id="field_{[=it.name]}_preview" style="float:right;padding:5px"></div>
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
    <input type="text" name="term_{[=it.meta.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="l5smgr-nodeset-tplterm_taxonomy" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.item.title]}</label>
    <div>
    <select class="form-control" name="term_{[=it.item.meta.name]}">
      {[~it.items :v]}
      <option value="{[=v.id]}" {[if (it.item.value == v.id) { ]}selected{[ } ]}>{[=v.title]}</option>
      {[~]}
    </select>
    </div>
  </div>
</script>


<script id="l5smgr-nodeset-tplext_comment_perentry" type="text/html">
  <div class="l4i-form-group">
    <label>Comment On/Off</label>
    <div>
    <select class="form-control" name="ext_comment_perentry">
      {[~it._general_onoff :gv]}
      <option value="{[=gv.type]}" {[ if (it.ext_comment_perentry == gv.type) { ]}selected{[ } ]}>{[=gv.name]}</option>
      {[~]}
    </select>
    </div>
  </div>
</script>
