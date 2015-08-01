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

package modset

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"../api"
	"../conf"
)

func SpecNew(entry api.Spec) error {
	return nil
}

func SpecFetch(modname string) (api.Spec, error) {

	var entry api.Spec

	file := fmt.Sprintf("%s/modules/%s/spec.json", conf.Config.Prefix, modname)
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return entry, errors.New("Error: config file is not exists")
	}

	fp, err := os.Open(file)
	if err != nil {
		return entry, errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	str, err := ioutil.ReadAll(fp)
	if err != nil {
		return entry, err
	}

	if err := utils.JsonDecode(str, &entry); err != nil {
		return entry, err
	}

	return entry, nil
}

func SpecInfoNew(entry api.Spec) error {

	if entry.Meta.Name == "" {
		return errors.New("Name Not Found")
	}

	if entry.Title == "" {
		return errors.New("Title Not Found")
	}

	_, err := SpecFetch(entry.Meta.Name)
	if err == nil {
		return errors.New("Spec Already Exists ")
	}

	entry.Meta.ResourceVersion = "1"
	entry.Meta.Created = utilx.TimeNow("atom")

	dir := fmt.Sprintf("%s/modules/%s", conf.Config.Prefix, entry.Meta.Name)
	if err := os.Mkdir(dir, 0750); err != nil {
		return err
	}

	return _specSet(entry)
}

func SpecInfoSet(entry api.Spec) error {

	if entry.Meta.Name == "" {
		return errors.New("Name Not Found")
	}

	if entry.Title == "" {
		return errors.New("Title Not Found")
	}

	prev, err := SpecFetch(entry.Meta.Name)
	if err != nil {
		return err
	}

	if prev.Title != entry.Title {

		prev.Title = entry.Title
		prev.Meta.Updated = utilx.TimeNow("atom")

		if err := _specSet(prev); err != nil {
			return err
		}
	}

	return err
}

func _specSet(entry api.Spec) error {

	jsb, _ := utils.JsonEncodeIndent(entry, "  ")

	//
	file := fmt.Sprintf("%s/modules/%s/spec.json", conf.Config.Prefix, entry.Meta.Name)

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(int64(len(jsb)))

	_, err = fp.Write(jsb)

	return err
}

func SpecSchemaSync(spec api.Spec) error {

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
	if spec.SrvName == "" {
		spec.SrvName = spec.Meta.Name
	}

	jsb, _ := utils.JsonEncodeIndent(spec, "  ")
	set := map[string]interface{}{
		"srvname": spec.SrvName,
		"status":  1,
		"title":   spec.Title,
		"version": spec.Meta.ResourceVersion,
		"updated": timenow,
		"body":    string(jsb),
	}

	q := rdobase.NewQuerySet().From("modules")
	q.Where.And("name", spec.Meta.Name)

	if _, err := dcn.Base.Fetch(q); err == nil {

		fr := rdobase.NewFilter()
		fr.And("name", spec.Meta.Name)

		_, err = dcn.Base.Update("modules", set, fr)

	} else {

		set["name"] = spec.Meta.Name
		set["created"] = timenow

		_, err = dcn.Base.Insert("modules", set)
	}

	if err != nil {
		return err
	}

	//
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

	//
	if err := dcn.Dialect.SchemaSync(conf.Config.Database.Dbname, ds); err != nil {
		return err
	}

	return nil
}
