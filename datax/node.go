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
	"sync"
	"time"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"code.hooto.com/hooto/alphapress/api"
	"code.hooto.com/hooto/alphapress/config"
	"code.hooto.com/hooto/alphapress/store"
)

var (
	worker_counter_locker sync.Mutex
	worker_pending        = false
)

func Worker() {

	worker_counter_locker.Lock()
	defer worker_counter_locker.Unlock()

	if worker_pending {
		return
	}

	worker_pending = true

	go func() {

		limit := 1000

		for {

			time.Sleep(60e9)
			if store.CacheDB == nil {
				continue
			}

			dcn, err := rdo.ClientPull("def")
			if err != nil {
				continue
			}

			for {

				ls := store.CacheDB.KvScan([]byte("access_counter"), []byte{}, uint32(limit)).Hash()

				imap := map[string]int{}

				for _, v := range ls {

					s := strings.Split(string(v.Key), "/")

					if len(s) == 4 {

						key := s[1] + "/" + s[3]
						if _, ok := imap[key]; ok {
							imap[key]++
						} else {
							imap[key] = 1
						}
					}

					store.CacheDB.KvDel(v.Key)
				}

				for key, num := range imap {

					ks := strings.Split(key, "/")

					if len(ks) != 2 {
						continue
					}

					q := rdobase.NewQuerySet().From(ks[0]).Limit(1)
					q.Where.And("id", ks[1])

					if rs, err := dcn.Base.Query(q); err == nil && len(rs) > 0 {

						ft := rdobase.NewFilter()
						ft.And("id", ks[1])

						dcn.Base.Update(ks[0], map[string]interface{}{
							"ext_access_counter": rs[0].Field("ext_access_counter").Int() + num,
						}, ft)
					}
				}

				if len(ls) < limit {
					break
				}
			}
		}
	}()
}

func (q *QuerySet) NodeCount() (int64, error) {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return 0, err
	}

	table := fmt.Sprintf("nx%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	return dcn.Base.Count(table, q.filter)
}

func (q *QuerySet) NodeList() api.NodeList {

	rsp := api.NodeList{}

	model, err := config.SpecNodeModel(q.ModName, q.Table)
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
	termTaxonomy := map[string]api.Term{}

	if len(rs) > 0 {

		for _, v := range rs {

			item := api.Node{
				ID:      v.Field("id").String(),
				PID:     v.Field("pid").String(),
				Status:  v.Field("status").Int16(),
				UserID:  v.Field("userid").String(),
				Title:   v.Field("title").String(),
				Created: v.Field("created").TimeFormat("datetime", "atom"),
				Updated: v.Field("updated").TimeFormat("datetime", "atom"),
			}

			if model.Extensions.AccessCounter {
				item.ExtAccessCounter = v.Field("ext_access_counter").Uint32()
			}

			if model.Extensions.CommentPerEntry {
				item.ExtCommentPerEntry = v.Field("ext_comment_perentry").Bool()
			}

			if model.Extensions.Permalink != "" && v.Field("ext_permalink_name").String() != "" {
				item.ExtPermalinkName = v.Field("ext_permalink_name").String()
				item.SelfLink = fmt.Sprintf("%s", item.ExtPermalinkName)
			} else {
				item.SelfLink = fmt.Sprintf("%s.html", item.ID)
			}

			for _, field := range model.Fields {

				nodeField := api.NodeField{
					Name:        field.Name,
					Value:       v.Field("field_" + field.Name).String(),
					ValueCaches: map[string]string{},
				}

				if field.Type == "text" &&
					len(v.Field("field_"+field.Name+"_attrs").String()) > 10 {

					var attrs []api.KeyValue
					if err := v.Field("field_" + field.Name + "_attrs").Json(&attrs); err == nil {
						nodeField.Attrs = attrs
					}
				}

				item.Fields = append(item.Fields, &nodeField)
			}

			for _, term := range model.Terms {

				termItem := api.NodeTerm{
					Name:  term.Meta.Name,
					Value: v.Field("term_" + term.Meta.Name).String(),
					Type:  term.Type,
				}

				item.Terms = append(item.Terms, termItem)
				if term.Type == api.TermTaxonomy {

					if te := TermTaxonomyCacheEntry(q.ModName, term.Meta.Name, v.Field("term_"+term.Meta.Name).Uint32()); te != nil {

						termTaxonomy[v.Field("term_"+term.Meta.Name).String()] = api.Term{
							ID:    te.ID,
							Title: te.Title,
						}

					} else if !utilx.ArrayContain(termItem.Value, termBufs[termItem.Name]) {
						termBufs[termItem.Name] = append(termBufs[termItem.Name], termItem.Value)
					}
				}
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	// Fetch Terms
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
						ID:    v.Field("id").Uint32(),
						Title: v.Field("title").String(),
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

	// qryhash := q.Hash()

	// if rsp.Model.CacheTTL > 0 && entry.Title != "" {
	// 	store.CacheSetJson(qryhash, rsp, rsp.Model.CacheTTL)
	// }

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

	rsp.Model, err = config.SpecNodeModel(q.ModName, q.Table)
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
			Name:        field.Name,
			Value:       rs.Field("field_" + field.Name).String(),
			ValueCaches: map[string]string{},
		}

		if field.Type == "text" &&
			len(rs.Field("field_"+field.Name+"_attrs").String()) > 10 {

			var attrs []api.KeyValue
			if err := rs.Field("field_" + field.Name + "_attrs").Json(&attrs); err == nil {
				nodeField.Attrs = attrs
			}
		}

		rsp.Fields = append(rsp.Fields, &nodeField)
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
	rsp.Status = rs.Field("status").Int16()
	rsp.UserID = rs.Field("userid").String()
	rsp.Title = rs.Field("title").String()
	rsp.Created = rs.Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs.Field("updated").TimeFormat("datetime", "atom")

	if rsp.Model.Extensions.AccessCounter {
		rsp.ExtAccessCounter = rs.Field("ext_access_counter").Uint32()
	}

	if rsp.Model.Extensions.CommentPerEntry {
		rsp.ExtCommentPerEntry = rs.Field("ext_comment_perentry").Bool()
	}

	if rsp.Model.Extensions.Permalink != "" {
		rsp.ExtPermalinkName = rs.Field("ext_permalink_name").String()
	}

	rsp.Kind = "Node"

	// qryhash := q.Hash()

	// if rsp.Model.CacheTTL > 0 && entry.Title != "" {
	// 	store.CacheSetJson(qryhash, rsp, rsp.Model.CacheTTL)
	// }

	return rsp
}
