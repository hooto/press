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
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessids/idsapi"

	"../../api"
	"../../modset"
)

type ModSet struct {
	*httpsrv.Controller
}

func (c ModSet) SpecListAction() {

	c.AutoRender = false

	rsp := api.SpecList{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.list") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return
	}

	q := rdobase.NewQuerySet().From("modules").Limit(100)
	rs, err := dcn.Base.Query(q)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return
	}

	for _, v := range rs {

		var entry api.Spec

		if err := v.Field("body").Json(&entry); err == nil {
			entry.SrvName = v.Field("srvname").String()
			rsp.Items = append(rsp.Items, entry)
		}
	}

	rsp.Kind = "SpecList"
}

func (c ModSet) SpecEntryAction() {

	rsp := api.Spec{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("editor.read") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if c.Params.Get("name") == "" {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Object Not Found",
		}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return
	}

	q := rdobase.NewQuerySet().From("modules").Limit(1)
	q.Where.And("name", c.Params.Get("name"))
	rs, err := dcn.Base.Query(q)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return
	}

	if len(rs) < 1 {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Object Not Found",
		}
		return
	}

	if err := rs[0].Field("body").Json(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: err.Error(),
		}
		return
	}

	rsp.SrvName = rs[0].Field("srvname").String()

	rsp.Kind = "Spec"
}

func (c ModSet) SpecInfoSetAction() {

	var set api.Spec

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("editor.write") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument"}
		return
	}

	if set.Meta.Name == "" {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument"}
		return
	}

	if _, err := modset.SpecFetch(set.Meta.Name); err != nil {

		if err := modset.SpecInfoNew(set); err != nil {
			set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
			return
		}

	} else {

		if err := modset.SpecInfoSet(set); err != nil {
			set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
			return
		}
	}

	if seted, err := modset.SpecFetch(set.Meta.Name); err == nil {
		modset.SpecSchemaSync(seted)
	}

	set.Kind = "Spec"
}
