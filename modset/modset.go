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
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"../api"
	"../config"
)

var (
	modNamePattern        = regexp.MustCompile("^[0-9a-z/]{3,30}$")
	srvNamePattern        = regexp.MustCompile("^[0-9a-z/\\-_]{1,50}$")
	modelNamePattern      = regexp.MustCompile("^[a-z]{1,1}[0-9a-z_]{1,20}$")
	nodeFeildNamePattern  = regexp.MustCompile("^[a-z]{1,1}[0-9a-z_]{1,20}$")
	routePathPattern      = regexp.MustCompile("^[0-9a-zA-Z_/\\-:]{1,30}$")
	routeParamNamePattern = regexp.MustCompile("^[a-z]{1,1}[0-9a-zA-Z_]{0,29}$")
)

func ModNameFilter(name string) (string, error) {

	name = strings.Trim(filepath.Clean(strings.ToLower(name)), "/")

	if mat := modNamePattern.MatchString(name); !mat {
		return "", fmt.Errorf("Invalid Module Name (%s)", name)
	}

	return name, nil
}

func SrvNameFilter(name string) (string, error) {

	name = strings.Trim(filepath.Clean(strings.ToLower(name)), "/")

	if mat := srvNamePattern.MatchString(name); !mat {
		return "", fmt.Errorf("Invalid Service Name (%s)", name)
	}

	return name, nil
}

func ModelNameFilter(name string) (string, error) {

	name = strings.TrimSpace(strings.ToLower(name))

	if mat := modelNamePattern.MatchString(name); !mat {
		return "", fmt.Errorf("Invalid Model Name (%s)", name)
	}

	return name, nil
}

func RoutePathFilter(name string) (string, error) {

	name = strings.TrimSpace(name)

	if mat := routePathPattern.MatchString(name); !mat {
		return "", fmt.Errorf("Invalid Route Path (%s)", name)
	}

	return name, nil
}

