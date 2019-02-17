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

package main

import (
	"fmt"
	"strings"

	"github.com/lessos/lessgo/encoding/json"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/mysqlgo"
	"github.com/lynkdb/pgsqlgo"

	"github.com/hooto/hpress/config"
)

var (
	s2   = "{{lessos_storage_service_uri}}/a01/"
	Data rdb.Connector
	err  error
)

func main() {

	file := "etc/config.json"
	if err := json.DecodeFile(file, &config.Config); err != nil {
		fmt.Println(err)
		return
	}

	dbcfg := config.Config.IoConnectors.Options("hpress_database")
	if dbcfg == nil {
		return
	}

	switch dbcfg.Driver {

	case "lynkdb/mysqlgo":
		Data, err = mysqlgo.NewConnector(*dbcfg)

	case "lynkdb/pgsqlgo":
		Data, err = pgsqlgo.NewConnector(*dbcfg)

	default:
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	mr, err := Data.Modeler()
	if err != nil {
		fmt.Println(err)
		return
	}

	tbs, err := mr.TableDump()
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		limit      int64 = 100
		s2Replacer       = strings.NewReplacer(
			"{{lessos_storage_service_uri}}/a01", "{{hp_storage_service_endpoint}}",
			"{{lessos_storage_service_uri}}/a02", "{{hp_storage_service_endpoint}}",
		)
		num = 0
	)

	for _, vt := range tbs {

		if !strings.HasPrefix(vt.Name, "nx") &&
			!strings.HasPrefix(vt.Name, "hpn_") {
			continue
		}
		var (
			q      = Data.NewQueryer().From(vt.Name).Limit(limit)
			offset = int64(0)
		)

		for {

			rs, err := Data.Query(q)
			if err != nil {
				break
			}

			for _, v := range rs {

				sets := map[string]interface{}{}

				for k, f := range v.Fields {

					if !strings.Contains(f.String(), "lessos_storage_service_uri") {
						continue
					}

					sets[k] = s2Replacer.Replace(f.String())
				}

				if len(sets) > 0 {

					num += 1
					fmt.Println(num)

					ft := Data.NewFilter()
					ft.And("id", v.Field("id").String())

					// sets["updated"] = time.Now().Format("2006-01-02 15:04:05")

					Data.Update(vt.Name, sets, ft)
				}
			}

			if err != nil || len(rs) < int(limit) {
				break
			}

			offset += limit
			q.Offset(offset)
		}

	}
}
