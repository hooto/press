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

	var rsp api.NodeList

	defer func() {

		c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Out.Header().Set("Content-type", "application/json")

		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dq := datax.NewQuery(c.Params.Get("specid"), c.Params.Get("modelid"))
	dq.Limit(100)

	rsp = dq.NodeList()
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

	dq := datax.NewQuery(c.Params.Get("specid"), c.Params.Get("modelid"))
	dq.Limit(100)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.NodeEntry()
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

				attrs := []api.KeyValue{}

				for _, attr := range valField.Attrs {
					if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
						attrs = append(attrs, api.KeyValue{attr.Key, attr.Value})
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

				attrs := []api.KeyValue{}

				for _, attr := range valField.Attrs {
					if attr.Key == "format" && utilx.ArrayContain(attr.Value, []string{"md", "text", "html"}) {
						attrs = append(attrs, api.KeyValue{attr.Key, attr.Value})
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
