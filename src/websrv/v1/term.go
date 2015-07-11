package v1

import (
	"../../api"
	"../../conf"
	"../../datax"

	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/utils"
)

var (
	spaceReg = regexp.MustCompile(" +")
)

type Term struct {
	*httpsrv.Controller
}

func (c Term) ListAction() {

	c.AutoRender = false

	var rsp api.TermList

	defer func() {

		c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Set("Content-type", "application/json")

		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dq := datax.NewQuery(c.Params.Get("specid"), c.Params.Get("modelid"))
	dq.Limit(100)

	rsp = dq.TermList()
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

	dq := datax.NewQuery(c.Params.Get("specid"), c.Params.Get("modelid"))
	dq.Limit(100)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.TermEntry()
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
