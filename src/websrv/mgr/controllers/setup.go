package controllers

import (
	"fmt"
	"io"
	// "net/http"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/service/lessids"
	"github.com/lessos/lessgo/utils"

	"../../../conf"
	"../../../state"
)

type Setup struct {
	*httpsrv.Controller
}

func (c Setup) IndexAction() {

	// TODO
	// session, err := c.Session.SessionFetch()
	// if err != nil || session.Uid == 0 {
	// 	c.Redirect(lessids.LoginUrl(c.Request.RawAbsUrl()))
	// 	return
	// }

	// if c.Params.Get("access_token") != "" {

	// 	ck := &http.Cookie{
	// 		Name:     "access_token",
	// 		Value:    session.AccessToken,
	// 		Path:     "/",
	// 		HttpOnly: true,
	// 		Expires:  session.Expired.UTC(),
	// 	}
	// 	http.SetCookie(c.Response.Out, ck)

	// 	c.Redirect("/lesscms")
	// 	return
	// }

	if state.LessIdsState == state.LessIdsUnRegistered {

		c.Data["lessids_url"] = lessids.ServiceUrl

		c.Data["instance_id"] = conf.Config.InstanceID

		insturl := "http://" + conf.Config.HttpAddr
		if conf.Config.HttpPort != 80 {
			insturl += fmt.Sprintf(":%d", conf.Config.HttpPort)
		}
		c.Data["instance_url"] = insturl + "/lesscms"

		c.Data["app_id"] = "lesscms"
		c.Data["app_title"] = "lessCMS"
		c.Data["version"] = conf.Config.Version

		c.Render("setup/app-register.tpl")
		return
	}

	c.Redirect("/lesscms")
}

func (c Setup) AppRegisterPutAction() {

	c.AutoRender = false

	var rsp = struct {
		Status  int
		Message string
	}{
		400,
		"Bad Request",
	}

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	// session, err := c.Session.SessionFetch()
	// if err != nil || session.Uid == 0 {
	// 	rsp.Status = 401
	// 	rsp.Message = "Unauthorized"
	// 	return
	// }

	if c.Params.Get("instance_id") == "" ||
		c.Params.Get("app_id") == "" {
		return
	}

	var req struct {
		AccessToken string `json:"access_token"`
		Data        struct {
			InstanceID  string           `json:"instance_id"`
			InstanceUrl string           `json:"instance_url"`
			AppId       string           `json:"app_id"`
			AppTitle    string           `json:"app_title"`
			Version     string           `json:"version"`
			Privileges  []conf.Privilege `json:"privileges"`
		} `json:"data"`
	}

	req.AccessToken = c.Session.AccessToken
	req.Data.InstanceID = c.Params.Get("instance_id")
	req.Data.InstanceUrl = c.Params.Get("instance_url")
	req.Data.AppId = "lesscms"
	req.Data.AppTitle = c.Params.Get("app_title")
	req.Data.Version = conf.Config.Version
	req.Data.Privileges = conf.Privileges

	reqstr, err := utils.JsonEncode(req)
	if err != nil {
		return
	}

	hc := httpclient.Put(c.Params.Get("lessids_url") +
		"/app-auth/register?access_token=" + c.Session.AccessToken)
	hc.Body(reqstr)

	regstr, err := hc.ReplyString()
	var reg struct {
		Status  int
		Message string
	}
	err = utils.JsonDecode(regstr, &reg)
	if reg.Status != 200 {
		rsp.Status = reg.Status
		rsp.Message = reg.Message
	} else {
		state.Refresh()
		rsp.Status = 200
		rsp.Message = "Successfully registered"
	}
}
