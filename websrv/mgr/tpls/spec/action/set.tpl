<style>
.l5smgr-spec-action-datax-attr-item td {
  padding: 0 2px 4px;
}
</style>

<div id="l5smgr-spec-actionset-alert"></div>

<div id="l5smgr-spec-actionset">

  <input type="hidden" name="modname" value="{[=it._modname]}">

  {[? it.name]}
  <input type="hidden" name="name" value="{[=it.name]}">
  {[??]}
  <div class="form-group">
    <label>Name</label>
    <input type="text" class="form-control input-sm" name="name" 
      placeholder="Action Name" value="{[=it.name]}">
  </div>
  {[?]}

  <div class="form-group">
    <label>Datax</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Name</th>
          <th>Query Table</th>
          <th>Pager</th>
          <th>Type</th>
          <th>Limit</th>
          <th>Order</th>
          <th>Cache TTL</th>
        </tr>
      </thead>
      <tbody id="l5smgr-spec-action-dataxs">
        {[~it.datax :v]}
        <tr id="datax-seq-{[=v._seqid]}" class="l5smgr-spec-action-datax-item">
          <td><input type="text" class="form-control input-sm" name="datax_name" size="10" value="{[=v.name]}" readonly></td>
          <td>
            <select class="form-control input-sm" name="datax_query_table">
            {[~it._nodeModels :nmv]}
              <option value="node.{[=nmv.meta.name]}" 
                {[ if (v.type.substr(0,4) == "node" && nmv.meta.name == v.query.table) { ]}selected{[ } ]}>node : {[=nmv.meta.name]}
              </option>
            {[~]}
            {[~it._termModels :tmv]}
              <option value="term.{[=tmv.meta.name]}" 
                {[ if (v.type.substr(0,4) == "term" && tmv.meta.name == v.query.table) { ]}selected{[ } ]}>term : {[=tmv.meta.name]}
              </option>
            {[~]}
            </select>
          </td>
          <td>
            <select class="form-control input-sm" name="datax_pager">
              <option value="true" {[ if (v.pager) { ]}selected{[ } ]}>YES</option>
              <option value="false" {[ if (!v.pager) { ]}selected{[ } ]}>NO</option>
            </select>
          </td>
          <td>
            <select class="form-control input-sm" name="datax_type">
            {[~it._datax_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.type.slice(5)) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
          <td>
            <input type="text" class="form-control input-sm" name="datax_query_limit" size="4" value="{[=v.query.limit]}">
          </td>
          <td>
            <input type="text" class="form-control input-sm" name="datax_query_order" size="8" value="{[=v.query.order]}">
          </td>
          <td>
            <input type="text" class="form-control input-sm" name="datax_cache_ttl" size="4" value="{[=v.cache_ttl]}">
          </td>
        </tr>
        {[~]}
      </tbody>
      </table>
    </div>
  </div>

</div>

<script id="l5smgr-spec-action-datax-item-tpl" type="text/html">
  <tr id="datax-seq-{[=it._seqid]}" class="l5smgr-spec-action-datax-item">
    <td><input type="text" class="form-control input-sm" name="datax_name" size="10" value=""></td>
    <td>
      <select class="form-control input-sm" name="datax_query_table">
      {[~it._nodeModels :nmv]}
        <option value="node.{[=nmv.meta.name]}" >node : {[=nmv.meta.name]}
        </option>
      {[~]}
      {[~it._termModels :tmv]}
        <option value="term.{[=tmv.meta.name]}" >term : {[=tmv.meta.name]}
        </option>
      {[~]}
      </select>
    </td>
    <td>
      <select class="form-control input-sm" name="datax_pager">
        <option value="true">YES</option>
        <option value="false" selected>NO</option>
      </select>
    </td>
    <td>
      <select class="form-control input-sm" name="datax_type">
      {[~it._datax_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == "list") { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
    <td>
      <input type="text" class="form-control input-sm" name="datax_query_limit" size="4" value="1">
    </td>
    <td>
      <input type="text" class="form-control input-sm" name="datax_query_order" size="8" value="">
    </td>
    <td>
      <input type="text" class="form-control input-sm" name="datax_cache_ttl" size="4" value="0">
    </td>
  </tr>
</script>

