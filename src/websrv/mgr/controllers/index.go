package controllers

import (
	"io"
	// "net/http"

	"github.com/lessos/lessgo/httpsrv"
	// "github.com/lessos/lessgo/service/lessids"

	// "../../../state"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false

	// TODO
	// // Check if lessIDS Service Available
	// if state.LessIdsState == state.LessIdsUnavailable {
	// 	c.Data["lessids_url"] = lessids.ServiceUrl
	// 	c.Render("error/lessids.offline.tpl")
	// 	return
	// }

	// //
	// session, err := c.Session.SessionFetch()
	// if err != nil || session.Uid == 0 {
	// 	c.Redirect(lessids.LoginUrl(c.Request.RawAbsUrl()))
	// 	return
	// }

	// //
	// if c.Params.Get("access_token") != "" {

	// 	ck := &http.Cookie{
	// 		Name:     "access_token",
	// 		Value:    session.AccessToken,
	// 		Path:     "/",
	// 		HttpOnly: true,
	// 		Expires:  session.Expired.UTC(),
	// 	}
	// 	http.SetCookie(c.Response.Out, ck)

	// 	c.Redirect("/mgr")
	// 	return
	// }

	//
	// if state.LessIdsState == state.LessIdsUnRegistered {
	// 	c.Redirect("/mgr/setup/index")
	// 	return
	// }

	io.WriteString(c.Response.Out, `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>CMS</title>
  <script src="/mgr/~/lessui/js/sea.js"></script>
  <script src="/mgr/-/js/main.js"></script>
  <script type="text/javascript">
    window.onload = l5sMgr.Boot() ;
  </script>
</head>
<body id="body-content">
loading
</body>
</html>`)

}
