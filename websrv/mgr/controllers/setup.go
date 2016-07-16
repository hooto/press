package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"../../../config"
	"../../../status"
)

type Setup struct {
	*httpsrv.Controller
}

func (c Setup) IndexAction() {

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect(iamclient.LoginUrl(c.Request.RawAbsUrl()))
		return
	}

	if token := c.Params.Get("access_token"); token != "" {

		ck := &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  iamclient.Expired(864000),
		}
		http.SetCookie(c.Response.Out, ck)

		c.Redirect("/mgr")
		return
	}

	if status.IamServiceStatus == status.IamServiceUnRegistered {

		c.Data["iam_url"] = iamclient.ServiceUrl

		host := c.Request.Host
		if i := strings.Index(host, ":"); i > 0 {
			host = host[:i]
		}

		insturl := "http://" + host
		if config.Config.HttpPort != 80 {
			insturl += fmt.Sprintf(":%d", config.Config.HttpPort)
		}
		c.Data["instance_url"] = insturl

		c.Data["app_id"] = "lesscms"
		c.Data["app_title"] = config.Config.AppTitle
		c.Data["version"] = config.Version

		c.Render("setup/app-register.tpl")
		return
	}

	c.Redirect("/mgr")
}

func (c Setup) AppRegisterPutAction() {

	reg := iamapi.AppInstanceRegister{
		AccessToken: iamclient.SessionAccessToken(c.Session),
		Instance: iamapi.AppInstance{
			Meta: types.ObjectMeta{
				ID: config.Config.InstanceID,
			},
			AppID:      "lesscms",
			AppTitle:   c.Params.Get("app_title"),
			Version:    config.Version,
			Url:        c.Params.Get("instance_url"),
			Privileges: config.Perms,
		},
	}

	defer c.RenderJson(&reg)

	regjs, _ := utils.JsonEncode(reg)

	// fmt.Println(regjs)

	hc := httpclient.Put(iamclient.ServiceUrl + "/v1/app-auth/register")
	hc.Body(regjs)

	if err := hc.ReplyJson(&reg); err != nil {

		reg.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}

	} else if reg.Error == nil && reg.Kind == "AppInstanceRegister" {

		config.Config.InstanceID = reg.Instance.Meta.ID
		config.Config.AppTitle = reg.Instance.AppTitle

		status.IamServiceStatus = status.IamServiceOK

		if err := config.Save(); err != nil {
			reg.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}
		}
	}

	hc.Close()
}
