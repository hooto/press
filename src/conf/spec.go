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

func specInitialize() error {

	ids := []string{"general", "c8f0ltxp"}

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

		// Instances[specid] = spec

		if err := specSchemaSync(spec); err != nil {
			return err
		}

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

			var fields []api.NodeField
			vnx.Field("fields").Json(&fields)

			// fmt.Println(fields)

			nodeModels = append(nodeModels, api.NodeModel{
				Metadata: api.ObjectMeta{
					ID:      vnx.Field("id").String(),
					Name:    vnx.Field("name").String(),
					UserID:  vnx.Field("userid").String(),
					Created: vnx.Field("created").TimeFormat("datetime", "atom"),
					Updated: vnx.Field("updated").TimeFormat("datetime", "atom"),
				},
				Title:  vnx.Field("title").String(),
				Fields: fields,
			})
		}
	}

	termModels := []api.TermModel{}
	qtx := rdobase.NewQuerySet().From("nodex").Limit(100)
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
			})
		}
	}

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
	}

	// fmt.Println(nodeModels, termModels)
}

func specSchemaSync(spec api.Spec) error {

	var (
		ds                rdobase.DataSet
		nodeModelTemplate rdobase.Table
		termModelTemplate rdobase.Table
		timenow           = rdobase.TimeNow("datetime")
	)

	if err := utils.JsonDecode(dsNodeModels, &nodeModelTemplate); err != nil {
		return err
	}

	if err := utils.JsonDecode(dsTermModels, &termModelTemplate); err != nil {
		return err
	}

	//
	dc, err := rdo.ClientPull("def")
	if err != nil {
		return err
	}

	//
	for _, nodeModel := range spec.NodeModels {

		tbl := nodeModelTemplate
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
			case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
				tbl.AddColumn(&rdobase.Column{
					Name: "field_" + field.Name,
					Type: field.Type,
				})
			}
		}

		ds.Tables = append(ds.Tables, &tbl)

		fieldjs, _ := utils.JsonEncode(nodeModel.Fields)

		dc.Base.InsertIgnore("nodex", map[string]interface{}{
			"id":      spec.Metadata.ID + "." + nodeModel.Metadata.Name,
			"name":    nodeModel.Metadata.Name,
			"specid":  spec.Metadata.ID,
			"userid":  "sysadmin",
			"state":   1,
			"title":   nodeModel.Title,
			"comment": nodeModel.Comment,
			"fields":  fieldjs,
			"created": timenow,
			"updated": timenow,
		})
	}

	for _, termModel := range spec.TermModels {

		tbl := termModelTemplate
		tbl.Name = fmt.Sprintf("tx%s_%s", spec.Metadata.ID, termModel.Metadata.Name)

		ds.Tables = append(ds.Tables, &tbl)

		dc.Base.InsertIgnore("termx", map[string]interface{}{
			"id":      spec.Metadata.ID + "." + termModel.Metadata.Name,
			"name":    termModel.Metadata.Name,
			"specid":  spec.Metadata.ID,
			"userid":  "sysadmin",
			"type":    termModel.Type,
			"state":   1,
			"title":   termModel.Metadata.Name,
			"created": timenow,
			"updated": timenow,
		})
	}

	//
	if err := dc.Dialect.SchemaSync(Config.Database.Dbname, ds); err != nil {
		return err
	}

	dc.Base.InsertIgnore("spec", map[string]interface{}{
		"id":      spec.Metadata.ID,
		"userid":  "sysadmin",
		"state":   1,
		"title":   spec.Metadata.Name,
		"comment": spec.Comment,
		"created": timenow,
		"updated": timenow,
	})

	return nil
}
