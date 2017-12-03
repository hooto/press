// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/iomix/rdb/modeler"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/store"
)

var (
	locker  sync.Mutex
	Modules = map[string]*api.Spec{}
)

func SpecSet(spec *api.Spec) {

	locker.Lock()
	defer locker.Unlock()

	if strings.Contains(spec.SrvName, "/") {
		spec.SrvName, _ = api.SrvNameFilter(spec.SrvName)
	}

	Modules[spec.SrvName] = spec
}

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

	timenow := rdb.TimeNow("datetime")

	if store.Data == nil {
		return errors.New("No RDB Connector Found")
	}

	//
	q := rdb.NewQuerySet().From("modules").Limit(200)
	q.Where.And("status", 1)
	if rs, err := store.Data.Query(q); err == nil {

		for _, v := range rs {

			var mod api.Spec

			if err := v.Field("body").JsonDecode(&mod); err == nil && mod.Meta.Name != "" {
				if mod.SrvName == "" || strings.Contains(mod.SrvName, "/") {
					mod.SrvName, _ = api.SrvNameFilter(v.Field("srvname").String())
				}
				Modules[mod.SrvName] = &mod
			} else {
				hlog.Printf("error", "Module.Init(%s) Failed", v.Field("name").String())
			}
		}
	}

	//
	for _, modname := range coreModules {

		//
		file := fmt.Sprintf("%s/modules/%s/spec.json", Prefix, modname)
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
		if err := json.Decode([]byte(cfgstr), &spec); err != nil {
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

		if spec.SrvName == "" || strings.Contains(spec.SrvName, "/") {
			spec.SrvName, err = api.SrvNameFilter(spec.Meta.Name)
			if err != nil {
				return err
			}
		}

		//
		jsb, _ := json.Encode(spec, "  ")
		set := map[string]interface{}{
			"status":  1,
			"title":   spec.Title,
			"version": spec.Meta.ResourceVersion,
			"updated": timenow,
			"body":    string(jsb),
		}

		q = rdb.NewQuerySet().From("modules")
		q.Where.And("name", spec.Meta.Name)

		if _, err := store.Data.Fetch(q); err == nil {

			fr := rdb.NewFilter()
			fr.And("name", spec.Meta.Name)

			store.Data.Update("modules", set, fr)

		} else {

			set["name"] = spec.Meta.Name
			set["srvname"] = spec.SrvName
			set["created"] = timenow

			store.Data.Insert("modules", set)
		}

		Modules[spec.SrvName] = &spec
	}

	//
	for _, mod := range Modules {

		if err := _instance_schema_sync(mod); err != nil {
			return err
		}

		SpecSrvRefresh(mod.SrvName)
	}

	return nil
}

func SpecRefresh(modname string) {

	for srvname, spec := range Modules {

		if spec.Meta.Name == modname {
			SpecSrvRefresh(srvname)
			break
		}
	}
}

func SpecSrvRefresh(srvname string) {

	if strings.Contains(srvname, "/") {
		srvname, _ = api.SrvNameFilter(srvname)
	}

	spec, ok := Modules[srvname]
	if !ok {
		return
	}

	for i, v := range spec.Router.Routes {
		spec.Router.Routes[i].Tree = strings.Split(strings.Trim(filepath.Clean(v.Path), "/"), "/")
	}

	httpsrv.GlobalService.TemplateLoader.Clean(spec.Meta.Name)
	httpsrv.GlobalService.TemplateLoader.Set(spec.Meta.Name,
		[]string{fmt.Sprintf("%s/modules/%s/views", Prefix, spec.Meta.Name)})
}

