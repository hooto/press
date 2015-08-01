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
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessids/idsapi"

	"../../api"
	"../../conf"
	"../../datax"
)

var (
	spaceReg = regexp.MustCompile(" +")
)

type Term struct {
	*httpsrv.Controller
}

func (c Term) ListAction() {

	var rsp api.TermList

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.list") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)

	rsp = dq.TermList()
}

func (c Term) EntryAction() {

	rsp := api.Term{
		TypeMeta: types.TypeMeta{
			APIVersion: api.Version,
		},
	}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.read") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.TermEntry()
}

func (c Term) SetAction() {

	rsp := api.Term{}

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

	model, err := conf.SpecTermModel(c.Params.Get("modname"), c.Params.Get("modelid"))
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
			Message: "Bad Request " + err.Error(),
		}
		return
	}

	var (
		set   = map[string]interface{}{}
		table = fmt.Sprintf("tx%s_%s", utils.StringEncode16(c.Params.Get("modname"), 12), c.Params.Get("modelid"))
	)

	q := rdobase.NewQuerySet().From(table).Limit(1)

	switch model.Type {

	case api.TermTag:

		uniTitle := spaceReg.ReplaceAllString(strings.TrimSpace(strings.ToLower(rsp.Title)), " ")

		h := md5.New()
		io.WriteString(h, uniTitle)
		rsp.UID = fmt.Sprintf("%x", h.Sum(nil))[:16]
		rsp.ID = 0

		q.Where.And("uid", rsp.UID)

		rs, err := dcn.Base.Query(q)
		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: "Can not pull database instance",
			}
			return
		}

		if len(rs) == 1 {

			rsp.ID = rs[0].Field("id").Uint32()

			if rs[0].Field("title").String() != rsp.Title {
				set["title"] = rsp.Title
			}

			if rs[0].Field("state").Int16() != rsp.State {
				set["state"] = rsp.State
			}

		} else {

			set["uid"] = rsp.UID
			set["title"] = rsp.Title
			set["state"] = rsp.State
			set["created"] = rdobase.TimeNow("datetime")
			set["userid"] = "dr5a8pgv"
		}

	case api.TermTaxonomy:

		if rsp.ID > 0 {

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
					Message: "Term Not Found",
				}
				return
			}

			if rs[0].Field("title").String() != rsp.Title {
				set["title"] = rsp.Title
			}

			if rs[0].Field("state").Int16() != rsp.State {
				set["state"] = rsp.State
			}

			if rs[0].Field("pid").Uint32() != rsp.PID {
				set["pid"] = rsp.PID
			}

			if rs[0].Field("weight").Int32() != rsp.Weight {
				set["weight"] = rsp.Weight
			}

		} else {

			set["pid"] = rsp.PID
			set["title"] = rsp.Title
			set["state"] = rsp.State
			set["weight"] = rsp.Weight
			set["created"] = rdobase.TimeNow("datetime")
			set["userid"] = "dr5a8pgv"
		}

	default:
		rsp.Error = &types.ErrorMeta{
			Code:    "500",
			Message: "Server Error",
		}
		return
	}

	if len(set) > 0 {

		set["updated"] = rdobase.TimeNow("datetime")

		if rsp.ID > 0 {

			ft := rdobase.NewFilter()
			ft.And("id", rsp.ID)
			_, err = dcn.Base.Update(table, set, ft)

		} else {

			rs, err := dcn.Base.Insert(table, set)
			if err == nil {
				if incrid, err := rs.LastInsertId(); err == nil && incrid > 0 {
					rsp.ID = uint32(incrid)
				} else {
					err = errors.New("Can Not Get LastInsertId")
				}
			}
		}

		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: err.Error(),
			}
			return
		}
	}

	rsp.Model = model

	rsp.Kind = "Term"
}
