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

package store

import (
	"github.com/lessos/lessdb/skv"
	skvdrv "github.com/lessos/lessdb/skv/goleveldb"
)

var (
	CacheDB skv.DB
	err     error
	errInit = &skv.Reply{Status: "ClientError"}
)

func Init(cfg skv.Config) error {

	if CacheDB, err = skvdrv.Open(cfg); err != nil {
		return err
	}

	return nil
}

func CacheSetBytes(key, value []byte, ttl int64) *skv.Reply {

	if CacheDB == nil {
		return errInit
	}

	return CacheDB.KvPut(key, value, ttl)
}

func CacheSet(key, value string, ttl int64) *skv.Reply {
	return CacheSetBytes([]byte(key), []byte(value), ttl)
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
