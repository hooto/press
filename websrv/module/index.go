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

package module

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hooto/httpsrv"

	"github.com/hooto/hpress/config"
)

type Static struct {
	*httpsrv.Controller
}

func (c Static) IndexAction() {

	c.AutoRender = false

	var (
		object_path = strings.TrimPrefix(c.Request.UrlPath(), "/hp/-/static/")
	)

	n := strings.Index(object_path, "/")
	if n < 1 {
		return
	}
	srvname := object_path[:n]

	ext := strings.ToLower(filepath.Ext(object_path))
	switch ext {
	case ".js", ".css":

	case ".jpg", ".png", ".svg", ".git", ".ico":

	default:
		return
	}

	mod, ok := config.Modules[srvname]
	if !ok {
		return
	}

	abs_path := config.Prefix + "/modules/" + mod.Meta.Name + "/static/" + object_path[n+1:]

	if fp, err := os.Open(abs_path); err == nil {

		c.Response.Out.Header().Set("Cache-Control", "max-age=86400")
		http.ServeContent(c.Response.Out, c.Request.Request, object_path, time.Now(), fp)
		fp.Close()

	} else {
		c.RenderError(404, "Object Not Found")
	}
}
