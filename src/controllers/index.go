package controllers

import (
	"../api"
	"../conf"
	"../datax"
	"fmt"
	"github.com/lessos/lessgo/pagelet"
)

type Index struct {
	*pagelet.Controller
	SpecID string
}

func (c Index) IndexAction() {
	fmt.Println("Index", c.Params.Get("pagelet"))
}

func (c Index) PageletAction() {

	c.AutoRender = false

	var (
		specid     = c.Params.Get("specid")     // Check
		dataAction = c.Params.Get("dataAction") // Check
		Template   = c.Params.Get("template")
	)

	// fmt.Println(specid, dataAction, Template)

	spec, ok := conf.Instances[specid]
	if !ok {
		return
	}

	// fmt.Println(spec.Actions)

	for _, action := range spec.Actions {

		// fmt.Println(action.Name)

		if action.Name != dataAction {
			continue
		}

		// if c.Params.Get("start") != "" {
		// 	// action.Query.Offset =
		// }

		// fmt.Println(action.Name)

		for _, datax := range action.Datax {
			c.dataRender(specid, datax)
		}

		c.Render(specid, Template)

		break
	}
}

func (c Index) dataRender(specid string, ad api.ActionData) {

	fmt.Println("c Index dataRender", specid)

	qry := datax.NewQuery(specid, ad.Query.Table)
	// fmt.Println(qry)
	if ad.Query.Limit > 0 {
		qry.Limit(ad.Query.Limit)
	}

	// if ad.Typep[:5] == "node" {

	// }

	if ad.Type == "node.list" {
		qry.From("nx" + specid + "_" + qry.Table)
		c.ViewData[ad.Name] = qry.Query()
	} else if ad.Type == "node.entry" {
		qry.From("nx" + specid + "_" + qry.Table)
		c.ViewData[ad.Name] = qry.QueryEntry()
	} else if ad.Type == "term.list" {
		qry.From("tx" + specid + "_" + qry.Table)
		c.ViewData[ad.Name] = qry.Query()
	} else if ad.Type == "term.entry" {
		qry.From("tx" + specid + "_" + qry.Table)
		c.ViewData[ad.Name] = qry.QueryEntry()
	}
}
