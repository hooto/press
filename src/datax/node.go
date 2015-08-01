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

package datax

import (
	"fmt"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"../api"
	"../conf"
)

func (q *QuerySet) NodeCount() (int64, error) {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return 0, err
	}

	table := fmt.Sprintf("nx%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	fr := rdobase.NewFilter()
	fr.And("state", 1)

	return dcn.Base.Count(table, fr)
}

func (q *QuerySet) NodeList() api.NodeList {

	rsp := api.NodeList{}

	model, err := conf.SpecNodeModel(q.ModName, q.Table)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Spec Not Found",
		}
		return rsp
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return rsp
	}

	table := fmt.Sprintf("nx%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	qs := rdobase.NewQuerySet().
		Select(q.cols).
		From(table).
		Limit(q.limit).
		Offset(q.offset)

	if q.order != "" {
		qs.Order(q.order)
	} else {
		qs.Order("created desc")
	}

	q.Filter("state", 1)

	qs.Where = q.filter

	rs, err := dcn.Base.Query(qs)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return rsp
	}

	termBufs := map[string][]string{}

	if len(rs) > 0 {

		for _, v := range rs {

			item := api.Node{
				ID:      v.Field("id").String(),
				PID:     v.Field("pid").String(),
				State:   v.Field("state").Int16(),
				UserID:  v.Field("userid").String(),
				Title:   v.Field("title").String(),
				Created: v.Field("created").TimeFormat("datetime", "atom"),
				Updated: v.Field("updated").TimeFormat("datetime", "atom"),
			}

			for _, field := range model.Fields {

				nodeField := api.NodeField{
					Name:  field.Name,
					Value: v.Field("field_" + field.Name).String(),
				}

				if field.Type == "text" &&
					len(v.Field("field_"+field.Name+"_attrs").String()) > 10 {

					var attrs []api.KeyValue
					if err := v.Field("field_" + field.Name + "_attrs").Json(&attrs); err == nil {
						nodeField.Attrs = attrs
					}
				}

				item.Fields = append(item.Fields, nodeField)
			}

			for _, term := range model.Terms {

				termItem := api.NodeTerm{
					Name:  term.Meta.Name,
					Value: v.Field("term_" + term.Meta.Name).String(),
					Type:  term.Type,
				}

				item.Terms = append(item.Terms, termItem)
				if term.Type == api.TermTaxonomy {
					if !utilx.ArrayContain(termItem.Value, termBufs[termItem.Name]) {
						termBufs[termItem.Name] = append(termBufs[termItem.Name], termItem.Value)
					}
				}
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	// Fetch Terms
	termTaxonomy := map[string]api.Term{}
	for _, term := range model.Terms {

		termids, ok := termBufs[term.Meta.Name]
		if !ok || len(termids) < 1 {
			continue
		}
		ids := []interface{}{}
		for _, tv := range termids {
			ids = append(ids, tv)
		}

		switch term.Type {

		case api.TermTaxonomy:

			table := fmt.Sprintf("tx%s_%s", utils.StringEncode16(q.ModName, 12), term.Meta.Name)
			qs := rdobase.NewQuerySet().From(table).Limit(1000)
			qs.Where.And("id.in", ids...)

			if rs, err := dcn.Base.Query(qs); err == nil && len(rs) > 0 {

				for _, v := range rs {
					termTaxonomy[v.Field("id").String()] = api.Term{
						ID:    rs[0].Field("id").Uint32(),
						Title: rs[0].Field("title").String(),
					}
				}
			}
		}
	}

	//
	for k, v := range rsp.Items {

		for tk, tv := range v.Terms {

			if tv.Value == "" {
				continue
			}

			switch tv.Type {

			case api.TermTaxonomy:

				if tvs, ok := termTaxonomy[tv.Value]; ok {
					rsp.Items[k].Terms[tk].Items = append(rsp.Items[k].Terms[tk].Items, tvs)
				}

			case api.TermTag:

				tags := strings.Split(tv.Value, ",")

				for _, vtag := range tags {

					rsp.Items[k].Terms[tk].Items = append(rsp.Items[k].Terms[tk].Items, api.Term{
						Title: vtag,
					})
				}
			}
		}
	}

	rsp.Model = model

	rsp.Kind = "NodeList"

	if q.Pager {
		num, _ := dcn.Base.Count(table, q.filter)
		rsp.Meta.TotalResults = uint64(num)
		rsp.Meta.StartIndex = uint64(q.offset)
		rsp.Meta.ItemsPerList = uint64(q.limit)
	}

	return rsp
}

func (q *QuerySet) NodeEntry() api.Node {

	rsp := api.Node{}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return rsp
	}

	rsp.Model, err = conf.SpecNodeModel(q.ModName, q.Table)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Node Not Found",
		}
		return rsp
	}

	table := fmt.Sprintf("nx%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	qs := rdobase.NewQuerySet().
		Select(q.cols).
		From(table).
		Order(q.order).
		Limit(1).
		Offset(q.offset)

	qs.Where = q.filter
	// qs.Where.And("id", c.Params.Get("id"))

	rs, err := dcn.Base.Fetch(qs)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: err.Error(),
		}
		return rsp
	}

	for _, field := range rsp.Model.Fields {

		nodeField := api.NodeField{
			Name:  field.Name,
			Value: rs.Field("field_" + field.Name).String(),
		}

		if field.Type == "text" &&
			len(rs.Field("field_"+field.Name+"_attrs").String()) > 10 {

			var attrs []api.KeyValue
			if err := rs.Field("field_" + field.Name + "_attrs").Json(&attrs); err == nil {
				nodeField.Attrs = attrs
			}
		}

		rsp.Fields = append(rsp.Fields, nodeField)
	}

	for _, term := range rsp.Model.Terms {

		rsp.Terms = append(rsp.Terms, api.NodeTerm{
			Name:  term.Meta.Name,
			Value: rs.Field("term_" + term.Meta.Name).String(),
			Type:  term.Type,
		})
	}

	rsp.Terms = NodeTermQuery(q.ModName, rsp.Model, rsp.Terms)

	rsp.ID = rs.Field("id").String()
	rsp.State = rs.Field("state").Int16()
	rsp.UserID = rs.Field("userid").String()
	rsp.Title = rs.Field("title").String()
	rsp.Created = rs.Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs.Field("updated").TimeFormat("datetime", "atom")

	rsp.Kind = "Node"

	return rsp
}
