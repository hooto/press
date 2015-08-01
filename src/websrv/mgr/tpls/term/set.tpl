
<div id="l5smgr-termset" style="box-sizing: border-box;">loading</div>

<div>
  <button class="btn btn-default btn-primary" onclick="l5sTerm.SetCommit()">Save</button>
  <button class="btn btn-default" onclick="l5sTerm.List()">Cancel</button>
</div>

<script id="l5smgr-termset-tpl" type="text/html">
  <input type="hidden" name="model_type" value="{[=it.model.type]}">
  <input type="hidden" name="id" value="{[=it.id]}">
  <input type="hidden" name="pid" value="{[=it.pid]}">
  <input type="hidden" name="status" value="{[=it.status]}">
  <div class="l4i-form-group">
    <label>Title</label>
    <p><input name="title" type="text" value="{[=it.title]}" class="l4i-form-control"></p>
  </div>
  {[ if (it.model.type == "taxonomy") { ]}
  <div class="l4i-form-group">
    <label>Weight</label>
    <p><input name="weight" type="text" value="{[=it.weight]}" class="l4i-form-control"></p>
  </div>
  {[ } ]}
</script>
