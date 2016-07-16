// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"fmt"
	"io"

	"github.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"

	"../../../config"
	"../../../status"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false

	//
	status.Locker.RLock()
	defer status.Locker.RUnlock()

	if status.IamServiceStatus == status.IamServiceUnRegistered {
		c.Redirect(config.HttpSrvBasePath("mgr/setup/index"))
		return
	}

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect(iamclient.AuthServiceUrl(
			config.Config.InstanceID,
			fmt.Sprintf("//%s%s/auth/cb", c.Request.Host, config.HttpSrvBasePath("")),
			c.Request.RawAbsUrl()))
		return
	}

	//
	io.WriteString(c.Response.Out, `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>CMS</title>
  <script src="`+config.HttpSrvBasePath("mgr/~/lessui/js/sea.js")+`"></script>
  <script src="`+config.HttpSrvBasePath("mgr/-/js/main.js")+`"></script>
  <script type="text/javascript">
    window._basepath = "`+config.HttpSrvBasePath("")+`";
  </script>
</head>
<body id="body-content">
loading
</body>
<script type="text/javascript">
window.onload = l5sMgr.Boot();
</script>
</html>`)

}
