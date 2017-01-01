
<div id="htpm-specset-alert"></div>


<div id="htpm-specset">
  
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
</div>
