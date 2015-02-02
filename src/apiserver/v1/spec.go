package v1

import (
	"../../api"
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/pagelet"
	"github.com/lessos/lessgo/utils"
	"io"
)

type Spec struct {
	*pagelet.Controller
}

func (c Spec) ListAction() {

	c.AutoRender = false

	rsp := api.SpecList{
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

	q := rdobase.NewQuerySet().From("spec").Limit(100)
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

			//
			nodeModels := []api.NodeModel{}
			qnx := rdobase.NewQuerySet().From("nodex").Limit(100)
			qnx.Where.And("specid", v.Field("id").String())
			if rsnx, err := dcn.Base.Query(qnx); err == nil {
				for _, vnx := range rsnx {
					nodeModels = append(nodeModels, api.NodeModel{
						Metadata: api.ObjectMeta{
							ID:      vnx.Field("id").String(),
							Name:    vnx.Field("name").String(),
							UserID:  vnx.Field("userid").String(),
							Created: vnx.Field("created").TimeFormat("datetime", "atom"),
							Updated: vnx.Field("updated").TimeFormat("datetime", "atom"),
						},
						Title: vnx.Field("title").String(),
					})
				}
			}

			//
			termModels := []api.TermModel{}
			qtx := rdobase.NewQuerySet().From("nodex").Limit(100)
			qtx.Where.And("specid", v.Field("id").String())
			if rstx, err := dcn.Base.Query(qtx); err == nil {
				for _, vtx := range rstx {
					termModels = append(termModels, api.TermModel{
						Metadata: api.ObjectMeta{
							ID:      vtx.Field("id").String(),
							Name:    vtx.Field("name").String(),
							UserID:  vtx.Field("userid").String(),
							Created: vtx.Field("created").TimeFormat("datetime", "atom"),
							Updated: vtx.Field("updated").TimeFormat("datetime", "atom"),
						},
						Title: vtx.Field("title").String(),
					})
				}
			}

			//
			rsp.Items = append(rsp.Items, api.Spec{
				Metadata: api.ObjectMeta{
					ID:              v.Field("id").String(),
					UserID:          v.Field("userid").String(),
					Created:         v.Field("created").TimeFormat("datetime", "atom"),
					Updated:         v.Field("updated").TimeFormat("datetime", "atom"),
					ResourceVersion: v.Field("version").String(),
				},
				State:      v.Field("state").Int16(),
				Title:      v.Field("title").String(),
				Comment:    v.Field("comment").String(),
				NodeModels: nodeModels,
				TermModels: termModels,
			})
		}
	}

	rsp.Kind = "SpecList"
}

func (c Spec) EntryAction() {

	c.AutoRender = false

	rsp := api.Spec{
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

	if c.Params.Get("id") == "" {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Channel Not Found",
		}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return
	}

	q := rdobase.NewQuerySet().From("spec").Limit(1)
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
			Message: "Object Not Found",
		}
		return
	}

	nodeModels := []api.NodeModel{}
	qnx := rdobase.NewQuerySet().From("nodex").Limit(100)
	qnx.Where.And("specid", rs[0].Field("id").String())
	if rsnx, err := dcn.Base.Query(qnx); err == nil {
		for _, vnx := range rsnx {
			nodeModels = append(nodeModels, api.NodeModel{
				Metadata: api.ObjectMeta{
					ID:      vnx.Field("id").String(),
					Name:    vnx.Field("name").String(),
					UserID:  vnx.Field("userid").String(),
					Created: vnx.Field("created").TimeFormat("datetime", "atom"),
					Updated: vnx.Field("updated").TimeFormat("datetime", "atom"),
				},
				Title: vnx.Field("title").String(),
			})
		}
	}

	termModels := []api.TermModel{}
	qtx := rdobase.NewQuerySet().From("nodex").Limit(100)
	qtx.Where.And("specid", rs[0].Field("id").String())
	if rstx, err := dcn.Base.Query(qtx); err == nil {
		for _, vtx := range rstx {
			termModels = append(termModels, api.TermModel{
				Metadata: api.ObjectMeta{
					ID:      vtx.Field("id").String(),
					Name:    vtx.Field("name").String(),
					UserID:  vtx.Field("userid").String(),
					Created: vtx.Field("created").TimeFormat("datetime", "atom"),
					Updated: vtx.Field("updated").TimeFormat("datetime", "atom"),
				},
				Title: vtx.Field("title").String(),
			})
		}
	}

	rsp = api.Spec{
		Metadata: api.ObjectMeta{
			ID:              rs[0].Field("id").String(),
			UserID:          rs[0].Field("userid").String(),
			Created:         rs[0].Field("created").TimeFormat("datetime", "atom"),
			Updated:         rs[0].Field("updated").TimeFormat("datetime", "atom"),
			ResourceVersion: rs[0].Field("version").String(),
		},
		State:      rs[0].Field("state").Int16(),
		Title:      rs[0].Field("title").String(),
		Comment:    rs[0].Field("comment").String(),
		NodeModels: nodeModels,
		TermModels: termModels,
	}

	rsp.Kind = "Spec"
}
