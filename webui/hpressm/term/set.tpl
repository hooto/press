
<div id="hpressm-termset" style="box-sizing: border-box;">loading</div>

<div style="margin:20px 0;">
  <button class="pure-button btapm-btn btapm-btn-primary" onclick="hpressTerm.SetCommit()">Save</button>
  <button class="pure-button btapm-btn" onclick="hpressTerm.List()">Cancel</button>
</div>

<script id="hpressm-termset-tpl" type="text/html">
  <input type="hidden" name="model_type" value="{[=it.model.type]}">
  <input type="hidden" name="id" value="{[=it.id]}">
  <input type="hidden" name="status" value="{[=it.status]}">
  
  <div class="l4i-form-group">
    <label>Title</label>
    <p><input name="title" type="text" value="{[=it.title]}" class="l4i-form-control"></p>
  </div>

  {[ if (it.model.type == "taxonomy") { ]}
  <div class="l4i-form-group">
    <label>Relations</label>
    <p>
      <select name="pid" class="l4i-form-control">
        <option value="0" {[ if (it.pid == 0) { ]}selected{[ } ]}>ROOT</option>
        {[~it._taxonomy_ls.items :v]}
        {[ if (v.pid == 0 && v.id != it.id) { ]}
        <option value="{[=v.id]}" {[ if (it.pid == v.id) { ]}selected{[ } ]}>{[=v.title]}</option>
          {[? v._subs]}
          {[~v._subs :v2]}
          {[ if (v2.id != it.id) { ]}
          <option value="{[=v2.id]}" {[ if (it.pid == v2.id) { ]}selected{[ } ]}>{[=hpressTerm.Sprint(v2._dp)]}{[=v2.title]}</option>
          {[ } ]}
          {[~]}
          {[?]}
        {[ } ]}
        {[~]}
      </select>
    </p>
  </div>
  {[ } ]}

  {[ if (it.model.type == "taxonomy") { ]}
  <div class="l4i-form-group">
    <label>Weight</label>
    <p><input name="weight" type="text" value="{[=it.weight]}" class="l4i-form-control"></p>
  </div>
  {[ } ]}
</script>
