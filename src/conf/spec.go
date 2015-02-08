package conf

import (
	"../api"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/pagelet"
	"github.com/lessos/lessgo/utils"
)

var (
	Instances = map[string]api.Spec{}
)

func SpecNodeModel(specid, modelName string) (*api.NodeModel, error) {

	if spec, ok := Instances[specid]; ok {
		for _, nodeModel := range spec.NodeModels {
			if modelName == nodeModel.Metadata.Name {
				return &nodeModel, nil
			}
		}
	}

	return &api.NodeModel{}, errors.New("Spec Not Found")
}

func SpecTermModel(specid, modelName string) (*api.TermModel, error) {

	if spec, ok := Instances[specid]; ok {
		for _, termModel := range spec.TermModels {
			if modelName == termModel.Metadata.Name {
				return &termModel, nil
			}
		}
	}

	return &api.TermModel{}, errors.New("Spec Not Found")
}

func specInitialize() error {

	ids := []string{"general", "c8f0ltxp"}
	timenow := rdobase.TimeNow("datetime")

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return err
	}

	for _, specid := range ids {

		//
		file := fmt.Sprintf("%s/spec/%s/spec.json", Config.Prefix, specid)
		if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
			return errors.New("Error: config file is not exists")
		}

		fp, err := os.Open(file)
		if err != nil {
			return errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
		}
		defer fp.Close()

		cfgstr, err := ioutil.ReadAll(fp)
		if err != nil {
			return errors.New(fmt.Sprintf("Error: Can not read (%s)", file))
		}

		var spec api.Spec
		if err := utils.JsonDecode(cfgstr, &spec); err != nil {
			return err
		}

		//
		if err := specSchemaSync(spec); err != nil {
			return err
		}

		//
		pagelet.Config.ViewPath(spec.Metadata.ID,
			fmt.Sprintf("%s/spec/%s/views", Config.Prefix, spec.Metadata.ID))

		for _, route := range spec.Router.Routes {

			pagelet.Config.RouteAppend(spec.Metadata.ID, route.Path, map[string]string{
				"specid":     spec.Metadata.ID,
				"dataAction": route.DataAction,
				"template":   route.Template,
				"controller": "Index",
				"action":     "Pagelet",
			})
		}

		//
		set := map[string]interface{}{
			"userid":  "sysadmin",
			"state":   1,
			"title":   spec.Metadata.Name,
			"comment": spec.Comment,
			"version": spec.Metadata.ResourceVersion,
			"updated": timenow,
			"body":    cfgstr,
		}

		q := rdobase.NewQuerySet().From("spec")
		q.Where.And("id", spec.Metadata.ID)

		if rs1, err := dcn.Base.Fetch(q); err == nil {

			if spec.Metadata.Name != rs1.Field("title").String() ||
				spec.Comment != rs1.Field("comment").String() ||
				spec.Metadata.ResourceVersion != rs1.Field("version").String() {

				fr := rdobase.NewFilter()
				fr.And("id", spec.Metadata.ID)

				dcn.Base.Update("spec", set, fr)
			}

		} else {

			set["id"] = spec.Metadata.ID
			set["created"] = timenow

			dcn.Base.Insert("spec", set)
		}

	}

	// Refresh Spec Config
	for _, specid := range ids {
		SpecRefresh(specid)
	}

	return nil
}

func SpecRefresh(specid string) {
	//
	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := rdobase.NewQuerySet().From("spec").Limit(1)
	q.Where.And("id", specid)
	rs, err := dcn.Base.Query(q)
	if err != nil {
		return
	}

	if len(rs) < 1 {
		return
	}

	nodeModels := []api.NodeModel{}
	qnx := rdobase.NewQuerySet().From("nodex").Limit(100)
	qnx.Where.And("specid", specid)
	if rsnx, err := dcn.Base.Query(qnx); err == nil {
		for _, vnx := range rsnx {

			var fields []api.FieldModel
			vnx.Field("fields").Json(&fields)

			var terms []api.TermModel
			vnx.Field("terms").Json(&terms)

			nodeModels = append(nodeModels, api.NodeModel{
				Metadata: api.ObjectMeta{
					ID:      vnx.Field("id").String(),
					Name:    vnx.Field("name").String(),
					UserID:  vnx.Field("userid").String(),
					Created: vnx.Field("created").TimeFormat("datetime", "atom"),
					Updated: vnx.Field("updated").TimeFormat("datetime", "atom"),
				},
				SpecID: vnx.Field("specid").String(),
				Title:  vnx.Field("title").String(),
				Fields: fields,
				Terms:  terms,
			})
		}
	}

	termModels := []api.TermModel{}
	qtx := rdobase.NewQuerySet().From("termx").Limit(500)
	qtx.Where.And("specid", specid)
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
				Type:  vtx.Field("type").String(),
			})
		}
	}

	var specBody api.Spec
	rs[0].Field("body").Json(&specBody)

	Instances[rs[0].Field("id").String()] = api.Spec{
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
		Actions:    specBody.Actions,
	}

	// fmt.Println(nodeModels, termModels)
}