func SpecFetch(modname string) (api.Spec, error) {

	var entry api.Spec

	file := fmt.Sprintf("%s/modules/%s/spec.json", config.Config.Prefix, modname)
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

	if entry.SrvName == "" {
		return errors.New("SrvName Not Found")
	}

	_, err := SpecFetch(entry.Meta.Name)
	if err == nil {
		return errors.New("Spec Already Exists ")
	}

	entry.Meta.ResourceVersion = "1"
	entry.Meta.Created = utilx.TimeNow("atom")

	dir := fmt.Sprintf("%s/modules/%s", config.Config.Prefix, entry.Meta.Name)
	if err := os.MkdirAll(dir, 0750); err != nil {
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

	if prev.Title != entry.Title || prev.SrvName != entry.SrvName {

		ver, _ := strconv.ParseUint(prev.Meta.ResourceVersion, 10, 64)

		prev.Meta.ResourceVersion = strconv.FormatUint((ver + 1), 10)

		prev.Title = entry.Title
		prev.SrvName = entry.SrvName
		prev.Meta.Updated = utilx.TimeNow("atom")

		if err := _specSet(prev); err != nil {
			return err
		}
	}

	return err
}

func SpecTermSet(modname string, entry api.TermModel) error {

	if modname == "" {
		return errors.New("modname Not Found")
	}

	prev, err := SpecFetch(modname)
	if err != nil {
		return err
	}

	sync, found := false, false
	for i, termModel := range prev.TermModels {

		if termModel.Meta.Name == entry.Meta.Name {

			found = true

			if prev.Title != entry.Title {
				prev.TermModels[i].Title = entry.Title
				sync = true
			}
		}
	}

	if !found {
		entry.ModName = ""
		prev.TermModels = append(prev.TermModels, entry)
		sync = true
	}

	if sync {

		ver, _ := strconv.ParseUint(prev.Meta.ResourceVersion, 10, 64)

		prev.Meta.ResourceVersion = strconv.FormatUint((ver + 1), 10)

		prev.Meta.Updated = utilx.TimeNow("atom")

		if err := _specSet(prev); err != nil {
			return err
		}
	}

	return err
}

func _keyValueListEqual(ls1, ls2 []api.KeyValue) bool {

	if len(ls1) != len(ls2) {
		return false
	}

	for _, kv1 := range ls1 {

		found := false

		for _, kv2 := range ls2 {

			if kv1.Key == kv2.Key {

				if kv1.Value != kv2.Value {
					return false
				}

				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func _termListEqual(ls1, ls2 []api.TermModel) bool {

	if len(ls1) != len(ls2) {
		return false
	}

	for _, kv1 := range ls1 {

		found := false

		for _, kv2 := range ls2 {

			if kv1.Meta.Name == kv2.Meta.Name {

				if kv1.Type != kv2.Type ||
					kv1.Title != kv2.Title {
					return false
				}

				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func SpecNodeSet(modname string, entry api.NodeModel) error {

	if modname == "" {
		return errors.New("modname Not Found")
	}

	for i, field := range entry.Fields {

		if mat := nodeFeildNamePattern.MatchString(field.Name); !mat {
			return fmt.Errorf("Invalid Field Name (%s)", field.Name)
		}

		if field.Title == "" {
			entry.Fields[i].Title = field.Name
		}

		if !utilx.ArrayContain(field.Type, api.NodeFieldTypes) {
			return fmt.Errorf("Invalid Field Type (%s:%s)", field.Name, field.Type)
		}

		if field.IndexType != 0 && field.IndexType != 1 && field.IndexType != 2 {
			return fmt.Errorf("Invalid Field Index Type (%s:%s)", field.Name, field.IndexType)
		}

		if field.Type == "string" {

			length, _ := strconv.Atoi(field.Length)

			if length < 1 {
				entry.Fields[i].Length = "10"
			} else if length > 200 {
				entry.Fields[i].Length = "200"
			}
		}

		for _, attr := range field.Attrs {

			if mat := nodeFeildNamePattern.MatchString(attr.Key); !mat {
				return fmt.Errorf("Invalid Field Attribute Key (%s)", attr.Key)
			}
		}
	}

	prev, err := SpecFetch(modname)
	if err != nil {
		return err
	}

	sync, found := false, false
	for i, nodeModel := range prev.NodeModels {

		if nodeModel.Meta.Name == entry.Meta.Name {

			found = true

			if prev.Title != entry.Title {
				prev.NodeModels[i].Title = entry.Title
				sync = true
			}

			if nodeModel.Extensions.AccessCounter != entry.Extensions.AccessCounter {
				prev.NodeModels[i].Extensions.AccessCounter = entry.Extensions.AccessCounter
				sync = true
			}

			if nodeModel.Extensions.CommentPerEntry != entry.Extensions.CommentPerEntry {
				prev.NodeModels[i].Extensions.CommentPerEntry = entry.Extensions.CommentPerEntry
				sync = true
			}

			if nodeModel.Extensions.Permalink != entry.Extensions.Permalink {
				prev.NodeModels[i].Extensions.Permalink = entry.Extensions.Permalink
				sync = true
			}

			if len(nodeModel.Fields) != len(entry.Fields) && len(entry.Fields) > 0 {

				prev.NodeModels[i].Fields = entry.Fields
				sync = true

			} else {

				for _, prevField := range nodeModel.Fields {

					field_sync := true

					for _, curField := range entry.Fields {

						if curField.Name == prevField.Name {

							if curField.Title == prevField.Title &&
								curField.Type == prevField.Type &&
								curField.IndexType == prevField.IndexType &&
								curField.Length == prevField.Length &&
								_keyValueListEqual(curField.Attrs, prevField.Attrs) {

								field_sync = false
							}

							break
						}
					}

					if field_sync && len(entry.Fields) > 0 {

						sync = true
						prev.NodeModels[i].Fields = entry.Fields

						break
					}
				}
			}

			if !_termListEqual(nodeModel.Terms, entry.Terms) {

				prev.NodeModels[i].Terms = entry.Terms
				sync = true

				for _, sterm := range entry.Terms {

					ptermok := false

					for i, pterm := range prev.TermModels {

						if pterm.Meta.Name == sterm.Meta.Name {

							ptermok = true
							prev.TermModels[i] = sterm
							break
						}
					}

					if !ptermok {
						prev.TermModels = append(prev.TermModels, sterm)
					}
				}
			}
		}
	}

	if !found {
		entry.ModName = ""
		prev.NodeModels = append(prev.NodeModels, entry)

		sync = true
	}

	if sync {

		ver, _ := strconv.ParseUint(prev.Meta.ResourceVersion, 10, 64)

		prev.Meta.ResourceVersion = strconv.FormatUint((ver + 1), 10)

		prev.Meta.Updated = utilx.TimeNow("atom")

		if err := _specSet(prev); err != nil {
			return err
		}
	}

	return err
}

func SpecActionSet(modname string, entry api.Action) error {

	if modname == "" {
		return errors.New("modname Not Found")
	}

	if mat := modelNamePattern.MatchString(entry.Name); !mat {
		return fmt.Errorf("Invalid Action Name (%s)", entry.Name)
	}

	prev, err := SpecFetch(modname)
	if err != nil {
		return err
	}

	for i, dentry := range entry.Datax {

		if mat := modelNamePattern.MatchString(dentry.Name); !mat {
			return fmt.Errorf("Invalid Datax Name (%s)", dentry.Name)
		}

		types := strings.Split(dentry.Type, ".")
		if len(types) != 2 {
			return fmt.Errorf("Invalid Datax Type (%s:%s)", dentry.Name, dentry.Type)
		}

		if !utilx.ArrayContain(types[1], []string{"list", "entry"}) {
			return fmt.Errorf("Invalid Datax Type (%s:%s)", dentry.Name, dentry.Type)
		}

		if dentry.CacheTTL > (86400 * 30) {
			entry.Datax[i].CacheTTL = 86400 * 30
		}

		switch types[0] {

		case "node":

			if dentry.Query.Limit < 1 {
				entry.Datax[i].Query.Limit = 1
			} else if dentry.Query.Limit > 10000 {
				entry.Datax[i].Query.Limit = 10000
			}

			table_found := false
			for _, nodeModel := range prev.NodeModels {

				if nodeModel.Meta.Name == dentry.Query.Table {
					table_found = true
					break
				}
			}

			if !table_found {
				return fmt.Errorf("Query Table Not Found (%s)", dentry.Query.Table)
			}

		case "term":

			table_found := false
			for _, termModel := range prev.TermModels {

				if termModel.Meta.Name == dentry.Query.Table {
					table_found = true
					break
				}
			}

			if !table_found {
				return fmt.Errorf("Query Table Not Found (%s)", dentry.Query.Table)
			}

		default:
			return fmt.Errorf("Invalid Datax Type (%s:%s)", dentry.Name, dentry.Type)
		}
	}

	sync, found := false, false
	for i, action := range prev.Actions {

		if action.Name == entry.Name {

			found = true

			if len(action.Datax) != len(entry.Datax) && len(entry.Datax) > 0 {

				prev.Actions[i].Datax = entry.Datax
				sync = true

			} else {

				for _, prevDatax := range action.Datax {

					datax_sync := true

					for _, curField := range entry.Datax {

						if curField.Name == prevDatax.Name {

							if curField.Type == prevDatax.Type &&
								curField.Pager == prevDatax.Pager &&
								curField.CacheTTL == prevDatax.CacheTTL &&
								curField.Query.Table == prevDatax.Query.Table &&
								curField.Query.Limit == prevDatax.Query.Limit &&
								curField.Query.Order == prevDatax.Query.Order {

								datax_sync = false
							}

							break
						}
					}

					if datax_sync && len(entry.Datax) > 0 {

						sync = true
						prev.Actions[i].Datax = entry.Datax

						break
					}
				}
			}

		}
	}

	if !found {
		entry.ModName = ""
		prev.Actions = append(prev.Actions, entry)

		sync = true
	}

	if sync {

		ver, _ := strconv.ParseUint(prev.Meta.ResourceVersion, 10, 64)

		prev.Meta.ResourceVersion = strconv.FormatUint((ver + 1), 10)

		prev.Meta.Updated = utilx.TimeNow("atom")

		if err := _specSet(prev); err != nil {
			return err
		}
	}

	return err
}

func _routeParamsEqual(a1, a2 map[string]string) bool {

	if len(a1) != len(a2) {
		return false
	}

	for k, v := range a1 {

		found := false

		for k2, v2 := range a2 {

			if k == k2 {

				if v != v2 {
					return false
				}

				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func SpecRouteSet(modname string, entry api.Route) error {

	if modname == "" {
		return errors.New("modname Not Found")
	}

	var err error

	if entry.Path, err = RoutePathFilter(entry.Path); err != nil {
		return fmt.Errorf("Invalid Action Path (%s)", entry.Path)
	}

	for k, _ := range entry.Params {

		if mat := routeParamNamePattern.MatchString(k); !mat {
			return fmt.Errorf("Invalid Param Name (%s)", k)
		}
	}

	prev, err := SpecFetch(modname)
	if err != nil {
		return err
	}

	sync, found := true, false
	for i, prevRoute := range prev.Router.Routes {

		if prevRoute.Path == entry.Path {

			found = true

			if entry.DataAction == prevRoute.DataAction &&
				entry.Template == prevRoute.Template &&
				_routeParamsEqual(entry.Params, prevRoute.Params) {

				sync = false
			} else {
				entry.ModName = ""
				prev.Router.Routes[i] = entry
			}

			break
		}
	}

	if !found {
		entry.ModName = ""
		prev.Router.Routes = append(prev.Router.Routes, entry)

		sync = true
	}

	if sync {

		ver, _ := strconv.ParseUint(prev.Meta.ResourceVersion, 10, 64)

		prev.Meta.ResourceVersion = strconv.FormatUint((ver + 1), 10)

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
	file := fmt.Sprintf("%s/modules/%s/spec.json", config.Config.Prefix, entry.Meta.Name)

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(int64(len(jsb)))

	if _, err = fp.Write(jsb); err == nil {
		// fmt.Println("config.Modules refresh", entry.Meta.ResourceVersion)
		// config.SpecSet(&entry)
	}

	return err
}

func SpecSchemaSync(spec api.Spec) error {

	var (
		ds      rdobase.DataSet
		timenow = rdobase.TimeNow("datetime")
	)

	// TODO
	config.SpecSet(&spec)
	config.SpecSrvRefresh(spec.SrvName)

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

		if nodeModel.Extensions.AccessCounter {
			tbl.AddColumn(&rdobase.Column{
				Name: "ext_access_counter",
				Type: "uint32",
			})
		}

		if nodeModel.Extensions.CommentPerEntry {
			tbl.AddColumn(&rdobase.Column{
				Name:    "ext_comment_perentry",
				Type:    "uint8",
				Default: "1",
			})
		}

		if nodeModel.Extensions.Permalink != "" &&
			nodeModel.Extensions.Permalink != "off" {
			tbl.AddColumn(&rdobase.Column{
				Name:   "ext_permalink_name",
				Type:   "string",
				Length: "100",
			})
			tbl.AddColumn(&rdobase.Column{
				Name:   "ext_permalink_idx",
				Type:   "string",
				Length: "12",
			})
			tbl.AddIndex(&rdobase.Index{
				Name: "ext_permalink_idx",
				Type: rdobase.IndexTypeIndex,
				Cols: []string{"ext_permalink_idx"},
			})
		}

		for _, field := range nodeModel.Fields {

			switch field.Type {

			case "string":

				tbl.AddColumn(&rdobase.Column{
					Name:   "field_" + field.Name,
					Type:   "string",
					Length: field.Length,
				})

				switch field.IndexType {

				case rdobase.IndexTypeUnique, rdobase.IndexTypeIndex:
					tbl.AddIndex(&rdobase.Index{
						Name: "field_" + field.Name,
						Type: field.IndexType,
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
	if err := dcn.Dialect.SchemaSync(config.Config.Database.Dbname, ds); err != nil {
		return err
	}

	return nil
}
