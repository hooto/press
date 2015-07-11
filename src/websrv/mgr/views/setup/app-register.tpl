<div class="container">

<div class="panel panel-default">
  <div class="panel-heading">Registration lessfly to lessids</div>
  <div class="panel-body">

<form id="n51hwa" action="#">
  
  <div id="cdnxe5" class="alert alert-danger hide">...</div>

  <div class="form-group">
    <label>lessids service url</label>
    <input type="text" name="lessids_url" class="form-control" placeholder="Enter the lessIds URL" value="{{.lessids_url}}">
  </div>

  <div class="form-group">
    <label>Instance ID</label>
    <input type="text" name="instance_id" class="form-control" value="{{.instance_id}}" readonly>
  </div>

  <div class="form-group">
    <label>Instance URL</label>
    <input type="text" name="instance_url" class="form-control" value="{{.instance_url}}">
  </div>

  <div class="form-group">
    <label>Application ID</label>
    <input type="text" name="app_id" class="form-control" value="{{.app_id}}" readonly>
  </div>

  <div class="form-group">
    <label>Application Name</label>
    <input type="text" name="app_title" class="form-control" placeholder="Enter the name of application" value="{{.app_title}}">
  </div>

  <div class="form-group">
    <label>Application Version</label>
    <input type="text" name="version" class="form-control" value="{{.version}}" readonly>
  </div>

  <button type="submit" class="btn btn-primary">Commit</button>
</form>

  </div>
</div>

</div>

<script type="text/javascript">


$("#n51hwa").submit(function(event) {

    event.preventDefault();

    //console.log($("#n51hwa").serialize());

    $.ajax({
        type    : "POST",
        url     : "/lessfly/sysmgr/setup/app-register-put",
        data    : $("#n51hwa").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {
                
                lessAlert("#cdnxe5", 'alert-success', rsj.message);

                window.setTimeout(function(){
                    window.location = "/lessfly";
                }, 1500);

            } else {
                lessAlert("#cdnxe5", 'alert-danger', rsj.message);
            }

        },
        error   : function(xhr, textStatus, error) {
            lessAlert("#cdnxe5", 'alert-danger', textStatus+' '+xhr.responseText);
        }
    });
});


</script>