func specSchemaSync(spec api.Spec) error {

	var (
		ds      rdobase.DataSet
		timenow = rdobase.TimeNow("datetime")
	)

	//
	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return err
	}

	//
	for _, nodeModel := range spec.NodeModels {

		var tbl rdobase.Table

		if err := utils.JsonDecode(dsTplNodeModels, &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("nx%s_%s", spec.Metadata.ID, nodeModel.Metadata.Name)

		for _, field := range nodeModel.Fields {

			switch field.Type {

			case "string":

				tbl.AddColumn(&rdobase.Column{
					Name:   "field_" + field.Name,
					Type:   "string",
					Length: field.Length,
				})

			case "text":

				tbl.AddColumn(&rdobase.Column{
					Name: "field_" + field.Name,
					Type: "string-text",
				})

				tbl.AddColumn(&rdobase.Column{
					Name:   "field_" + field.Name + "_attrs",
					Type:   "string",
					Length: "200",
				})

			case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":

				tbl.AddColumn(&rdobase.Column{
					Name: "field_" + field.Name,
					Type: field.Type,
				})

			}
		}

		for _, term := range nodeModel.Terms {

			switch term.Type {

			case api.TermTag:

				tbl.AddColumn(&rdobase.Column{
					Name:   "term_" + term.Metadata.Name,
					Type:   "string",
					Length: "200",
				})

				tbl.AddColumn(&rdobase.Column{
					Name: "term_" + term.Metadata.Name + "_body",
					Type: "string-text",
				})

				tbl.AddColumn(&rdobase.Column{
					Name:   "term_" + term.Metadata.Name + "_idx",
					Type:   "string",
					Length: "100",
				})

				tbl.AddIndex(&rdobase.Index{
					Name: "term_" + term.Metadata.Name + "_idx",
					Type: rdobase.IndexTypeIndex,
					Cols: []string{"term_" + term.Metadata.Name + "_idx"},
				})

			case api.TermTaxonomy:

				tbl.AddColumn(&rdobase.Column{
					Name: "term_" + term.Metadata.Name,
					Type: "uint32",
				})

				tbl.AddIndex(&rdobase.Index{
					Name: "term_" + term.Metadata.Name,
					Type: rdobase.IndexTypeIndex,
					Cols: []string{"term_" + term.Metadata.Name},
				})
			}

		}

		ds.Tables = append(ds.Tables, &tbl)

		fieldsjs, _ := utils.JsonEncode(nodeModel.Fields)
		termsjs, _ := utils.JsonEncode(nodeModel.Terms)

		setID := spec.Metadata.ID + "." + nodeModel.Metadata.Name

		set := map[string]interface{}{
			"name":    nodeModel.Metadata.Name,
			"specid":  spec.Metadata.ID,
			"userid":  "sysadmin",
			"state":   1,
			"title":   nodeModel.Title,
			"comment": nodeModel.Comment,
			"fields":  fieldsjs,
			"terms":   termsjs,
			"updated": timenow,
		}

		q := rdobase.NewQuerySet().From("nodex")
		q.Where.And("id", setID)

		if rs1, err := dcn.Base.Fetch(q); err == nil {

			if nodeModel.Title != rs1.Field("title").String() ||
				nodeModel.Comment != rs1.Field("comment").String() ||
				fieldsjs != rs1.Field("fields").String() ||
				termsjs != rs1.Field("terms").String() {

				fr := rdobase.NewFilter()
				fr.And("id", setID)

				dcn.Base.Update("nodex", set, fr)
			}
		} else {

			set["id"] = setID
			set["created"] = timenow

			dcn.Base.Insert("nodex", set)
		}
	}

	for _, termModel := range spec.TermModels {

		var tbl rdobase.Table

		if err := utils.JsonDecode(dsTplTermModels, &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("tx%s_%s", spec.Metadata.ID, termModel.Metadata.Name)

		switch termModel.Type {

		case api.TermTag:

			tbl.AddColumn(&rdobase.Column{
				Name:   "uid",
				Type:   "string",
				Length: "16",
			})

			tbl.AddIndex(&rdobase.Index{
				Name: "uid",
				Type: rdobase.IndexTypeUnique,
				Cols: []string{"uid"},
			})

		case api.TermTaxonomy:

			tbl.AddColumn(&rdobase.Column{
				Name: "pid",
				Type: "uint32",
			})

			tbl.AddIndex(&rdobase.Index{
				Name: "pid",
				Type: rdobase.IndexTypeIndex,
				Cols: []string{"pid"},
			})

			tbl.AddColumn(&rdobase.Column{
				Name: "weight",
				Type: "int16",
			})

			tbl.AddIndex(&rdobase.Index{
				Name: "weight",
				Type: rdobase.IndexTypeIndex,
				Cols: []string{"weight"},
			})

		default:
			continue
		}

		ds.Tables = append(ds.Tables, &tbl)

		// dcn.Base.InsertIgnore("termx", map[string]interface{}{
		// 	"id":      spec.Metadata.ID + "." + termModel.Metadata.Name,
		// 	"name":    termModel.Metadata.Name,
		// 	"specid":  spec.Metadata.ID,
		// 	"userid":  "sysadmin",
		// 	"type":    termModel.Type,
		// 	"state":   1,
		// 	"title":   termModel.Title,
		// 	"created": timenow,
		// 	"updated": timenow,
		// })

		setID := spec.Metadata.ID + "." + termModel.Metadata.Name

		set := map[string]interface{}{
			"name":    termModel.Metadata.Name,
			"specid":  spec.Metadata.ID,
			"userid":  "sysadmin",
			"type":    termModel.Type,
			"state":   1,
			"title":   termModel.Title,
			"updated": timenow,
		}

		q := rdobase.NewQuerySet().From("termx")
		q.Where.And("id", setID)

		if rs1, err := dcn.Base.Fetch(q); err == nil {

			if termModel.Type != rs1.Field("type").String() ||
				termModel.Title != rs1.Field("title").String() {

				fr := rdobase.NewFilter()
				fr.And("id", setID)

				dcn.Base.Update("termx", set, fr)
			}

		} else {

			set["id"] = setID
			set["created"] = timenow

			dcn.Base.Insert("termx", set)
		}
	}

	//
	if err := dcn.Dialect.SchemaSync(Config.Database.Dbname, ds); err != nil {
		return err
	}

	return nil
}
