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

package v1

import (
	"fmt"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"
	"github.com/lessos/lessids/idsapi"

	"../../api"
	"../../conf"
	"../../datax"
)

type Node struct {
	*httpsrv.Controller
}

func (c Node) ListAction() {

	rsp := api.NodeList{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.list") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)

	rsp = dq.NodeList()
}

func (c Node) EntryAction() {

	rsp := api.Node{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.read") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.NodeEntry()
}

func (c Node) SetAction() {

	rsp := api.Node{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.write") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	model, err := conf.SpecNodeModel(c.Params.Get("modname"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "404",
			Message: "Spec or Model Not Found",
		}
		return
	}

	if err := c.Request.JsonDecode(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Bad Request: " + err.Error(),
		}
		return
	}

	var (
		set = map[string]interface{}{}
	)

	//
	table := fmt.Sprintf("nx%s_%s", utils.StringEncode16(c.Params.Get("modname"), 12), c.Params.Get("modelid"))

	if len(rsp.ID) > 0 {

		q := rdobase.NewQuerySet().From(table).Limit(1)
		q.Where.And("id", rsp.ID)
		rs, err := dcn.Base.Query(q)
		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: "Can not pull database instance",
			}
			return
		}

		if len(rs) < 1 {
			rsp.Error = &types.ErrorMeta{
				Code:    "404",
				Message: "Node Not Found",
			}
			return
		}

		if rs[0].Field("title").String() != rsp.Title {
			set["title"] = rsp.Title
		}

		if rs[0].Field("state").Int16() != rsp.State {
			set["state"] = rsp.State
		}

		//
		for _, valField := range rsp.Fields {

			for _, modField := range model.Fields {

				if modField.Name != valField.Name {
					continue
				}

				if rs[0].Field("field_"+valField.Name).String() != valField.Value {

					set["field_"+valField.Name] = valField.Value

					if modField.Type == "text" {

						attrs := []api.KeyValue{}

						for _, attr := range valField.Attrs {
							if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
								attrs = append(attrs, api.KeyValue{attr.Key, attr.Value})
							}
						}

						set["field_"+valField.Name+"_attrs"], _ = utils.JsonEncode(attrs)
					}
				}

				break
			}
		}

		//
		for _, modTerm := range model.Terms {

			for _, term := range rsp.Terms {

				if modTerm.Meta.Name != term.Name {
					continue
				}

				switch modTerm.Type {

				case api.TermTag:

					tags, _ := datax.TermSync(c.Params.Get("modname"), modTerm.Meta.Name, term.Value)

					if rs[0].Field("term_"+term.Name).String() != term.Value {
						set["term_"+modTerm.Meta.Name] = tags.Content()
						set["term_"+modTerm.Meta.Name+"_idx"] = tags.Index()
					}

				case api.TermTaxonomy:

					set["term_"+modTerm.Meta.Name] = term.Value
				}
			}
		}

	} else {

		set["id"] = utils.StringNewRand(12)
		set["title"] = rsp.Title
		set["state"] = rsp.State
		set["created"] = rdobase.TimeNow("datetime")
		set["userid"] = "dr5a8pgv"

		//
		for _, valField := range rsp.Fields {

			for _, modField := range model.Fields {

				if modField.Name != valField.Name {
					continue
				}

				set["field_"+valField.Name] = valField.Value

				if modField.Type == "text" {
					attrs := []api.KeyValue{}

					for _, attr := range valField.Attrs {
						if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
							attrs = append(attrs, api.KeyValue{attr.Key, attr.Value})
						}
					}

					set["field_"+valField.Name+"_attrs"], _ = utils.JsonEncode(attrs)
				}

				break
			}
		}

		//
		for _, modTerm := range model.Terms {

			for _, term := range rsp.Terms {

				if modTerm.Meta.Name != term.Name {
					continue
				}

				switch modTerm.Type {

				case api.TermTag:

					tags, _ := datax.TermSync(c.Params.Get("modname"), modTerm.Meta.Name, term.Value)
					set["term_"+modTerm.Meta.Name] = tags.Content()
					set["term_"+modTerm.Meta.Name+"_idx"] = tags.Index()

				case api.TermTaxonomy:

					set["term_"+modTerm.Meta.Name] = term.Value
				}
			}
		}
	}

	if len(set) > 0 {

		set["updated"] = rdobase.TimeNow("datetime")

		if len(rsp.ID) > 0 {

			ft := rdobase.NewFilter()
			ft.And("id", rsp.ID)
			_, err = dcn.Base.Update(table, set, ft)

		} else {
			rsp.ID = set["id"].(string)
			_, err = dcn.Base.Insert(table, set)
		}

		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: err.Error(),
			}
			return
		}
	}

	rsp.Kind = "Node"
}
