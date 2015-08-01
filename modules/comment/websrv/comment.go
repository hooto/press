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

package websrv

import (
	"strings"

	"github.com/eryx/hcaptcha/captcha"
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"../../../src/conf"
	"../../../src/datax"
)

const (
	nsModName          = "comment"
	errCaptchaNotMatch = "CaptchaNotMatch"
)

type Comment struct {
	*httpsrv.Controller
}

func (c Comment) EmbedAction() {

	c.AutoRender = false

	if c.Params.Get("refer_modname") == "" || c.Params.Get("refer_id") == "" {
		return
	}

	qry := datax.NewQuery(nsModName, "entry")
	qry.Limit(500)
	qry.Order("created asc")
	qry.Filter("field_refer_id", c.Params.Get("refer_id"))
	qry.Filter("field_refer", c.Params.Get("refer_modname")+"."+c.Params.Get("refer_datax_table"))

	c.Data["list"] = qry.NodeList()

	c.Data["new_form_refer_id"] = c.Params.Get("refer_id")
	c.Data["new_form_refer_modname"] = c.Params.Get("refer_modname")
	c.Data["new_form_refer_datax_table"] = c.Params.Get("refer_datax_table")

	c.Data["new_form_author"] = "Guest"

	c.Render(nsModName, "embed.tpl")
}

func (c Comment) SetAction() {

	var set TypeComment

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: err.Error(),
		}
		return
	}

	set.Content = strings.TrimSpace(set.Content)
	set.Author = strings.TrimSpace(set.Author)

	if _, ok := conf.Modules[set.ReferModName]; !ok {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Spec Not Found",
		}
		return
	}

	if set.ReferID == "" || set.ReferDataxTable == "" {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "ReferID or ReferDataxTable Can Not be Null",
		}
		return
	}

	re_title := "Re: "
	prevq := datax.NewQuery(set.ReferModName, set.ReferDataxTable)
	prevq.Filter("id", set.ReferID)
	if rs := prevq.NodeEntry(); rs.Error != nil || rs.Kind != "Node" {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Refer Content Not Found",
		}
		return
	} else {
		re_title += rs.Title
	}

	if set.Content == "" {
		set.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Content Can Not be Null",
		}
		return
	}

	if set.Error = captcha.Verify(set.CaptchaToken, set.CaptchaWord); set.Error != nil {

		set.Error.Code = errCaptchaNotMatch
		set.Error.Message = "Word Verification do not match"

		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		set.Error = &types.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	set.Meta.ID = utils.StringNewRand(16)
	set.Meta.Created = rdobase.TimeNow("datetime")

	//
	item := map[string]interface{}{
		"id":             set.Meta.ID,
		"title":          re_title,
		"state":          1,
		"userid":         utils.StringEncode16("guest", 8),
		"field_refer_id": set.ReferID,
		"field_refer":    set.ReferModName + "." + set.ReferDataxTable,
		"field_author":   set.Author,
		"field_content":  set.Content,
		"created":        set.Meta.Created,
		"updated":        set.Meta.Created,
	}

	if _, err := dcn.Base.Insert("nx"+utils.StringEncode16("comment", 12)+"_entry", item); err != nil {
		set.Error = &types.ErrorMeta{
			Code:    "500",
			Message: err.Error(),
		}
		return
	} else {
		set.Kind = "Comment"
	}
}
