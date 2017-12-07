
<div id="hpressm-specset-alert"></div>


<div id="hpressm-specset">
  
  {[? it.meta.name]}
  <input type="hidden" name="name" value="{[=it.meta.name]}">
  {[??]}
  <div class="form-group">
    <label>Name</label>
    <input type="text" class="form-control" name="name" 
      placeholder="Module Name" value="{[=it.meta.name]}">
  </div>
  {[?]}

  <div class="form-group">
    <label>URL Prefix Name of Http Service</label>
    <input type="text" class="form-control" name="srvname" 
      placeholder="URL Prefix Name of Http Service" value="{[=it.srvname]}">
  </div>

  <div class="form-group">
    <label>Title</label>
    <input type="text" class="form-control" name="title" 
      placeholder="Title" value="{[=it.title]}">
  </div>

  {[if (it.meta.name != "core/general") {]}
  <div class="form-group">
    <label>Status</label>
    <select class="form-control" name="status">
      <option value="1" {[if (it.status) { ]}selected{[ } ]}>Enable</option>
      <option value="0" {[if (!it.status) { ]}selected{[ } ]}>Disable</option>
    </select>
  </div>
  {[}]}
</div>
