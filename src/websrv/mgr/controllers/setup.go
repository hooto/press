package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"

	"../../../conf"
	"../../../status"
)

type Setup struct {
	*httpsrv.Controller
}

func (c Setup) IndexAction() {

	if !idclient.IsLogin(c.Session.AccessToken) {
		c.Redirect(idclient.LoginUrl(c.Request.RawAbsUrl()))
		return
	}

	if token := c.Params.Get("access_token"); token != "" {

		ck := &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  idclient.Expired(864000),
		}
		http.SetCookie(c.Response.Out, ck)

		c.Redirect("/mgr")
		return
	}

	if status.IdentityServiceStatus == status.IdentityServiceUnRegistered {

		c.Data["ids_url"] = idclient.ServiceUrl

		host := c.Request.Host
		if i := strings.Index(host, ":"); i > 0 {
			host = host[:i]
		}

		insturl := "http://" + host
		if conf.Config.HttpPort != 80 {
			insturl += fmt.Sprintf(":%d", conf.Config.HttpPort)
		}
		c.Data["instance_url"] = insturl

		c.Data["app_id"] = "lesscms"
		c.Data["app_title"] = conf.Config.AppTitle
		c.Data["version"] = conf.Version

		c.Render("setup/app-register.tpl")
		return
	}

	c.Redirect("/mgr")
}

func (c Setup) AppRegisterPutAction() {

	reg := idsapi.AppInstanceRegister{
		AccessToken: c.Session.AccessToken,
		Instance: idsapi.AppInstance{
			Meta: types.ObjectMeta{
				ID: conf.Config.InstanceID,
			},
			AppID:      "lesscms",
			AppTitle:   c.Params.Get("app_title"),
			Version:    conf.Version,
			Url:        c.Params.Get("instance_url"),
			Privileges: conf.Perms,
		},
	}

	defer c.RenderJson(&reg)

	regjs, _ := utils.JsonEncode(reg)

	// fmt.Println(regjs)

	hc := httpclient.Put(idclient.ServiceUrl + "/v1/app-auth/register")
	hc.Body(regjs)

	if err := hc.ReplyJson(&reg); err != nil {

		reg.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, err.Error()}

	} else if reg.Error == nil && reg.Kind == "AppInstanceRegister" {

		conf.Config.InstanceID = reg.Instance.Meta.ID
		conf.Config.AppTitle = reg.Instance.AppTitle

		status.IdentityServiceStatus = status.IdentityServiceOK

		if err := conf.Save(); err != nil {
			reg.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, err.Error()}
		}
	}
}
