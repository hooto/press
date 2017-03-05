// Copyright 2015~2017 hooto Author, All rights reserved.
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
	"net/http"

	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"

	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/status"
)

type Setup struct {
	*httpsrv.Controller
}

func (c Setup) IndexAction() {

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect(iamclient.LoginUrl(c.Request.RawAbsUrl()))
		return
	}

	if token := c.Params.Get(iamclient.AccessTokenKey); len(token) >= 16 {
		ck := &http.Cookie{
			Name:     iamclient.AccessTokenKey,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  iamclient.Expired(864000),
		}
		http.SetCookie(c.Response.Out, ck)

		c.Redirect("mgr")
		return
	}

	if status.IamServiceStatus == status.IamServiceUnRegistered {
		c.Data["iam_url"] = iamclient.ServiceUrl

		c.Data["instance_url"] = c.UrlBase("")

		c.Data["app_id"] = config.AppName
		c.Data["app_title"] = config.Config.AppTitle
		c.Data["version"] = config.Version

		c.Render("setup/app-register.tpl")
		return
	}

	c.Redirect("mgr")
}

func (c Setup) AppRegisterPutAction() {

	reg := iamapi.AppInstanceRegister{
		AccessToken: iamclient.SessionAccessToken(c.Session),
		Instance: iamapi.AppInstance{
			Meta: types.ObjectMeta{
				ID: config.Config.InstanceID,
			},
			AppID:      config.AppName,
			AppTitle:   c.Params.Get("app_title"),
			Version:    config.Version,
			Url:        c.Params.Get("instance_url"),
			Privileges: config.Perms,
		},
	}

	defer c.RenderJson(&reg)

	regjs, _ := json.Encode(reg, "")

	hc := httpclient.Put(iamclient.ServiceUrl + "/v1/app-auth/register")
	hc.Body(regjs)

	if err := hc.ReplyJson(&reg); err != nil {

		reg.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}

	} else if reg.Error == nil && reg.Kind == "AppInstanceRegister" {

		config.Config.InstanceID = reg.Instance.Meta.ID
		iamclient.InstanceID = reg.Instance.Meta.ID

		config.Config.AppTitle = reg.Instance.AppTitle

		status.IamServiceStatus = status.IamServiceOK

		if err := config.Save(); err != nil {
			reg.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}
		}
	}

	hc.Close()
}
