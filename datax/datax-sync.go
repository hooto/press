// Copyright 2018 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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
	"fmt"
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/mysqlgo"
	"github.com/lynkdb/postgrego"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/store"
)

func data_sync_pull() error {

	if len(config.Config.ExtUpDatabases) == 0 {
		return nil
	}

	var cfgs types.KvPairs
	if rs := store.LocalCache.KvGet(api.NsSysDataPull()); rs.OK() {
		rs.Decode(&cfgs)
	}

	var (
		limit int64 = 100
		src   rdb.Connector
		err   error
		tng   = uint32(time.Now().Unix())
		dtbs  types.ArrayString
	)

	dmr, err := store.Data.Modeler()
	if err != nil {
		return err
	}

	if tbs, err := dmr.TableDump(); err != nil {
		return err
	} else {

		for _, vt := range tbs {
			dtbs.Set(vt.Name)
		}
	}

	for _, cv := range config.Config.ExtUpDatabases {

		// fmt.Println("\n\ndb sync", cv.Name)

		if src != nil {
			// TODO
			// src.Close()
		}

		switch cv.Driver {
		case "lynkdb/mysqlgo":
			src, err = mysqlgo.NewConnector(*cv)

		case "lynkdb/postgrego":
			src, err = postgrego.NewConnector(*cv)

		default:
			continue
		}

		if err != nil {
			hlog.Printf("warn", "data connect ((%s) error : %s",
				cv.Name, err.Error())
			continue
		}

		if src == nil {
			continue
		}

		mr, err := src.Modeler()
		if err != nil {
			return err
		}

		tbs, err := mr.TableDump()
		if err != nil {
			return err
		}

		for _, vt := range tbs {

			if !strings.HasPrefix(vt.Name, "hpt_") &&
				!strings.HasPrefix(vt.Name, "hpn_") {
				continue
			}

			if !dtbs.Has(vt.Name) {
				continue
			}

			var (
				cn, cu  = 0, 0
				q       = src.NewQueryer().From(vt.Name).Limit(limit)
				offset  = int64(0)
				up_name = fmt.Sprintf("sync-time/%s:%s/%s",
					cv.Value("host"), cv.Value("port"), vt.Name)
			)
			err = nil

			if pv := cfgs.Get(up_name); len(pv.String()) > 10 {
				q.Where().And("updated.ge", pv.String())
			}

			// fmt.Println("\nTABLE", vt.Name, tn, tng)

			for {

				rs, err := src.Query(q)
				if err != nil {
					break
				}

				for _, v := range rs {

					sets := map[string]interface{}{}
					for k, f := range v.Fields {
						if k == "ext_access_counter" {
							continue
						}
						sets[k] = f.String()
					}

					qr := store.Data.NewQueryer().From(vt.Name)
					fr := store.Data.NewFilter().And("id", v.Field("id").String())
					qr.SetFilter(fr)
					rsi, err := store.Data.Fetch(qr)

					if rsi.NotFound() {
						if _, err = store.Data.Insert(vt.Name, sets); err != nil {
							// fmt.Println("  ER INSERT", vt.Name, v.Field("id").String(), err.Error())
							break
						} else {
							// fmt.Println("  OK INSERT", vt.Name, v.Field("id").String())
							cn += 1
						}
					} else if err != nil {
						// fmt.Printf(" TABLE %s, ID %s, ER %s\n", vt.Name, v.Field("id").String(), err.Error())
						break
					} else {

						var (
							tup = v.Field("updated").Uint32()
							tlc = rsi.Field("updated").Uint32()
						)

						if tup > tlc {

							if _, err = store.Data.Update(vt.Name, sets, fr); err != nil {
								// fmt.Println("  ER UPDATE", vt.Name, v.Field("id").String())
								break
							} else {
								// fmt.Println("  OK UPDATE", vt.Name, v.Field("id").String())
								cu += 1
							}
						}

						continue
					}

				}

				if err != nil || len(rs) < int(limit) {
					// fmt.Printf("  DONE INSERT/IGNORE %d, UPDATE %d, ALL %d\n",
					// 	cn, cu, int(offset)+len(rs))
					break
				}

				offset += limit
			}

			if err == nil {
				if cn > 0 || cu > 0 {
					hlog.Printf("info", "data sync (%s) INSERT %d, UPDATE %d", up_name, cn, cu)
					cfgs.Set(up_name, tng)
				}
			} else {
				hlog.Printf("warn", "data sync ((%s) error : %s",
					up_name, err.Error())
			}
		}
	}

	if rs := store.LocalCache.KvPut(api.NsSysDataPull(), cfgs, nil); !rs.OK() {
		// fmt.Println("  DATA PULL TAG ERROR")
	}

	return nil
}
