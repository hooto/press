// Copyright 2015~2017 hooto Author, All rights reserved.
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

	"code.hooto.com/lynkdb/iomix/connect"
	"code.hooto.com/lynkdb/iomix/skv"
	"code.hooto.com/lynkdb/kvgo"
)

var (
	err        error
	LocalCache skv.Connector
)

func Init(cfg connect.MultiConnOptions) error {

	opts := cfg.Options("htp_local_cache")
	if opts == nil {
		return errors.New("No htp_local_cache Config.IoConnectors Found")
	}

	if LocalCache, err = kvgo.Open(*opts); err != nil {
		return fmt.Errorf("Can Not Connect To %s, Error: %s", opts.Name, err.Error())
	}

	return nil
}

/*
func CacheSet(key, value string, ttl int64) *skv.Reply {
	return CacheSetBytes([]byte(key), []byte(value), ttl)
}

func CacheSetBytes(key, value []byte, ttl int64) *skv.Reply {

	if CacheDB == nil {
		return errInit
	}

	return CacheDB.KvPut(key, value, ttl)
}

func CacheSetJson(key string, value interface{}, ttl int64) *skv.Reply {

	if CacheDB == nil {
		return errInit
	}

	return CacheDB.KvPutJson([]byte(key), value, ttl)
}

func CacheGet(key string) *skv.Reply {

	if CacheDB == nil {
		return errInit
	}

	return CacheDB.KvGet([]byte(key))
}

func CacheDel(key string) *skv.Reply {

	if CacheDB == nil {
		return errInit
	}

	return CacheDB.KvDel([]byte(key))
}
*/
