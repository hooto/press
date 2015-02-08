<nav class="navbar navbar-default" role="navigation">
  <div class="container-fluid">

    <div class="navbar-header" style="">
        <!-- <a class="navbar-brand" href="#">
            <img alt="" src="/~/cmf/logo.png" width="20px" height="20px">
        </a> -->
        <span class="navbar-brand">Project name</span>
    </div>

    <div class="collapse navbar-collapse">
        <ul class="nav navbar-nav">
            {{range $v := .topnav}}
            <li class=""><a href="{{Field $v "content"}}">{{Field $v "title"}}</a></li>
            {{end}}
        </ul>

        <ul class="nav navbar-nav navbar-right">
          <li><a href="#">Login</a></li>
        </ul>
    </div>

  </div>
</nav>