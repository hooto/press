<div id="hpm-nodeset-layout" style="display: flex; align-items: flex-start;">
  <div id="hpm-nodeset-laymain" style="flex:3">loading</div>
  <div id="hpm-nodeset-layside" class="" style="flex:1;"></div>
</div>

<div style="margin:20px 0;">
  <button class="pure-button btapm-btn btapm-btn-primary" onclick="hpNode.SetCommit()">Save</button>
  <button class="pure-button btapm-btn" onclick="hpNode.List()">Cancel</button>
</div>

<div id="hpm-node-set-opts" class="hpm-hide">
  <li id="hpm-node-set-opts-label" style="font-weight:bold;">Content</li>
</div>

<script id="hpm-nodeset-tpl" type="text/html">
  <input type="hidden" name="id" value="{[=it.id]}">
  <div class="l4i-form-group">
    <label>Title</label>
    <p><input name="title" type="text" value="{[=it.title]}" class="l4i-form-control"></p>
  </div>
  <div id="hpm-nodeset-tops"></div>
  <div id="hpm-nodeset-fields"></div>
</script>

<script id="hpm-nodeset-tplstatus" type="text/html">
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

<script id="hpm-nodeset-tpltext" type="text/html">
  <div class="l4i-form-group hpm-nodeset-tpltext">
    <label id="field_{[=it.name]}_tools">
      <span>{[=it.title]}</span>
      <span id="field_{[=it.name]}_editor_nav" class="editor-nav">
        {[~it._formats :v]}
        <a class="tpltext-editor-item editor-nav-{[=v.name]}" href="#" 
          onclick="hpEditor.Open('{[=it.name]}', '{[=v.name]}')">{[=v.value]}</a>
        {[~]}
      </span>
      <span id="field_{[=it.name]}_editor_mdr" class="editor_mdr" style="display:none">
        <a class="tpltext-editor-item preview_open" href="#preview-open" onclick="hpEditor.PreviewOpen('{[=it.name]}')" style="display:none">Open Markdown Preview</a>
        <a class="tpltext-editor-item preview_close" href="#preview-close" onclick="hpEditor.PreviewClose('{[=it.name]}')" style="display:none">Close Markdown Preview</a>
      </span>
    </label>
    <input type="hidden" id="field_{[=it.name]}_attr_format" name="field_{[=it.name]}_attr_format" value="{[=it.attr_format]}">
    
    <table width="100%" id="field_{[=it.name]}_layout" class="editor-fra">
    <tr>
      <td id="field_{[=it.name]}_editor" valign="top">
        <textarea class="l4i-form-control" id="field_{[=it.name]}" name="field_{[=it.name]}" rows="{[if (it.attr_ui_rows) {]}{[=it.attr_ui_rows]}{[} else {]}6{[}]}">{[=it.value]}</textarea>
      </td>
      <td id="field_{[=it.name]}_colspace" style="display:none"></td>
      <td id="field_{[=it.name]}_colpreview" valign="top" style="display:none" classs="hp-scroll">
        <div class="markdown-body hp-scroll" 
          id="field_{[=it.name]}_preview" style="float:right;padding:5px"></div>
      </td>
    </tr>
    </table>
  </div>
</script>

<script id="hpm-nodeset-tplint" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplstring" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplterm_tag" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.title]}</label>
    <input type="text" name="term_{[=it.meta.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplterm_taxonomy" type="text/html">
  <div class="l4i-form-group">
    <label>{[=it.model.title]}</label>
    <div>
    <select class="form-control" name="term_{[=it.model.meta.name]}">
    {[~it.items :v]}
      {[ if (v.pid == 0) { ]}
      <option value="{[=v.id]}" {[if (it.item.value == v.id) { ]}selected{[ } ]}>{[=v.title]}</option>
      {[? v._subs]}
        {[~v._subs :v2]}
        <option value="{[=v2.id]}" {[if (it.item.value == v2.id) { ]}selected{[ } ]}>{[=hpTerm.Sprint(v2._dp)]}{[=v2.title]}</option>
        {[~]}
      {[}]}
      {[ } ]}
    {[~]}
    </select>
    </div>
  </div>
</script>


<script id="hpm-nodeset-tplext_comment_perentry" type="text/html">
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


<script id="hpm-nodeset-tplext_permalink" type="text/html">
  <div class="l4i-form-group">
    <label>Permalink Name</label>
    <div>
      <input type="text" name="ext_permalink_name" class="l4i-form-control" value="{[=it.ext_permalink_name]}">
    </div>
  </div>
</script>


<script id="hpm-nodeset-tplext_node_refer" type="text/html">
  <div class="l4i-form-group">
    <label>Refer ID</label>
    <div>
      <input type="text" name="ext_node_refer" class="l4i-form-control" value="{[=it.ext_node_refer]}">
    </div>
  </div>
</script>

