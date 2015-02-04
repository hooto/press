package v1

import (
	"../../api"
	"../../conf"

	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/pagelet"
	"github.com/lessos/lessgo/utils"
)

var (
	spaceReg = regexp.MustCompile(" +")
)

type Term struct {
	*pagelet.Controller
}

func (c Term) ListAction() {

	c.AutoRender = false

	rsp := api.TermList{
		TypeMeta: api.TypeMeta{
			APIVersion: api.Version,
		},
	}

	defer func() {

		c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Set("Content-type", "application/json")

		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	model, err := conf.SpecTermModel(c.Params.Get("specid"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return
	}

	table := fmt.Sprintf("tx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))

	q := rdobase.NewQuerySet().From(table).Limit(100)

	if model.Type == api.TermTag {
		q.Order("updated desc")
	} else if model.Type == api.TermTaxonomy {
		q.Order("weight asc")
	}

	rs, err := dcn.Base.Query(q)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
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
}

func (c Term) EntryAction() {

	c.AutoRender = false

	rsp := api.Term{
		TypeMeta: api.TypeMeta{
			APIVersion: api.Version,
		},
	}

	defer func() {

		c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Set("Content-type", "application/json")

		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	table := fmt.Sprintf("tx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))

	q := rdobase.NewQuerySet().From(table).Limit(1)
	q.Where.And("id", c.Params.Get("id"))
	rs, err := dcn.Base.Query(q)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	if len(rs) < 1 {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return
	}

	rsp.Model, err = conf.SpecTermModel(c.Params.Get("specid"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Term Not Found",
		}
		return
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
		return
	}

	rsp.ID = rs[0].Field("id").Uint32()
	rsp.State = rs[0].Field("state").Int16()
	rsp.UserID = rs[0].Field("userid").String()
	rsp.Title = rs[0].Field("title").String()
	rsp.Created = rs[0].Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs[0].Field("updated").TimeFormat("datetime", "atom")

	rsp.Kind = "Term"
}

func (c Term) SetAction() {

	c.AutoRender = false

	rsp := api.Term{
		TypeMeta: api.TypeMeta{
			APIVersion: api.Version,
		},
	}

	defer func() {

		c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Set("Content-type", "application/json")

		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	model, err := conf.SpecTermModel(c.Params.Get("specid"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Spec or Model Not Found",
		}
		return
	}

	if err := utils.JsonDecode(c.Request.RawBody, &rsp); err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "400",
			Message: "Bad Request",
		}
		return
	}

	var (
		set   = map[string]interface{}{}
		table = fmt.Sprintf("tx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))
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
			rsp.Error = &api.ErrorMeta{
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
				rsp.Error = &api.ErrorMeta{
					Code:    "500",
					Message: "Can not pull database instance",
				}
				return
			}

			if len(rs) < 1 {
				rsp.Error = &api.ErrorMeta{
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
		rsp.Error = &api.ErrorMeta{
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
			rsp.Error = &api.ErrorMeta{
				Code:    "500",
				Message: err.Error(),
			}
			return
		}
	}

	rsp.Model = model

	rsp.Kind = "Term"
}
