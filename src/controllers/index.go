package controllers

import (
	"fmt"
	"strconv"

	"../api"
	"../conf"
	"../datax"

	"github.com/lessos/lessgo/pagelet"
	"github.com/lessos/lessgo/pagelet/vutils"
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

	// fmt.Println(c.Params)
	// for key, val := range c.Params.Values {

	// 	if len(key) > 5 && key[:5] == "term_" {

	// 	}

	// 	// fmt.Println(key, val)
	// }

	c.ViewData["baseuri"] = "/" + specid
	c.ViewData["specid"] = specid

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

	// fmt.Println("c Index dataRender", specid)

	qry := datax.NewQuery(specid, ad.Query.Table)
	if ad.Query.Limit > 0 {
		qry.Limit(ad.Query.Limit)
	}

	if id := c.Params.Get("id"); id != "" {
		if len(id) > 5 && id[len(id)-5:] == ".html" {
			id = id[:len(id)-5]
		}
		qry.Filter("id", id)
	}

	qry.Pager = ad.Pager

	switch ad.Type {

	case "node.list":

		if spec, ok := conf.Instances[specid]; ok {

			for _, modNode := range spec.NodeModels {

				if ad.Query.Table != modNode.Metadata.Name {
					continue
				}

				for _, term := range modNode.Terms {

					if termVal := c.Params.Get("term_" + term.Metadata.Name); termVal != "" {

						switch term.Type {
						case api.TermTaxonomy:
							qry.Filter("term_"+term.Metadata.Name, termVal)
							c.ViewData["term_"+term.Metadata.Name] = termVal
						case api.TermTag:
							// TOPO
							qry.Filter("term_"+term.Metadata.Name+".like", "%"+termVal+"%")
							c.ViewData["term_"+term.Metadata.Name] = termVal
						}
					}
				}

				break
			}
		}

		page := 1
		if c.Params.Get("page") != "" {
			page, _ = strconv.Atoi(c.Params.Get("page"))
			if page > 1 {
				qry.Offset(ad.Query.Limit * (int64(page) - 1))
			}
		}

		if c.Params.Get("qry_text") != "" {
			qry.Filter("title.like", "%"+c.Params.Get("qry_text")+"%")
			c.ViewData["qry_text"] = c.Params.Get("qry_text")
		}

		ls := qry.NodeList()

		c.ViewData[ad.Name] = ls

		if qry.Pager {
			pager := vutils.NewPager(0,
				uint64(ls.Metadata.TotalResults),
				uint64(ls.Metadata.ItemsPerList),
				10)
			pager.CurrentPageNumber = uint64(page)
			c.ViewData[ad.Name+"_pager"] = pager
		}

	case "node.entry":

		c.ViewData[ad.Name] = qry.NodeEntry()

	case "term.list":

		ls := qry.TermList()
		c.ViewData[ad.Name] = ls

		if qry.Pager {
			c.ViewData[ad.Name+"_pager"] = vutils.NewPager(0,
				uint64(ls.Metadata.TotalResults),
				uint64(ls.Metadata.ItemsPerList),
				10)
		}

	case "term.entry":

		c.ViewData[ad.Name] = qry.TermEntry()

	}
}
