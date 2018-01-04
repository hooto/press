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

package frontend

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/x/webui"
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/datax"
	"github.com/hooto/hpress/store"
)

type Index struct {
	*httpsrv.Controller
	hookPosts []func()
	us        iamapi.UserSession
}

func (c *Index) Init() int {
	c.us, _ = iamclient.SessionInstance(c.Session)
	return 0
}

func (c Index) filter(rt []string, spec *api.Spec) (string, string, bool) {

	for _, route := range spec.Router.Routes {

		matlen, params := 0, map[string]string{}

		for i, node := range route.Tree {

			if len(node) < 1 || i >= len(rt) {
				break
			}

			if node[0] == ':' {

				params[node[1:]] = rt[i]

				matlen++

			} else if node == rt[i] {

				matlen++
			}
		}

		if matlen == len(route.Tree) {

			for k, v := range params {
				c.Params.Values[k] = append(c.Params.Values[k], v)
			}

			return route.DataAction, route.Template, true
		}
	}

	for _, route := range spec.Router.Routes {
		if route.Default {
			return route.DataAction, route.Template, true
		}
	}

	return "", "", false
}

var (
	srvname_default = "core-genereal"
	uris_default    = []string{"core-general"}
)

func (c Index) IndexAction() {

	c.AutoRender = false
	start := time.Now().UnixNano()

	if v := config.SysConfigList.FetchString("http_h_ac_allow_origin"); v != "" {
		c.Response.Out.Header().Set("Access-Control-Allow-Origin", v)
	}

	var (
		reqpath = filepath.Clean("/" + c.Request.RequestPath)
		uris    = []string{}
	)
	if reqpath == "" || reqpath == "." {
		reqpath = "/"
	}
	if len(reqpath) > 0 && reqpath != "/" {
		uris = strings.Split(strings.Trim(reqpath, "/"), "/")
	}

	if len(uris) < 1 {
		if config.RouterBasepathDefault != "/" {
			reqpath = config.RouterBasepathDefault
			uris = config.RouterBasepathDefaults
		} else {
			uris = uris_default
		}
	}
	srvname := uris[0]

	if len(uris) < 2 {
		uris = append(uris, "")
	}
	// fmt.Println(uris, srvname, c.Params.Get("referid"), c.Params.Get("id"))

	mod, ok := config.Modules[srvname]
	if !ok {
		srvname = srvname_default
		mod, ok = config.Modules[srvname]
		if !ok {
			return
		}
	}

	dataAction, template, mat := c.filter(uris[1:], mod)
	if !mat {
		if uris[1] == "" {
			template = "index.tpl"
		} else {
			template = "404.tpl"
		}
	}

	// if session, err := c.Session.Instance(); err == nil {
	// 	c.Data["session"] = session
	// }

	c.Data["baseuri"] = "/" + srvname
	c.Data["http_request_path"] = reqpath
	c.Data["srvname"] = srvname
	c.Data["modname"] = mod.Meta.Name
	c.Data["sys_version_sign"] = config.SysVersionSign
	if c.us.IsLogin() {
		c.Data["s_user"] = c.us.UserName
	}

	if dataAction != "" {

		for _, action := range mod.Actions {

			if action.Name != dataAction {
				continue
			}

			for _, datax := range action.Datax {
				c.dataRender(srvname, action.Name, datax)
				c.Data["__datax_table__"] = datax.Query.Table
			}

			break
		}
	}

	// render_start := time.Now()
	c.Render(mod.Meta.Name, template)

	// fmt.Println("render in-time", mod.Meta.Name, template, time.Since(render_start))

	c.RenderString(fmt.Sprintf("<!-- rt-time/db+render : %d ms -->", (time.Now().UnixNano()-start)/1e6))

	// fmt.Println("hookPosts", len(c.hookPosts))
	for _, fn := range c.hookPosts {
		fn()
	}
}

