
<nav class="navbar navbar-inverse navbar-static-top" style="margin-bottom:10px;">
  <div class="container-fluid">

    <div class="navbar-header" style="">
        <a class="navbar-brand" href="#">
            <!-- <img alt="Brand" src="/mgr/-/img/logo.png"> -->
        </a>
        <span class="navbar-brand">CMS</span>
    </div>

    <div id="l5s-nav" class="collapse navbar-collapse">
        <ul class="nav navbar-nav">
            <li class="l5s-nav-item active"><a href="#content" onclick="l5sMgr.NodeList()">Content</a></li>
            <li class="l5s-nav-item "><a href="#structure" onclick="l5sMgr.SpecList()">Spec</a></li>
        </ul>

        <ul class="nav navbar-nav navbar-right">
          <li><a href="#signout" onclick="l5sIDS.Logout()">Logout</a></li>
        </ul>
    </div>

  </div>
</nav>

<div id="com-content" class="container-fluid">loading</div>
