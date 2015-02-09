package datax

import (
	"fmt"

	"../api"
	"../conf"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
)

func (q *QuerySet) NodeList() api.NodeList {

	rsp := api.NodeList{}

	model, err := conf.SpecNodeModel(q.SpecID, q.Table)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Spec Not Found",
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

	table := fmt.Sprintf("nx%s_%s", q.SpecID, q.Table)
	// q.Limit(100)

	// q := rdobase.NewQuerySet().From(table).Limit(100)
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
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
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

	return rsp
}

func (q *QuerySet) NodeEntry() api.Node {

	rsp := api.Node{
		TypeMeta: api.TypeMeta{
			APIVersion: api.Version,
		},
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
			Message: "Can not pull database instance",
		}
		return rsp
	}

	rsp.Model, err = conf.SpecNodeModel(q.SpecID, q.Table)
	if err != nil {
		rsp.Error = &api.ErrorMeta{
			Code:    "404",
			Message: "Node Not Found",
		}
		return rsp
	}

	table := fmt.Sprintf("nx%s_%s", q.SpecID, q.Table)

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
		rsp.Error = &api.ErrorMeta{
			Code:    "500",
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
			Name:  term.Metadata.Name,
			Value: rs.Field("term_" + term.Metadata.Name).String(),
		})
	}

	rsp.Terms = NodeTermQuery(rsp.Model, rsp.Terms)

	rsp.ID = rs.Field("id").String()
	rsp.State = rs.Field("state").Int16()
	rsp.UserID = rs.Field("userid").String()
	rsp.Title = rs.Field("title").String()
	rsp.Created = rs.Field("created").TimeFormat("datetime", "atom")
	rsp.Updated = rs.Field("updated").TimeFormat("datetime", "atom")

	rsp.Kind = "Node"

	return rsp
}
