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

package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lynkdb/iomix/connect"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/iomix/skv"
	"github.com/lynkdb/kvgo"
	"github.com/lynkdb/mysqlgo"
	"github.com/lynkdb/pgsqlgo"
)

var (
	err         error
	Data        rdb.Connector
	DataOptions *connect.ConnOptions
	LocalCache  skv.Connector
)

func Init(cfg connect.MultiConnOptions) error {

	opts := cfg.Options("hpress_local_cache")
	if opts == nil {
		return errors.New("No hpress_local_cache Config.IoConnectors Found")
	}

	if LocalCache, err = kvgo.Open(*opts); err != nil {
		return fmt.Errorf("Can Not Connect To %s, Error: %s", opts.Name, err.Error())
	}

	if opts = cfg.Options("hpress_database"); opts == nil {
		hlog.Print("error", err.Error())
		return errors.New("No hpress_database Config.IoConnectors Found")
	}

	switch opts.Driver {

	case "lynkdb/mysqlgo":
		Data, err = mysqlgo.NewConnector(*opts)

	case "lynkdb/pgsqlgo":
		Data, err = pgsqlgo.NewConnector(*opts)

	default:
		return errors.New("Invalid lynkdb/driver")
	}

	if err != nil {
		hlog.Printf("error", "store_init %s", err.Error())
		return err
	}

	DataOptions = opts

	if err = db_upgrade_0_5(Data); err != nil {
		return err
	}

	return nil
}

func db_upgrade_0_5(data rdb.Connector) error {

	mdr, _ := data.Modeler()

	tbls, _ := mdr.TableDump()
	for _, tbl := range tbls {

		if strings.HasPrefix(tbl.Name, "nx") ||
			strings.HasPrefix(tbl.Name, "tx") ||
			tbl.Name == "modules" {

			for _, cv := range tbl.Columns {

				if cv.Name != "created" && cv.Name != "updated" {
					continue
				}

				if cv.Type != "datetime" {
					continue
				}

				sqls := []string{}

				hlog.Printf("warn", "store_init upgrade table %s, colume %s, to int",
					tbl.Name, cv.Name)

				switch DataOptions.Driver {

				case "lynkdb/mysqlgo":
					sqls = []string{
						fmt.Sprintf("ALTER TABLE %s ADD time_tmp int", tbl.Name),
						fmt.Sprintf("UPDATE %s SET time_tmp = UNIX_TIMESTAMP(%s)", tbl.Name, cv.Name),
						fmt.Sprintf("ALTER TABLE %s DROP column %s", tbl.Name, cv.Name),
						fmt.Sprintf("ALTER TABLE %s CHANGE time_tmp %s int", tbl.Name, cv.Name),
					}

				case "lynkdb/pgsqlgo":
					sqls = []string{
						fmt.Sprintf("ALTER TABLE %s ADD COLUMN time_tmp bigint", tbl.Name),
						fmt.Sprintf("UPDATE %s SET time_tmp = extract(epoch from %s)", tbl.Name, cv.Name),
						fmt.Sprintf("ALTER TABLE %s DROP column %s", tbl.Name, cv.Name),
						fmt.Sprintf("ALTER TABLE %s RENAME time_tmp TO %s", tbl.Name, cv.Name),
					}
				}

				for _, sql := range sqls {
					if _, err := data.ExecRaw(sql); err != nil {
						return err
					}
				}

				hlog.Printf("warn", "store_init upgrade table %s, colume %s, to int, DONE",
					tbl.Name, cv.Name)
			}
		}

		if strings.HasPrefix(tbl.Name, "nx") ||
			strings.HasPrefix(tbl.Name, "tx") ||
			tbl.Name == "sys_config" ||
			tbl.Name == "modules" {

			tbl_name_new := ""
			if tbl.Name[:2] == "nx" {
				tbl_name_new = "hpn_" + tbl.Name[2:]
			} else if tbl.Name[:2] == "tx" {
				tbl_name_new = "hpt_" + tbl.Name[2:]
			} else {
				tbl_name_new = "hp_" + tbl.Name
			}

			hlog.Printf("warn", "store_init rename table %s to %s", tbl.Name, tbl_name_new)

			sql := ""

			switch DataOptions.Driver {

			case "lynkdb/mysqlgo":
				sql = fmt.Sprintf("RENAME TABLE %s TO %s", tbl.Name, tbl_name_new)

			case "lynkdb/pgsqlgo":
				sql = fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tbl.Name, tbl_name_new)
			}

			if _, err := data.ExecRaw(sql); err != nil {
				return err
			}

			hlog.Printf("warn", "store_init rename table %s to %s, DONE", tbl.Name, tbl_name_new)
		}

	}

	return nil
}
