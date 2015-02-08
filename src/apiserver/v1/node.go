package v1

import (
	"fmt"
	"io"

	"../../api"
	"../../conf"
	"../../datax"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/pagelet"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"
)

type Node struct {
	*pagelet.Controller
}

func (c Node) ListAction() {

	c.AutoRender = false

	rsp := api.NodeList{
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

	model, err := conf.SpecNodeModel(c.Params.Get("specid"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Spec Not Found",
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

	table := fmt.Sprintf("nx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))

	q := rdobase.NewQuerySet().From(table).Limit(100)
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

			item := api.Node{
				ID:      v.Field("id").String(),
				State:   v.Field("state").Int16(),
				UserID:  v.Field("userid").String(),
				Title:   v.Field("title").String(),
				Created: v.Field("created").TimeFormat("datetime", "atom"),
				Updated: v.Field("updated").TimeFormat("datetime", "atom"),
			}

			for _, field := range model.Fields {

				item.Fields = append(item.Fields, api.NodeField{
					Name:  field.Name,
					Value: v.Field("field_" + field.Name).String(),
				})
			}

			for _, term := range model.Terms {

				item.Terms = append(item.Terms, api.NodeTerm{
					Name:  term.Metadata.Name,
					Value: v.Field("term_" + term.Metadata.Name).String(),
				})
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	rsp.Model = model

	rsp.Kind = "NodeList"
}

func (c Node) EntryAction() {

	c.AutoRender = false

	rsp := api.Node{
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

	table := fmt.Sprintf("nx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))

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
			Message: "Node Not Found",
		}
		return
	}

	rsp.Model, err = conf.SpecNodeModel(c.Params.Get("specid"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Node Not Found",
		}
		return
	}

	for _, field := range rsp.Model.Fields {

		rsp.Fields = append(rsp.Fields, api.NodeField{
			Name:  field.Name,
			Value: rs[0].Field("field_" + field.Name).String(),
		})
	}

	for _, term := range rsp.Model.Terms {

		rsp.Terms = append(rsp.Terms, api.NodeTerm{
			Name:  term.Metadata.Name,
			Value: rs[0].Field("term_" + term.Metadata.Name).String(),
		})
	}

	rsp.Terms = datax.NodeTermsQuery(rsp.Model, rsp.Terms)

	rsp.ID = rs[0].Field("id").String()
	rsp.State = rs[0].Field("state").Int16()
	rsp.UserID = rs[0].Field("userid").String()
	rsp.Title = rs[0].Field("title").String()
	rsp.Created = rs[0].Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs[0].Field("updated").TimeFormat("datetime", "atom")

	rsp.Kind = "Node"
}

func (c Node) SetAction() {

	c.AutoRender = false

	rsp := api.Node{
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

	model, err := conf.SpecNodeModel(c.Params.Get("specid"), c.Params.Get("modelid"))
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
		set = map[string]interface{}{}
		fns = []string{}
	)

	for _, modField := range model.Fields {
		fns = append(fns, modField.Name)
	}

	//
	table := fmt.Sprintf("nx%s_%s", c.Params.Get("specid"), c.Params.Get("modelid"))

	if len(rsp.ID) > 0 {

		q := rdobase.NewQuerySet().From(table).Limit(1)
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
				Message: "Node Not Found",
			}
			return
		}

		if rs[0].Field("title").String() != rsp.Title {
			set["title"] = rsp.Title
		}

		if rs[0].Field("state").Int16() != rsp.State {
			set["state"] = rsp.State
		}

		//
		for _, valField := range rsp.Fields {

			if utilx.ArrayContain(valField.Name, fns) &&
				rs[0].Field("field_"+valField.Name).String() != valField.Value {
				set["field_"+valField.Name] = valField.Value

				attrs := map[string]string{}

				for _, attr := range valField.Attrs {
					if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
						attrs["format"] = attr.Value
					}
				}

				set["field_"+valField.Name+"_attrs"], _ = utils.JsonEncode(attrs)
			}
		}

		//
		for _, modTerm := range model.Terms {

			for _, term := range rsp.Terms {

				if modTerm.Metadata.Name != term.Name {
					continue
				}

				switch modTerm.Type {

				case api.TermTag:

					tags, _ := datax.TermSync(model.SpecID, modTerm.Metadata.Name, term.Value)

					if rs[0].Field("term_"+term.Name).String() != term.Value {
						set["term_"+modTerm.Metadata.Name] = tags.Content()
						set["term_"+modTerm.Metadata.Name+"_idx"] = tags.Index()
					}

				case api.TermTaxonomy:

					set["term_"+modTerm.Metadata.Name] = term.Value
				}
			}
		}

	} else {

		set["id"] = utils.StringNewRand36(12)
		set["title"] = rsp.Title
		set["state"] = rsp.State
		set["created"] = rdobase.TimeNow("datetime")
		set["userid"] = "dr5a8pgv"

		//
		for _, valField := range rsp.Fields {

			if utilx.ArrayContain(valField.Name, fns) {

				set["field_"+valField.Name] = valField.Value

				attrs := map[string]string{}

				for _, attr := range valField.Attrs {
					if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
						attrs["format"] = attr.Value
					}
				}

				set["field_"+valField.Name+"_attrs"], _ = utils.JsonEncode(attrs)
			}
		}

		//
		for _, modTerm := range model.Terms {

			for _, term := range rsp.Terms {

				if modTerm.Metadata.Name != term.Name {
					continue
				}

				switch modTerm.Type {

				case api.TermTag:

					tags, _ := datax.TermSync(model.SpecID, modTerm.Metadata.Name, term.Value)
					set["term_"+modTerm.Metadata.Name] = tags.Content()
					set["term_"+modTerm.Metadata.Name+"_idx"] = tags.Index()

				case api.TermTaxonomy:

					set["term_"+modTerm.Metadata.Name] = term.Value
				}
			}
		}
	}

	if len(set) > 0 {

		set["updated"] = rdobase.TimeNow("datetime")

		if len(rsp.ID) > 0 {

			ft := rdobase.NewFilter()
			ft.And("id", rsp.ID)
			_, err = dcn.Base.Update(table, set, ft)

		} else {
			rsp.ID = set["id"].(string)
			_, err = dcn.Base.Insert(table, set)
		}

		if err != nil {
			rsp.Error = &api.ErrorMeta{
				Code:    "500",
				Message: err.Error(),
			}
			return
		}
	}

	rsp.Kind = "Node"
}
