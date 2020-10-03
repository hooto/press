
<div id="hpm-specset-alert"></div>


<div id="hpm-specset">
<table class="hpm-formtable">
  <tr>
    <td width="180px">Name</td>
    <td>
    {[? it.meta.name]}
    <input type="text" class="form-control" name="name" value="{[=it.meta.name]}" readonly>
    {[??]}
    <input type="text" class="form-control" name="name" 
      placeholder="Module Name" value="{[=it.meta.name]}">
    {[?]}
    </td>
  </tr>

  <tr>
    <td>Service Path</td>
    <td>
      <input type="text" class="form-control" name="srvname" 
      placeholder="URL Prefix Name of Http Service" value="{[=it.srvname]}">
	  <small>URL Prefix Name of Http Service</small>
    </td>
  </tr>

  <tr>
    <td>Title</td>
    <td>
      <input type="text" class="form-control" name="title" 
        placeholder="Title" value="{[=it.title]}">
    </td>
  </tr>

  {[if (it.meta.name != "core/general") {]}
  <tr>
    <td>Status</td>
    <td>
	<select class="form-control" name="status">
      <option value="1" {[if (it.status) { ]}selected{[ } ]}>Enable</option>
      <option value="0" {[if (!it.status) { ]}selected{[ } ]}>Disable</option>
    </select>
    </td>
  </tr>
  {[}]}

  <tr>
    <td>Theme</td>
    <td>
      <textarea class="form-control" name="theme_config" rows="8">{[? it.theme_config]}{[=it.theme_config]}{[?]}</textarea> 
    </td>
  </tr>

</table>
</div>
