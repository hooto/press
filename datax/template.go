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

package datax

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/store"
)

func FilterUri(data map[string]interface{}, args ...interface{}) template.URL {

	uris := []string{}

	for key, val := range data {

		if len(key) > 5 && key[:5] == "term_" {
			uris = append(uris, fmt.Sprintf("%s=%v", key, val))
		}
	}

	if len(args) > 1 {
		for i := 0; i < len(args); i += 2 {
			uris = append(uris, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		}
	}

	if len(uris) > 0 {
		return template.URL(strings.Join(uris, "&"))
	}

	return ""
}

func Pagelet(data map[string]interface{}, args ...string) template.HTML {

	defer func() {
		if err := recover(); err != nil {
			hlog.Printf("error", "Pagelet Panic %s", err)
		}
	}()

	//
	if len(args) < 2 || len(args) > 3 {
		return ""
	}

	//
	modname, templatePath := args[0], args[1]
	if len(args) == 2 {
		// fmt.Println("Pagelet args=2", modname, args)
		return templateRender(data, modname, templatePath)
	}
	// fmt.Println("Pagelet", modname, args)

	//
	user, _ := data["s_user"]
	if user == "" {
		user = "guest"
	}

	//
	for _, spec := range config.Modules {

		if spec.Meta.Name != modname {
			continue
		}

		dataAction := args[2]

		for _, action := range spec.Actions {

			if action.Name != dataAction {
				continue
			}

			for _, datax := range action.Datax {

				qry := NewQuery(modname, datax.Query.Table)

				if datax.Query.Limit > 0 {
					qry.Limit(datax.Query.Limit)
				}

				if datax.Query.Order != "" {
					qry.Order(datax.Query.Order)
				}

				qry.Filter("status", 1)

				switch datax.Type {

				case "node.list":

					var ls api.NodeList
					qryhash := qry.Hash()
					if datax.CacheTTL > 0 && user != config.Config.AppInstance.Meta.User {
						if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
							rs.Decode(&ls)
						}
					}

					if len(ls.Items) == 0 {
						ls = qry.NodeList([]string{}, []string{})
						if datax.CacheTTL > 0 && len(ls.Items) > 0 {
							store.LocalCache.KvPut([]byte(qryhash), ls, &skv.KvWriteOptions{
								Ttl: datax.CacheTTL,
							})
						}
					}

					data[datax.Name] = ls

				case "node.entry":

					var entry api.Node
					qryhash := qry.Hash()
					if datax.CacheTTL > 0 && user != config.Config.AppInstance.Meta.User {
						if rs := store.LocalCache.KvGet([]byte(qryhash)); rs.OK() {
							rs.Decode(&entry)
						}
					}

					if entry.Title == "" {
						entry = qry.NodeEntry()
						if datax.CacheTTL > 0 && entry.Title != "" {
							store.LocalCache.KvPut([]byte(qryhash), entry, &skv.KvWriteOptions{
								Ttl: datax.CacheTTL,
							})
						}
					}

					data[datax.Name] = entry
				}
			}

			return templateRender(data, spec.Meta.Name, templatePath)
		}

		return templateRender(data, spec.Meta.Name, templatePath)
	}

	//
	return templateRender(data, modname, templatePath)
}

func templateRender(data map[string]interface{}, module, templatePath string) template.HTML {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tplset, err := httpsrv.GlobalService.TemplateLoader.Template(module, templatePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var out bytes.Buffer
	if err = tplset.Render(&out, data); err != nil {
		hlog.Printf("error", "tplset.Render Error %v", err)
		return ""
	}

	return template.HTML(out.String())
}
