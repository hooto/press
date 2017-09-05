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

	"github.com/lynkdb/iomix/connect"
	"github.com/lynkdb/iomix/rdb"
	"github.com/lynkdb/iomix/skv"
	"github.com/lynkdb/kvgo"
	"github.com/lynkdb/mysqlgo"
)

var (
	err        error
	Data       rdb.Connector
	LocalCache skv.Connector
)

func Init(cfg connect.MultiConnOptions) error {

	opts := cfg.Options("hpress_local_cache")
	if opts == nil {
		return errors.New("No hpress_local_cache Config.IoConnectors Found")
	}

	if LocalCache, err = kvgo.Open(*opts); err != nil {
		return fmt.Errorf("Can Not Connect To %s, Error: %s", opts.Name, err.Error())
	}

	opts = cfg.Options("hpress_database")
	if opts == nil {
		return errors.New("No hpress_database Config.IoConnectors Found")
	}

	conn, err := mysqlgo.NewConnector(*opts)
	if err != nil {
		return err
	}
	Data = conn

	return nil
}
