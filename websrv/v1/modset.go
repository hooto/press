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
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"

	"code.hooto.com/hooto/hootopress/api"
	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/modset"
)

type ModSet struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *ModSet) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c ModSet) SpecListAction() {

	c.AutoRender = false

	rsp := api.SpecList{}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.list", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
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

	if !iamclient.SessionAccessAllowed(c.Session, "editor.read", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if c.Params.Get("name") == "" {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Object Not Found",
		}
		return
	}

	name, err := modset.ModNameFilter(c.Params.Get("name"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
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
	q.Where.And("name", name)
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

	if !iamclient.SessionAccessAllowed(c.Session, "editor.write", config.Config.InstanceID) {

		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	set.Meta.Name, err = modset.ModNameFilter(set.Meta.Name)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.SrvName, err = modset.SrvNameFilter(set.SrvName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	if _, err = modset.SpecFetch(set.Meta.Name); err != nil {

		if err = modset.SpecInfoNew(set); err != nil {
			set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
			return
		}

	} else {

		if err = modset.SpecInfoSet(set); err != nil {
			set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
			return
		}
	}

	seted, err := modset.SpecFetch(set.Meta.Name)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	modset.SpecSchemaSync(seted)

	set.Kind = "Spec"
}

func (c ModSet) SpecTermSetAction() {

	var set api.TermModel

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	set.Meta.Name, err = modset.ModelNameFilter(set.Meta.Name)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.ModName, err = modset.ModNameFilter(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.Type = strings.ToLower(set.Type)
	if set.Type != "tag" && set.Type != "taxonomy" {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Invalid Type"}
		return
	}

	_, err = modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "ModName Not Found"}
		return
	}

	err = modset.SpecTermSet(set.ModName, set)

	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	seted, err := modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	modset.SpecSchemaSync(seted)

	set.Kind = "TermModel"
}

func (c ModSet) SpecNodeSetAction() {

	var set api.NodeModel

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	set.Meta.Name, err = modset.ModelNameFilter(set.Meta.Name)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.ModName, err = modset.ModNameFilter(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	_, err = modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "ModName Not Found"}
		return
	}

	err = modset.SpecNodeSet(set.ModName, set)

	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}
	seted, err := modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	modset.SpecSchemaSync(seted)

	set.Kind = "NodeModel"
}

func (c ModSet) SpecActionSetAction() {

	var set api.Action

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	set.Name, err = modset.ModelNameFilter(set.Name)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.ModName, err = modset.ModNameFilter(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	_, err = modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "ModName Not Found"}
		return
	}

	err = modset.SpecActionSet(set.ModName, set)

	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	seted, err := modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	modset.SpecSchemaSync(seted)

	set.Kind = "Action"
}

func (c ModSet) SpecRouteSetAction() {

	var set api.Route

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&set)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	set.Path, err = modset.RoutePathFilter(set.Path)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	set.ModName, err = modset.ModNameFilter(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, err.Error()}
		return
	}

	_, err = modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "ModName Not Found"}
		return
	}

	err = modset.SpecRouteSet(set.ModName, set)

	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	seted, err := modset.SpecFetch(set.ModName)
	if err != nil {
		set.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	modset.SpecSchemaSync(seted)

	set.Kind = "SpecRoute"
}