func _instance_schema_sync(spec *api.Spec) error {

	if store.Data == nil {
		return errors.New("No RDB Connector Found")
	}

	ds := modeler.DatabaseEntry{}

	// nodes
	for _, nodeModel := range spec.NodeModels {

		var tbl modeler.Table

		if err := json.Decode([]byte(dsTplNodeModels), &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("nx%s_%s", idhash.HashToHexString([]byte(spec.Meta.Name), 12), nodeModel.Meta.Name)

		if nodeModel.Extensions.AccessCounter {
			tbl.AddColumn(&modeler.Column{
				Name: "ext_access_counter",
				Type: "uint32",
			})
		}

		if nodeModel.Extensions.CommentPerEntry {
			tbl.AddColumn(&modeler.Column{
				Name:    "ext_comment_perentry",
				Type:    "uint8",
				Default: "1",
			})
		}

		if nodeModel.Extensions.Permalink != "" &&
			nodeModel.Extensions.Permalink != "off" {
			tbl.AddColumn(&modeler.Column{
				Name:   "ext_permalink_name",
				Type:   "string",
				Length: "100",
			})
			tbl.AddColumn(&modeler.Column{
				Name:   "ext_permalink_idx",
				Type:   "string",
				Length: "12",
			})
			tbl.AddIndex(&modeler.Index{
				Name: "ext_permalink_idx",
				Type: modeler.IndexTypeIndex,
				Cols: []string{"ext_permalink_idx"},
			})
		}

		for _, field := range nodeModel.Fields {

			switch field.Type {

			case "string":

				tbl.AddColumn(&modeler.Column{
					Name:   "field_" + field.Name,
					Type:   "string",
					Length: field.Length,
				})

				switch field.IndexType {
				case modeler.IndexTypeUnique, modeler.IndexTypeIndex:
					tbl.AddIndex(&modeler.Index{
						Name: "field_" + field.Name,
						Type: field.IndexType,
						Cols: []string{"field_" + field.Name},
					})
				}

			case "text":

				tbl.AddColumn(&modeler.Column{
					Name: "field_" + field.Name,
					Type: "string-text",
				})

				tbl.AddColumn(&modeler.Column{
					Name:   "field_" + field.Name + "_attrs",
					Type:   "string",
					Length: "200",
				})

			case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":

				tbl.AddColumn(&modeler.Column{
					Name: "field_" + field.Name,
					Type: field.Type,
				})

			}
		}

		for _, term := range nodeModel.Terms {

			switch term.Type {

			case api.TermTag:

				tbl.AddColumn(&modeler.Column{
					Name:   "term_" + term.Meta.Name,
					Type:   "string",
					Length: "200",
				})

				// tbl.AddColumn(&modeler.Column{
				// 	Name: "term_" + term.Meta.Name + "_body",
				// 	Type: "string-text",
				// })

				tbl.AddColumn(&modeler.Column{
					Name:   "term_" + term.Meta.Name + "_idx",
					Type:   "string",
					Length: "100",
				})

				tbl.AddIndex(&modeler.Index{
					Name: "term_" + term.Meta.Name + "_idx",
					Type: modeler.IndexTypeIndex,
					Cols: []string{"term_" + term.Meta.Name + "_idx"},
				})

			case api.TermTaxonomy:

				tbl.AddColumn(&modeler.Column{
					Name: "term_" + term.Meta.Name,
					Type: "uint32",
				})

				tbl.AddIndex(&modeler.Index{
					Name: "term_" + term.Meta.Name,
					Type: modeler.IndexTypeIndex,
					Cols: []string{"term_" + term.Meta.Name},
				})
			}

		}

		ds.Tables = append(ds.Tables, &tbl)
	}

	// terms
	for _, termModel := range spec.TermModels {

		var tbl modeler.Table

		if err := json.Decode([]byte(dsTplTermModels), &tbl); err != nil {
			continue
		}

		tbl.Name = fmt.Sprintf("tx%s_%s", idhash.HashToHexString([]byte(spec.Meta.Name), 12), termModel.Meta.Name)

		switch termModel.Type {

		case api.TermTag:

			tbl.AddColumn(&modeler.Column{
				Name:   "uid",
				Type:   "string",
				Length: "16",
			})

			tbl.AddIndex(&modeler.Index{
				Name: "uid",
				Type: modeler.IndexTypeUnique,
				Cols: []string{"uid"},
			})

		case api.TermTaxonomy:

			tbl.AddColumn(&modeler.Column{
				Name: "pid",
				Type: "uint32",
			})

			tbl.AddIndex(&modeler.Index{
				Name: "pid",
				Type: modeler.IndexTypeIndex,
				Cols: []string{"pid"},
			})

			tbl.AddColumn(&modeler.Column{
				Name: "weight",
				Type: "int16",
			})

			tbl.AddIndex(&modeler.Index{
				Name: "weight",
				Type: modeler.IndexTypeIndex,
				Cols: []string{"weight"},
			})

		default:
			continue
		}

		ds.Tables = append(ds.Tables, &tbl)
	}

	// sync
	ms, err := store.Data.Modeler()
	if err != nil {
		return err
	}
	opts := Config.IoConnectors.Options("hpress_database")
	if opts == nil {
		return errors.New("No Database Setup")
	}
	return ms.Sync(opts.Value("dbname"), ds)
}
