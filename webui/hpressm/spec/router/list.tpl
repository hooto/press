
<div id="hpressm-spec-routels-alert"></div>

<table class="table table-hover">
  <thead>
    <tr>
      <th>Path</th>
      <th>Data Action</th>
      <th>Template</th>
      <th>Params</th>
      <th>Default</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="hpressm-spec-routels"></tbody>
</table>

<script id="hpressm-spec-routels-tpl" type="text/html">  
{[~it.router.routes :v]}
<tr>
  <td class="hpressm-font-fixspace">{[=v.path]}</td>
  <td>{[=v.dataAction]}</td>
  <td>{[=v.template]}</td>
  <td>{[=v._paramsNum]}</td>
  <td>
    {[if (v.default) {]}Yes{[} else {]}No{[}]}
  </td>
  <td align="right">
    <button class="btn btn-default" onclick="hpressSpec.RouteSet('{[=it._modname]}', '{[=v.path]}')">Setting</button>
  </td>
</tr>
{[~]}
</script>
