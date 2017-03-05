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

package v1

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"code.hooto.com/hooto/hootopress/api"
	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/datax"
)

type Node struct {
	*httpsrv.Controller
}

var (
	node_list_limit int64 = 15
	reg_permalink         = regexp.MustCompile("^[0-9a-z_-]{1,100}$")
	node_set_lock   sync.Mutex
)

func (c Node) ListAction() {

	ls := api.NodeList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.list", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(node_list_limit)
	dq.Filter("status.gt", 0)

	page := c.Params.Int64("page")
	if page < 1 {
		page = 1
	}

	if page > 1 {
		dq.Offset(int64((page - 1) * node_list_limit))
	}

	dqc := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dqc.Filter("status.gt", 0)

	if c.Params.Get("qry_text") != "" {
		dq.Filter("title.like", "%"+c.Params.Get("qry_text")+"%")
		dqc.Filter("title.like", "%"+c.Params.Get("qry_text")+"%")
	}

	count, err := dqc.NodeCount()
	if err != nil {
		ls.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	ls = dq.NodeList()

	ls.Meta.TotalResults = uint64(count)
	ls.Meta.StartIndex = uint64((page - 1) * node_list_limit)
	ls.Meta.ItemsPerList = uint64(node_list_limit)
}

func (c Node) EntryAction() {

	rsp := api.Node{}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.read", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)
	dq.Filter("status.gt", 0)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.NodeEntry()
}

func (c Node) SetAction() {

	rsp := api.Node{}

	defer c.RenderJson(&rsp)

	node_set_lock.Lock()
	defer node_set_lock.Unlock()

	if !iamclient.SessionAccessAllowed(c.Session, "editor.write", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
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

	model, err := config.SpecNodeModel(c.Params.Get("modname"), c.Params.Get("modelid"))
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
	table := fmt.Sprintf("nx%s_%s", idhash.HashToHexString([]byte(c.Params.Get("modname")), 12), c.Params.Get("modelid"))

	if model.Extensions.Permalink != "" {

		rsp.ExtPermalinkName = strings.Replace(strings.ToLower(strings.TrimSpace(rsp.ExtPermalinkName)), " ", "-", -1)

		if !reg_permalink.MatchString(rsp.ExtPermalinkName) {
			rsp.Error = &types.ErrorMeta{
				Code:    "400",
				Message: "Invalid Permalink Name",
			}
			return
		}

		// set["ext_permalink_name"] = rsp.ExtPermalinkName

		permaname := rsp.ExtPermalinkName

		for i := 0; i < 10; i++ {

			if i > 0 {
				permaname = fmt.Sprintf("%s-%d", rsp.ExtPermalinkName, i)
			}

			permaidx := idhash.HashToHexString([]byte(permaname), 12)

			q := rdobase.NewQuerySet().From(table).Limit(1)
			q.Where.And("ext_permalink_idx", permaidx)

			if len(rsp.ID) > 0 {
				q.Where.And("id.ne", rsp.ID)
			}

			if rs, err := dcn.Base.Query(q); err == nil && len(rs) < 1 {

				set["ext_permalink_name"] = permaname
				set["ext_permalink_idx"] = permaidx
				break
			}
		}

		if _, ok := set["ext_permalink_idx"]; !ok {

			rsp.Error = &types.ErrorMeta{
				Code:    "400",
				Message: "Permalink Name Conflict",
			}
			return
		}
	}

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

		if rs[0].Field("status").Int16() != rsp.Status {
			set["status"] = rsp.Status
		}

		//
		for _, valField := range rsp.Fields {

			for _, modField := range model.Fields {

				if modField.Name != valField.Name {
					continue
				}

				if rs[0].Field("field_"+modField.Name).String() != valField.Value {
					set["field_"+modField.Name] = valField.Value
				}

				if modField.Type == "text" {

					attrs := []api.KeyValue{}

					for _, attr := range valField.Attrs {
						if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
							attrs = append(attrs, api.KeyValue{attr.Key, attr.Value})
						}
					}

					attrs_js, _ := json.Encode(attrs, "  ")

					if string(attrs_js) != rs[0].Field("field_"+modField.Name+"_attrs").String() {
						set["field_"+modField.Name+"_attrs"] = string(attrs_js)
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

		set["id"] = idhash.RandHexString(12)
		set["title"] = rsp.Title
		set["status"] = rsp.Status
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

					jsb, _ := json.Encode(attrs, "  ")
					set["field_"+valField.Name+"_attrs"] = string(jsb)
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

	if model.Extensions.CommentPerEntry {
		if rsp.ExtCommentPerEntry {
			set["ext_comment_perentry"] = 1
		} else {
			set["ext_comment_perentry"] = 0
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

func (c Node) DelAction() {

	rsp := api.Node{}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.write", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
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

	if _, err := config.SpecNodeModel(c.Params.Get("modname"), c.Params.Get("modelid")); err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "404",
			Message: "Spec or Model Not Found",
		}
		return
	}

	//
	set := map[string]interface{}{
		"updated": rdobase.TimeNow("datetime"),
		"status":  0,
	}

	//
	table := fmt.Sprintf("nx%s_%s", idhash.HashToHexString([]byte(c.Params.Get("modname")), 12), c.Params.Get("modelid"))

	//
	ids := strings.Split(c.Params.Get("id"), ",")

	for _, id := range ids {

		q := rdobase.NewQuerySet().From(table).Limit(1)
		q.Where.And("id", id)

		if rs, err := dcn.Base.Query(q); err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: "Can not pull database instance",
			}
			return
		} else if len(rs) < 1 {
			rsp.Error = &types.ErrorMeta{
				Code:    "404",
				Message: "Node Not Found",
			}
			return
		}

		ft := rdobase.NewFilter()
		ft.And("id", id)

		if _, err = dcn.Base.Update(table, set, ft); err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: fmt.Sprintf("id:%s err:%s", id, err.Error()),
			}
			return
		}
	}

	rsp.Kind = "Node"
}
