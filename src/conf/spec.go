// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessgo/utils"

	"../api"
)

var (
	Modules = map[string]*api.Spec{}
)

func SpecNodeModel(modname, modelName string) (*api.NodeModel, error) {

	for _, mod := range Modules {

		if mod.Meta.Name != modname {
			continue
		}

		for _, nodeModel := range mod.NodeModels {

			if modelName == nodeModel.Meta.Name {
				return &nodeModel, nil
			}
		}
	}

	return &api.NodeModel{}, errors.New("Spec Not Found")
}

func SpecTermModel(modname, modelName string) (*api.TermModel, error) {

	for _, mod := range Modules {

		if mod.Meta.Name != modname {
			continue
		}

		for _, termModel := range mod.TermModels {

			if modelName == termModel.Meta.Name {
				return &termModel, nil
			}
		}
	}

	return &api.TermModel{}, errors.New("Spec Not Found")
}

func module_init() error {

	timenow := rdobase.TimeNow("datetime")

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return err
	}

	//
	q := rdobase.NewQuerySet().From("modules").Limit(200)
	q.Where.And("status", 1)
	if rs, err := dcn.Base.Query(q); err == nil {

		for _, v := range rs {

			var mod api.Spec

			if err := v.Field("body").Json(&mod); err == nil && mod.Meta.Name != "" {
				mod.SrvName = v.Field("srvname").String()
				Modules[v.Field("srvname").String()] = &mod
			} else {
				logger.Printf("error", "Module.Init(%s) Failed", v.Field("name").String())
			}
		}
	}

	//
	for _, modname := range coreModules {

		//
		file := fmt.Sprintf("%s/modules/%s/spec.json", Config.Prefix, modname)
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

		specResVersion, _ := strconv.Atoi(spec.Meta.ResourceVersion)
		instResVersion := 0

		for _, mod := range Modules {

			if mod.Meta.Name == modname {

				instResVersion, _ = strconv.Atoi(mod.Meta.ResourceVersion)

				break
			}
		}

		if specResVersion <= instResVersion {
			continue
		}

		//
		jsb, _ := utils.JsonEncodeIndent(spec, "  ")
		set := map[string]interface{}{
			"status":  1,
			"title":   spec.Title,
			"version": spec.Meta.ResourceVersion,
			"updated": timenow,
			"body":    string(jsb),
		}

		q = rdobase.NewQuerySet().From("modules")
		q.Where.And("name", spec.Meta.Name)

		if _, err := dcn.Base.Fetch(q); err == nil {

			fr := rdobase.NewFilter()
			fr.And("name", spec.Meta.Name)

			dcn.Base.Update("modules", set, fr)

		} else {

			set["name"] = spec.Meta.Name
			set["srvname"] = spec.Meta.Name
			set["created"] = timenow

			dcn.Base.Insert("modules", set)
		}

		Modules[spec.Meta.Name] = &spec
	}

	//
	for _, mod := range Modules {

		if err := _instance_schema_sync(mod); err != nil {
			return err
		}

		SpecRefresh(mod.SrvName)
	}

	return nil
}

func SpecRefresh(srvname string) {

	spec, ok := Modules[srvname]
	if !ok {
		return
	}

	for i, v := range spec.Router.Routes {
		spec.Router.Routes[i].Tree = strings.Split(strings.Trim(filepath.Clean(v.Path), "/"), "/")
	}

	httpsrv.GlobalService.TemplateLoader.Set(spec.Meta.Name,
		[]string{fmt.Sprintf("%s/modules/%s/views", Config.Prefix, spec.Meta.Name)})
}

func _instance_schema_sync(spec *api.Spec) error {

	//
	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return err
	}

	ds := rdobase.DataSet{}

	// nodes
	for _, nodeModel := range spec.NodeModels {

		var tbl rdobase.Table

		if err := utils.JsonDecode(dsTplNodeModels, &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("nx%s_%s", utils.StringEncode16(spec.Meta.Name, 12), nodeModel.Meta.Name)

		for _, field := range nodeModel.Fields {

			switch field.Type {

			case "string":

				tbl.AddColumn(&rdobase.Column{
					Name:   "field_" + field.Name,
					Type:   "string",
					Length: field.Length,
				})

				it, _ := strconv.Atoi(field.IndexType)
				switch it {
				case rdobase.IndexTypeUnique, rdobase.IndexTypeIndex:
					tbl.AddIndex(&rdobase.Index{
						Name: "field_" + field.Name,
						Type: it,
						Cols: []string{"field_" + field.Name},
					})
				}

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
					Name:   "term_" + term.Meta.Name,
					Type:   "string",
					Length: "200",
				})

				tbl.AddColumn(&rdobase.Column{
					Name: "term_" + term.Meta.Name + "_body",
					Type: "string-text",
				})

				tbl.AddColumn(&rdobase.Column{
					Name:   "term_" + term.Meta.Name + "_idx",
					Type:   "string",
					Length: "100",
				})

				tbl.AddIndex(&rdobase.Index{
					Name: "term_" + term.Meta.Name + "_idx",
					Type: rdobase.IndexTypeIndex,
					Cols: []string{"term_" + term.Meta.Name + "_idx"},
				})

			case api.TermTaxonomy:

				tbl.AddColumn(&rdobase.Column{
					Name: "term_" + term.Meta.Name,
					Type: "uint32",
				})

				tbl.AddIndex(&rdobase.Index{
					Name: "term_" + term.Meta.Name,
					Type: rdobase.IndexTypeIndex,
					Cols: []string{"term_" + term.Meta.Name},
				})
			}

		}

		ds.Tables = append(ds.Tables, &tbl)
	}

	// terms
	for _, termModel := range spec.TermModels {

		var tbl rdobase.Table

		if err := utils.JsonDecode(dsTplTermModels, &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("tx%s_%s", utils.StringEncode16(spec.Meta.Name, 12), termModel.Meta.Name)

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
	}

	// sync
	if err := dcn.Dialect.SchemaSync(Config.Database.Dbname, ds); err != nil {
		return err
	}

	return nil
}
