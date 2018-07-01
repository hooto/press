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
<div id="hpm-nodeset-top-title"></div>
<div id="hpm-nodeset-tops"></div>
<div id="hpm-nodeset-fields"></div>
</script>

<script id="hpm-nodeset-tplstatus" type="text/html">
<div class="hpm-nodeset-tplx">
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
<div class="hpm-nodeset-tplx hpm-nodeset-tpltext">

  <label>
    <span>{[=it.title]}</span>

    {[? it.attr_lang_list]}
    <select id="field_{[=it.name]}_langs" class="field-nav-lang form-control" onchange="hpNode.SetFieldLang('{[=it.name]}')">
      {[~it.attr_lang_list :v]}
  	<option value="{[=v.id]}">{[=v.name]}</option>
      {[~]}
    </select>
    {[?]}
  </label>

  <input type="hidden" id="field_{[=it.name]}_attr_format" name="field_{[=it.name]}_attr_format" value="{[=it.attr_format]}">

  <div class="editor-outbox">  
    <div id="field_{[=it.name]}_inner_toolbar" class="editor-inner-toolbar">

      <span id="field_{[=it.name]}_editor_nav" class="editor-nav pure-button-group">
        {[~it._formats :v]}
        <button class="tpltext-editor-item editor-nav-{[=v.name]} pure-button button-xsmall" 
          onclick="hpEditor.Open('{[=it.name]}', '{[=v.name]}')">{[=v.value]}</button>
        {[~]}
      </span>

	  <span class="vline"></span>

      <span id="field_{[=it.name]}_editor_mdr" class="" style="display:none">
        <button class="pure-button button-xsmall preview_open" onclick="hpEditor.PreviewOpen('{[=it.name]}')" style="display:none">Preview</button>
        <button class="pure-button button-xsmall preview_close" onclick="hpEditor.PreviewClose('{[=it.name]}')" style="display:none">Close Preview</button>
        <button class="pure-button button-xsmall storage-image-insert" onclick="hpEditor.StorageImageSelector('{[=it.name]}')">Image</button>
      </span>

	</div>

    <div id="field_{[=it.name]}_layout" class="editor-fra">
      <div id="field_{[=it.name]}_editor" class="editor-fra-item">
        <textarea class="l4i-form-control" id="field_{[=it.name]}" name="field_{[=it.name]}" rows="{[if (it.attr_ui_rows) {]}{[=it.attr_ui_rows]}{[} else {]}6{[}]}">{[=it.value]}</textarea>
      </div>
      <div id="field_{[=it.name]}_colpreview" style="display:none" classs="editor-fra-item">
        <div class="markdown-body hp-scroll" 
          id="field_{[=it.name]}_preview" style="padding:10px"></div>
      </div>
    </div>
 
  </div>
</div>
</script>

<script id="hpm-nodeset-tplint" type="text/html">
  <div class="hpm-nodeset-tplx">
    <label>{[=it.title]}</label>
    <input type="text" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplstring" type="text/html">
  <div class="hpm-nodeset-tplx hpm-nodeset-tplstring">
    <label>
      <span>{[=it.title]}</span>
      {[? it.attr_lang_list]}
      <select id="field_{[=it.name]}_langs" class="field-nav-lang form-control" onchange="hpNode.SetFieldLang('{[=it.name]}')">
        {[~it.attr_lang_list :v]}
		<option value="{[=v.id]}">{[=v.name]}</option>
        {[~]}
	  </select>
      {[?]}
    </label>
    <input type="text" id="field_{[=it.name]}" name="field_{[=it.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplterm_tag" type="text/html">
  <div class="hpm-nodeset-tplx">
    <label>{[=it.title]}</label>
    <input type="text" name="term_{[=it.meta.name]}" class="l4i-form-control" value="{[=it.value]}">
  </div>
</script>

<script id="hpm-nodeset-tplterm_taxonomy" type="text/html">
  <div class="hpm-nodeset-tplx">
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
  <div class="hpm-nodeset-tplx">
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
  <div class="hpm-nodeset-tplx">
    <label>Permalink Name</label>
    <div>
      <input type="text" name="ext_permalink_name" class="l4i-form-control" value="{[=it.ext_permalink_name]}">
    </div>
  </div>
</script>


<script id="hpm-nodeset-tplext_node_refer" type="text/html">
  <div class="hpm-nodeset-tplx">
    <label>Refer ID</label>
    <div>
      <input type="text" name="ext_node_refer" class="l4i-form-control" value="{[=it.ext_node_refer]}">
    </div>
  </div>
</script>

