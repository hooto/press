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

package frontend

import (
	"path/filepath"
	"strings"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/x/webui"

	"../../api"
	"../../conf"
	"../../datax"
)

type Index struct {
	*httpsrv.Controller
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

	return "", "", false
}

func (c Index) IndexAction() {

	c.AutoRender = false

	var (
		srvname = "general"
		uris    = strings.Split(strings.Trim(filepath.Clean(c.Request.RequestPath), "/"), "/")
	)

	if len(uris) < 1 {
		uris = append(uris, "general")
	} else {
		srvname = uris[0]
	}

	if len(uris) < 2 {
		uris = append(uris, "")
	}

	mod, ok := conf.Modules[srvname]
	if !ok {
		srvname = "general"
		mod, ok = conf.Modules[srvname]
		if !ok {
			return
		}
	}

	dataAction, template, mat := c.filter(uris[1:], mod)
	if !mat {

		srvname = "general"

		if uris[1] == "" {
			template = "index.tpl"
		} else {
			template = "404.tpl"
		}
	}

	c.Data["baseuri"] = "/" + srvname
	c.Data["modname"] = mod.SrvName

	if dataAction != "" {

		for _, action := range mod.Actions {

			if action.Name != dataAction {
				continue
			}

			for _, datax := range action.Datax {
				c.dataRender(srvname, datax)
				c.Data["__datax_table__"] = datax.Query.Table
			}

			break
		}
	}

	c.Render(mod.Meta.Name, template)
}

func (c Index) dataRender(srvname string, ad api.ActionData) {

	mod, ok := conf.Modules[srvname]
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

	if id := c.Params.Get("id"); id != "" {
		if len(id) > 5 && id[len(id)-5:] == ".html" {
			id = id[:len(id)-5]
		}
		qry.Filter("id", id)
	}

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
						qry.Filter("term_"+term.Meta.Name, termVal)
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

		ls := qry.NodeList()

		c.Data[ad.Name] = ls

		if qry.Pager {
			pager := webui.NewPager(0,
				uint64(ls.Meta.TotalResults),
				uint64(ls.Meta.ItemsPerList),
				10)
			pager.CurrentPageNumber = uint64(page)
			c.Data[ad.Name+"_pager"] = pager
		}

	case "node.entry":

		c.Data[ad.Name] = qry.NodeEntry()

	case "term.list":

		ls := qry.TermList()
		c.Data[ad.Name] = ls

		if qry.Pager {
			c.Data[ad.Name+"_pager"] = webui.NewPager(0,
				uint64(ls.Meta.TotalResults),
				uint64(ls.Meta.ItemsPerList),
				10)
		}

	case "term.entry":

		c.Data[ad.Name] = qry.TermEntry()
	}
}
