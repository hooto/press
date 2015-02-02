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

	for _, action := range spec.Actions {

		if action.Name != dataAction {
			continue
		}

		// if c.Params.Get("start") != "" {
		// 	// action.Query.Offset =
		// }

		for _, datax := range action.Datax {
			c.dataRender(specid, datax)
		}

		c.Render(specid, Template)

		break
	}
}

func (c Index) dataRender(specid string, ad api.ActionData) {

	qry := datax.NewQuery(specid, ad.Query.Table)

	if ad.Query.Limit > 0 {
		qry.Limit(ad.Query.Limit)
	}

	if ad.Type == "list" {
		c.ViewData[ad.Name] = qry.Query()
	} else if ad.Type == "entry" {
		c.ViewData[ad.Name] = qry.QueryEntry()
	}
}
