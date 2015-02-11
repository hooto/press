package datax

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strings"

	"../api"
	"../conf"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
)

var (
	spaceReg = regexp.MustCompile(" +")
)

func (q *QuerySet) TermList() api.TermList {

	rsp := api.TermList{}

	model, err := conf.SpecTermModel(q.SpecID, q.Table)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return rsp
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
	}

	q.limit = 100
	table := fmt.Sprintf("tx%s_%s", q.SpecID, q.Table)

	qs := rdobase.NewQuerySet().
		Select(q.cols).
		From(table).
		Limit(q.limit).
		Offset(q.offset)

	if model.Type == api.TermTag {
		qs.Order("updated desc")
	} else if model.Type == api.TermTaxonomy {
		qs.Order("weight asc")
	}

	qs.Where = q.filter

	rs, err := dcn.Base.Query(qs)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
	}

	if len(rs) > 0 {

		for _, v := range rs {

			item := api.Term{
				ID:      v.Field("id").Uint32(),
				State:   v.Field("state").Int16(),
				UserID:  v.Field("userid").String(),
				Title:   v.Field("title").String(),
				Created: v.Field("created").TimeFormat("datetime", "atom"),
				Updated: v.Field("updated").TimeFormat("datetime", "atom"),
			}

			switch model.Type {
			case api.TermTag:
				item.UID = v.Field("uid").String()
			case api.TermTaxonomy:
				item.PID = v.Field("pid").Uint32()
				item.Weight = v.Field("weight").Int32()
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	rsp.Model = model
	rsp.Kind = "TermList"

	if q.Pager {
		num, _ := dcn.Base.Count(table, q.filter)
		rsp.Metadata.TotalResults = uint64(num)
		rsp.Metadata.StartIndex = uint64(q.offset)
		rsp.Metadata.ItemsPerList = uint64(q.limit)
	}

	return rsp
}

func (q *QuerySet) TermEntry() api.Term {

	rsp := api.Term{}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
	}

	rsp.Model, err = conf.SpecTermModel(q.SpecID, q.Table)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return rsp
	}

	table := fmt.Sprintf("tx%s_%s", q.SpecID, q.Table)

	qs := rdobase.NewQuerySet().
		Select(q.cols).
		From(table).
		Order(q.order).
		Limit(1).
		Offset(q.offset)

	qs.Where = q.filter

	rs, err := dcn.Base.Query(qs)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
	}

	if len(rs) < 1 {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return rsp
	}

	switch rsp.Model.Type {
	case api.TermTaxonomy:
		rsp.PID = rs[0].Field("pid").Uint32()
		rsp.Weight = rs[0].Field("weight").Int32()
	case api.TermTag:
		rsp.UID = rs[0].Field("uid").String()
	default:
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Server Error",
		}
		return rsp
	}

	rsp.ID = rs[0].Field("id").Uint32()
	rsp.State = rs[0].Field("state").Int16()
	rsp.UserID = rs[0].Field("userid").String()
	rsp.Title = rs[0].Field("title").String()
	rsp.Created = rs[0].Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs[0].Field("updated").TimeFormat("datetime", "atom")

	rsp.Kind = "Term"

	return rsp
}

type TermList api.TermList

func (t *TermList) Index() string {

	if len(t.Items) < 1 {
		return ""
	}

	idxs := []string{}
	for _, v := range t.Items {
		idxs = append(idxs, fmt.Sprintf("%v", v.ID))
	}

	return strings.Join(idxs, ",")
}

func (t *TermList) Content() string {

	if len(t.Items) < 1 {
		return ""
	}

	ts := []string{}
	for _, v := range t.Items {
		ts = append(ts, v.Title)
	}

	return strings.Join(ts, ",")
}

func NodeTermQuery(model *api.NodeModel, terms []api.NodeTerm) []api.NodeTerm {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return terms
	}

	for _, modTerm := range model.Terms {

		for k, term := range terms {

			if modTerm.Metadata.Name != term.Name {
				continue
			}

			table := fmt.Sprintf("tx%s_%s", model.SpecID, modTerm.Metadata.Name)

			q := rdobase.NewQuerySet().From(table)

			switch modTerm.Type {

			case api.TermTag:

				// TODO

			case api.TermTaxonomy:

				q.Limit(1)
				q.Where.And("id", term.Value)

				if rs, err := dcn.Base.Query(q); err == nil && len(rs) > 0 {

					terms[k].Items = append(terms[k].Items, api.Term{
						ID:    rs[0].Field("id").Uint32(),
						Title: rs[0].Field("title").String(),
					})
				}
			}

			break
		}
	}

	return terms
}

func TermSync(specid, modelid, terms string) (TermList, error) {

	ls := TermList{}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return ls, err
	}

	terms = spaceReg.ReplaceAllString(terms, " ")

	tars := strings.Split(terms, ",")

	ids := []interface{}{}

	for _, term := range tars {

		tag := api.Term{
			Title: strings.TrimSpace(term),
		}

		if len(tag.Title) < 1 {
			continue
		}

		h := md5.New()
		io.WriteString(h, strings.ToLower(tag.Title))
		tag.UID = fmt.Sprintf("%x", h.Sum(nil))[:16]

		exist := false
		for _, prev := range ids {
			if prev.(string) == tag.UID {
				exist = true
				break
			}
		}
		if exist {
			continue
		}

		ls.Items = append(ls.Items, tag)

		ids = append(ids, tag.UID)
	}

	table := fmt.Sprintf("tx%s_%s", specid, modelid)

	if len(ids) > 0 {

		q := rdobase.NewQuerySet().From(table).Limit(int64(len(ids)))
		q.Where.And("uid.in", ids...)

		if rs, err := dcn.Base.Query(q); err == nil {
			for _, v := range rs {

				for tk, tv := range ls.Items {
					if v.Field("uid").String() == tv.UID {
						ls.Items[tk].ID = v.Field("id").Uint32()
						break
					}
				}
			}
		}
	}

	timenow := rdobase.TimeNow("datetime")

	for tk, tv := range ls.Items {

		if tv.ID > 0 {
			continue
		}

		rs, err := dcn.Base.Insert(table, map[string]interface{}{
			"uid":     tv.UID,
			"title":   tv.Title,
			"userid":  "sysadmin",
			"state":   1,
			"created": timenow,
			"updated": timenow,
		})

		if err == nil {
			if incrid, err := rs.LastInsertId(); err == nil && incrid > 0 {
				ls.Items[tk].ID = uint32(incrid)
			}
		}
	}

	return ls, nil
}
