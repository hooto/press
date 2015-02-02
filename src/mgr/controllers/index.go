package controllers

import (
	"github.com/lessos/lessgo/pagelet"
	"io"
)

type Index struct {
	*pagelet.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false

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
