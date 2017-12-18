// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamclient"

	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/status"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	status.Locker.RLock()
	defer status.Locker.RUnlock()

	if c.Params.Get("_iam_out") != "" {
		c.Redirect(c.UrlBase(""))
		return
	}

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect(iamclient.AuthServiceUrl(
			config.Config.InstanceID,
			c.UrlBase("hp/auth/cb"),
			c.Request.RawAbsUrl(),
		))
		return
	}

	if status.IamServiceStatus == status.IamServiceUnRegistered {
		c.Redirect("hp/mgr/setup/index")
		return
	}

	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	if v := config.SysConfigList.FetchString("http_h_ac_allow_origin"); v != "" {
		c.Response.Out.Header().Set("Access-Control-Allow-Origin", v)
	}

	c.Data["sys_version_sign"] = config.SysVersionSign

	c.Render("index.tpl")
}
