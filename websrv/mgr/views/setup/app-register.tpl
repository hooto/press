<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Setup : Registration to IAM Service</title>
  <link rel="stylesheet" href="{{HttpSrvBasePath "hp/~/bs/3.3/css/bootstrap.css"}}" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="{{HttpSrvBasePath "hp/~/hp/img/ap.ico"}}">
  <script src="{{HttpSrvBasePath "hp/~/hp/js/jquery.js"}}"></script>
  <script src="{{HttpSrvBasePath "hp/~/lessui/js/lessui.js"}}"></script>
  <script type="text/javascript">
    window.onload_hooks = [];
  </script>
</head>

<body>

<div class="container" style="width:600px">

<div class="hp-mgr-setup-logo" style="text-align:center;padding:30px;">
  <img src="{{HttpSrvBasePath "hp/~/hp/img/alpha2.png"}}">
</div>

<div class="panel panel-default">
  
  <div class="panel-heading">
    <h3 style="margin: 10px 0;">Setup</h3>
    <strong>Register Application to IAM (Identity &amp; Access Management) Service</strong>
  </div>
  
  <div class="panel-body">

    <form id="hp-app-reg" action="#">
      
      <div id="hp-app-reg-alert" class="alert alert-danger hide">...</div>
    
      <div class="form-group">
        <label>IAM Service URL</label>
        <input type="text" name="iam_url" class="form-control" placeholder="Enter the IAM Service URL" value="{{.iam_url}}" readonly>
      </div>
    
<!--       <div class="form-group">
        <label>Instance ID</label>
        <input type="text" name="instance_id" class="form-control" value="{{.instance_id}}" readonly>
      </div>
 -->    
      <div class="form-group">
        <label>Instance Frontend URL</label>
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
    
      <button type="submit" class="btn btn-success btn-block" style="margin-top:30px">Commit</button>
    </form>
  </div>

</div>

</div>

</body>
</html>

<script type="text/javascript">


$("#hp-app-reg").submit(function(event) {

    event.preventDefault();

    var alertid = "#hp-app-reg-alert";

    $.ajax({
        type    : "POST",
        url     : "{{HttpSrvBasePath "hp/mgr/setup/app-register-sync"}}",
        data    : $(this).serialize(),
        timeout : 3000,
        success : function(data) {

            if (!data || data.kind != "AppInstanceRegister") {

                if (data.error) {
                    l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
                } else {
                    l4i.InnerAlert(alertid, 'alert-danger', 'Network Exception');
                }
            
            } else {

                l4i.InnerAlert(alertid, 'alert-success', 'Successfully registered ...');

                window.setTimeout(function(){
                    window.location = "{{HttpSrvBasePath "hp/mgr"}}";
                }, 1500);
            }
        },
        error   : function(xhr, textStatus, error) {
            l4i.InnerAlert(alertid, 'alert-danger', textStatus+' '+xhr.responseText);
        }
    });
});


</script>