func (c *Index) dataRender(srvname, action_name string, ad api.ActionData) {

	mod, ok := config.Modules[srvname]
	if !ok {
		return
	}

	qry := datax.NewQuery(mod.Meta.Name, ad.Query.Table)
	if ad.Query.Limit > 0 {
		qry.Limit(ad.Query.Limit)
	}

	if ad.Query.Order != "" {
		qry.Order(ad.Query.Order)
	}

	qry.Filter("status", 1)

	qry.Pager = ad.Pager

	switch ad.Type {

	case "node.list":

		for _, modNode := range mod.NodeModels {

			if ad.Query.Table != modNode.Meta.Name {
				continue
			}

			for _, term := range modNode.Terms {

				if termVal := c.Params.Get("term_" + term.Meta.Name); termVal != "" {

					switch term.Type {

					case api.TermTaxonomy:

						if idxs := datax.TermTaxonomyCacheIndexes(mod.Meta.Name, term.Meta.Name, termVal); len(idxs) > 1 {
							args := []interface{}{}
							for _, idx := range idxs {
								args = append(args, idx)
							}
							qry.Filter("term_"+term.Meta.Name+".in", args...)
						} else {
							qry.Filter("term_"+term.Meta.Name, termVal)
						}

						c.Data["term_"+term.Meta.Name] = termVal

					case api.TermTag:
						// TOPO
						qry.Filter("term_"+term.Meta.Name+".like", "%"+termVal+"%")
						c.Data["term_"+term.Meta.Name] = termVal
					}
				}
			}

			break
		}

		page := c.Params.Int64("page")
		if page > 1 {
			qry.Offset(ad.Query.Limit * (page - 1))
		}

		if c.Params.Get("qry_text") != "" {
			qry.Filter("title.like", "%"+c.Params.Get("qry_text")+"%")
			c.Data["qry_text"] = c.Params.Get("qry_text")
		}

		var ls api.NodeList
		qryhash := qry.Hash()

		if ad.CacheTTL > 0 && (!c.us.IsLogin() || c.us.UserName != config.Config.AppInstance.Meta.User) {
			if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
				rs.Decode(&ls)
			}
		}

		if len(ls.Items) == 0 {

			ls = qry.NodeList([]string{}, []string{})
			// fmt.Println("index node.list")
			if ad.CacheTTL > 0 && len(ls.Items) > 0 {
				c.hookPosts = append(
					c.hookPosts,
					func() {
						store.LocalCache.KvPut([]byte(qryhash), &ls, &skv.KvWriteOptions{
							Ttl: ad.CacheTTL,
						})
					},
				)
			}
		}

		c.Data[ad.Name] = ls

		if qry.Pager {
			pager := webui.NewPager(uint64(page),
				uint64(ls.Meta.TotalResults),
				uint64(ls.Meta.ItemsPerList),
				10)
			pager.CurrentPageNumber = uint64(page)
			c.Data[ad.Name+"_pager"] = pager
		}

	case "node.entry":

		id := c.Params.Get(ad.Name + "_id")
		if id == "" {
			id = c.Params.Get("id")
			if id == "" {
				return
			}
		}

		nodeModel, err := config.SpecNodeModel(mod.Meta.Name, ad.Query.Table)
		if err != nil {
			return
		}

		node_refer := ""
		if nodeModel.Extensions.NodeRefer != "" {
			if mv, ok := c.Data[action_name+"_nsr_"+nodeModel.Extensions.NodeRefer]; ok {
				node_refer = mv.(string)
			}
		}

		id_ext := ""
		if i := strings.LastIndex(id, "."); i > 0 {
			id_ext = id[i:]
			id = id[:i]
		}

		if id_ext == ".html" {
			qry.Filter("id", id)
		} else if nodeModel.Extensions.Permalink != "" {
			if nodeModel.Extensions.NodeRefer != "" && node_refer == "" {
				return
			}
			qry.Filter("ext_permalink_idx", idhash.HashToHexString([]byte(node_refer+id), 12))
		} else {
			return
		}

		var entry api.Node
		qryhash := qry.Hash()
		if ad.CacheTTL > 0 && (!c.us.IsLogin() || c.us.UserName != config.Config.AppInstance.Meta.User) {
			if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
				rs.Decode(&entry)
			}
		}

		if entry.ID == "" {
			entry = qry.NodeEntry()
			if ad.CacheTTL > 0 && entry.Title != "" {
				c.hookPosts = append(
					c.hookPosts,
					func() {
						store.LocalCache.KvPut([]byte(qryhash), &entry, &skv.KvWriteOptions{
							Ttl: ad.CacheTTL,
						})
					},
				)
			}
		}

		if entry.ID == "" {
			return
		}

		if nodeModel.Extensions.AccessCounter {

			if ips := strings.Split(c.Request.RemoteAddr, ":"); len(ips) > 1 {

				table := fmt.Sprintf("nx%s_%s", idhash.HashToHexString([]byte(mod.Meta.Name), 12), ad.Query.Table)
				store.LocalCache.KvPut([]byte("access_counter/"+table+"/"+ips[0]+"/"+entry.ID), "1", nil)
			}
		}

		if nodeModel.Extensions.NodeSubRefer != "" {
			// fmt.Println("setting", action_name, ad.Query.Table, nodeModel.Extensions.NodeSubRefer, "_id", entry.ID)
			c.Data[action_name+"_nsr_"+ad.Query.Table] = entry.ID
		}

		if entry.Title != "" {
			c.Data["__html_head_title__"] = datax.StringSub(datax.TextHtml2Str(entry.Title), 0, 50)
		}

		c.Data[ad.Name] = entry

	case "term.list":

		var ls api.TermList
		qryhash := qry.Hash()
		if ad.CacheTTL > 0 {
			if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
				rs.Decode(&ls)
			}
		}

		if len(ls.Items) == 0 {
			ls = qry.TermList()
			if ad.CacheTTL > 0 && len(ls.Items) > 0 {
				store.LocalCache.KvPut([]byte(qryhash), ls, &skv.KvWriteOptions{
					Ttl: ad.CacheTTL,
				})
			}
		}

		c.Data[ad.Name] = ls

		if qry.Pager {
			c.Data[ad.Name+"_pager"] = webui.NewPager(0,
				uint64(ls.Meta.TotalResults),
				uint64(ls.Meta.ItemsPerList),
				10)
		}

	case "term.entry":

		var entry api.Term
		qryhash := qry.Hash()

		if ad.CacheTTL > 0 {
			if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
				rs.Decode(&entry)
			}
		}

		if entry.Title == "" {
			entry = qry.TermEntry()
			if ad.CacheTTL > 0 && entry.Title != "" {
				store.LocalCache.KvPut([]byte(qryhash), entry, &skv.KvWriteOptions{
					Ttl: ad.CacheTTL,
				})
			}
		}

		c.Data[ad.Name] = entry
	}
}
