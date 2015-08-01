<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Setup : Registration to Identity Service</title>
  <link rel="stylesheet" href="/~/bootstrap/3.3/css/bootstrap.min.css" type="text/css">
  <link rel="stylesheet" href="/~/css/main.css" type="text/css">
  <script src="/~/jquery/1.11/jquery.min.js"></script>
  <script src="/~/lessui/js/lessui.js"></script>
  <script src="/~/lessui/js/sea.js"></script>
  <script src="/~/js/main.js"></script>
  <script type="text/javascript">
    window.onload_hooks = [];
  </script>
</head>

<body>

<div class="container">

<div class="panel panel-default">
  
  <div class="panel-heading"><strong>Setup: Register Application to Identity Service</strong></div>
  
  <div class="panel-body">

    <form id="l5s-app-reg" action="#">
      
      <div id="l5s-app-reg-alert" class="alert alert-danger hide">...</div>
    
      <div class="form-group">
        <label>Identity Service Url</label>
        <input type="text" name="ids_url" class="form-control" placeholder="Enter the Identity Service URL" value="{{.ids_url}}" readonly>
      </div>
    
<!--       <div class="form-group">
        <label>Instance ID</label>
        <input type="text" name="instance_id" class="form-control" value="{{.instance_id}}" readonly>
      </div>
 -->    
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

</body>
</html>

<script type="text/javascript">


$("#l5s-app-reg").submit(function(event) {

    event.preventDefault();

    var alertid = "#l5s-app-reg-alert";

    $.ajax({
        type    : "POST",
        url     : "/mgr/setup/app-register-put",
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
                    window.location = "/mgr";
                }, 1500);
            }
        },
        error   : function(xhr, textStatus, error) {
            l4i.InnerAlert(alertid, 'alert-danger', textStatus+' '+xhr.responseText);
        }
    });
});


</script>

