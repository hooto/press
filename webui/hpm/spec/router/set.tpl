<style type="text/css">
#hpm-spec-route-params td {
  padding: 0 10px 10px 0;
}
</style>

<div id="hpm-spec-routeset-alert"></div>
<div id="hpm-spec-routeset"></div>

<script id="hpm-spec-routeset-tpl" type="text/html">
  <input type="hidden" name="modname" value="{[=it._modname]}">

  {[? it.path]}
  <input type="hidden" name="path" value="{[=it.path]}">
  {[??]}
  <div class="form-group">
    <label>Path</label>
    <input type="text" class="form-control" name="path" 
      placeholder="Route Path" value="{[=it.path]}">
  </div>
  {[?]}

  <div class="form-group">
    <label>Data Action</label>
    <select class="form-control" name="data_action">
      {[~it._actions :v]}
      <option value="{[=v.name]}" {[if (v.name == it.dataAction) { ]}selected{[ } ]}>{[=v.name]}</option>
      {[~]}
    </select>
  </div>

  <div class="form-group">
    <label>Template</label>   
    <div class="input-group">
      <input id="hpm-spec-routeset-template" type="text" class="form-control" name="template" placeholder="Template Path" value="{[=it.template]}">
      <span class="input-group-btn">
        <button class="btn btn-default" type="button" onclick="hpSpec.RouteSetTemplateSelect('{[=it._modname]}')">Select a Template</button>
      </span>
    </div>
  </div>

  <div class="form-group">
    <label>Params</label>
    <table id="hpm-spec-route-params" width="100%"></table>
  </div>
  
  <div class="form-group">
    <label>Default</label>
    <select class="form-control" name="default">
      <option value="1" {[if (it.default) { ]}selected{[ } ]}>Yes</option>
      <option value="0" {[if (!it.default) { ]}selected{[ } ]}>No</option>
    </select>
  </div>
</script>

<script id="hpm-spec-route-param-item-tpl" type="text/html">
<tr class="hpm-spec-route-param-item">
  <td>
    <input type="text" class="form-control input-sm" name="param_key" size="16" placeholder="Param Name" value="{[=it._key]}">
  </td>
  <td>
    <input type="text" class="form-control input-sm" name="param_value" size="32" placeholder="Param Value" value="{[=it._value]}">
  </td>
</tr>
</script>
