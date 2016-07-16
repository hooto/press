<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <title>IAM Service Offline</title>
</head>

<body>

<style>
body {
    width: 680px;
    margin: 0 auto;
    background-color: #eee;
}
.err-brw {
    margin: 40px;
    padding: 20px;
    width: 600px;
    border: 1px solid #ccc;
    background-color: #fff;
    border-radius: 4px;
}
.err-brw td {
    padding: 10px 20px 10px 0;
}
.err-brw .imgs1 {
    width: 32px; height: 32px; 
}
.err-footer {
    width: 600px;
    text-align: center;
}
</style>

<div class="err-brw">
  <div class="">
    <div class="alert alert-danger">{{T . "Service Unavailable"}}</div>
    
    <p>{{T . "iam-unavailable-desc" .iam_url}}</p>    
  </div>

</div>

<div class="err-footer">
    &copy; 2015 <a href="http://lessos.com" target="_blank">lessOS.com</a>
</div>

</body>

</html>
