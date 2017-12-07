<style>
.hpressm-spec-node-field-attr-item td {
  padding: 0 2px 4px;
}
</style>

<div id="hpressm-spec-nodeset-alert"></div>

<div id="hpressm-spec-nodeset">

  <input type="hidden" name="modname" value="{[=it._modname]}">

  {[? it.meta.name]}
  <input type="hidden" name="name" value="{[=it.meta.name]}">
  {[??]}
  <div class="form-group">
    <label>Name</label>
    <input type="text" class="form-control input-sm" name="name" 
      placeholder="Node Name" value="{[=it.meta.name]}">
  </div>
  {[?]}

  <div class="form-group">
    <label>Title</label>
    <input type="text" class="form-control input-sm" name="title" 
      placeholder="Title" value="{[=it.title]}">
  </div>

  <div class="form-group">
    <label>Fields</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Name</th>
          <th>Title</th>
          <th>Type</th>
          <th>Length</th>
          <th>Index Type</th>
          <th>Extended attributes</th>
          <th></th>
        </tr>
      </thead>
      <tbody id="hpressm-spec-node-fields">
        {[~it.fields :v]}
        <tr id="field-seq-{[=v._seqid]}" class="hpressm-spec-node-field-item">
          <td><input type="text" class="form-control input-sm" name="field_name" size="10" value="{[=v.name]}" readonly></td>
          <td><input type="text" class="form-control input-sm" name="field_title" size="20" value="{[=v.title]}"></td>
          <td>
            <select class="form-control input-sm" name="field_type">
            {[~it._field_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.type) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
          <td><input class="form-control input-sm" type="text" name="field_length" size="3" value="{[=v.length]}"></td>
          <td>
            <select class="form-control input-sm" name="field_index_type">
            {[~it._field_idx_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.indexType) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
          <td>
            <table><tbody class="hpressm-spec-node-field-attrs">
              {[~v.attrs :atv]}
              <tr class="hpressm-spec-node-field-attr-item">
                <td><input type="text" class="form-control input-sm" name="field_attr_key" size="8" value="{[=atv.key]}"></td>
                <td><input type="text" class="form-control input-sm" name="field_attr_value" size="16" value="{[=atv.value]}"></td>
              </tr>
              {[~]}
            </tbody></table>
          </td>
          <td>
            <button class="btn btn-default btn-sm" onclick="hpressSpec.NodeSetFieldAttrAppend('{[=v._seqid]}')">+ Attribute</button>
          </td>
        </tr>
        {[~]}
      </tbody>
      </table>
    </div>
  </div>

  <div class="form-group">
    <label>Terms</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Name</th>
          <th>Title</th>
          <th>Type</th>
        </tr>
      </thead>
      <tbody id="hpressm-spec-node-terms">
        {[~it.terms :v]}
        <tr id="field-seq-{[=v._seqid]}" class="hpressm-spec-node-term-item">
          <td><input type="text" class="form-control input-sm" name="term_name" size="20" value="{[=v.meta.name]}" readonly></td>
          <td><input type="text" class="form-control input-sm" name="term_title" size="30" value="{[=v.title]}"></td>
          <td>
            <select class="form-control input-sm" name="term_type">
            {[~it._term_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.type) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
        {[~]}
      </table>
    </div>
  </div>

  <div class="form-group">
    <label>Extensions</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Option</th>
          <th>Attributes</th>
        </tr>
      </thead>
      <tbody id="hpressm-spec-node-exts">
        <tr>
          <td>Access Counter</td>
          <td>
            <select class="form-control input-sm" name="ext_access_counter">
            {[~it._general_onoff :gv]}
            <option value="{[=gv.type]}" {[ if (it.extensions.access_counter == gv.type) { ]}selected{[ } ]}>{[=gv.name]}</option>
            {[~]}
            </select>
          </td>
        </tr>
        <tr>
          <td>Comment Enable</td>
          <td>
            <select class="form-control input-sm" name="ext_comment_enable">
            {[~it._general_onoff :gv]}
            <option value="{[=gv.type]}" {[ if (it.extensions.comment_enable == gv.type) { ]}selected{[ } ]}>{[=gv.name]}</option>
            {[~]}
            </select>
          </td>
        </tr>
        <tr>
          <td>Comment On/Off Per Entry</td>
          <td>
            <select class="form-control input-sm" name="ext_comment_perentry">
            {[~it._general_onoff :gv]}
            <option value="{[=gv.type]}" {[ if (it.extensions.comment_perentry == gv.type) { ]}selected{[ } ]}>{[=gv.name]}</option>
            {[~]}
            </select>
          </td>
        </tr>
        <tr>
          <td>Permalink Settings</td>
          <td>
            <select class="form-control input-sm" name="ext_permalink">
            {[~it._permalink_def :gv]}
            <option value="{[=gv.type]}" {[ if (it.extensions.permalink == gv.type) { ]}selected{[ } ]}>{[=gv.name]}</option>
            {[~]}
            </select>
          </td>
        </tr>
        {[if (it.extensions.node_sub_refer) {]}
        <tr>
          <td>Node Sub Refer</td>
          <td>
		    {[=it.extensions.node_sub_refer]}
          </td>
        </tr>
        {[} else {]}
        <tr>
          <td>Refer to Node Name</td>
          <td>
		    <input type="text" class="form-control input-sm" name="ext_node_refer" value="{[=it.extensions.node_refer]}">
          </td>
        </tr>
        {[}]}
      </tbody>
      </table>
    </div>
  </div>

</div>

<script id="hpressm-spec-node-field-item-tpl" type="text/html">
  <tr id="field-seq-{[=it._seqid]}" class="hpressm-spec-node-field-item">
    <td><input type="text" class="form-control input-sm" name="field_name" size="10" value=""></td>
    <td><input type="text" class="form-control input-sm" name="field_title" size="16" value=""></td>
    <td>
      <select class="form-control input-sm" name="field_type">
      {[~it._field_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._type) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
    <td><input type="text" class="form-control input-sm" name="field_length" size="5" value="0"></td>
    <td>
      <select class="form-control input-sm" name="field_index_type">
      {[~it._field_idx_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._indexType) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
    <td>
      <table><tbody class="hpressm-spec-node-field-attrs"></tbody></table>
    </td>
    <td>
      <button class="btn btn-default btn-sm" onclick="hpressSpec.NodeSetFieldAttrAppend('{[=it._seqid]}')">+ Attribute</button>
    </td>
  </tr>
</script>

<script id="hpressm-spec-node-field-attr-item-tpl" type="text/html">
  <tr class="hpressm-spec-node-field-attr-item">
    <td><input type="text" class="form-control input-sm" name="field_attr_key" size="8" value=""></td>
    <td><input type="text" class="form-control input-sm" name="field_attr_value" size="12" value=""></td>
  </tr>
</script>

<script id="hpressm-spec-node-term-item-tpl" type="text/html">
  <tr id="field-seq-{[=it._seqid]}" class="hpressm-spec-node-term-item">
    <td><input type="text" class="form-control input-sm" name="term_name" size="20" value=""></td>
    <td><input type="text" class="form-control input-sm" name="term_title" size="30" value=""></td>
    <td>
      <select class="form-control input-sm" name="term_type">
      {[~it._term_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._type) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
  </tr>
</script>
