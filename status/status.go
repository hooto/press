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

package status

import (
	"sync"
	"time"

	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/utilx"
	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"

	"../config"
)

const (
	IdentityServiceOK           int = 1200
	IdentityServiceUnavailable  int = 1503
	IdentityServiceUnRegistered int = 1501
)

var (
	Uptime                string
	IdentityServiceStatus int
	Locker                sync.RWMutex
)

func init() {
	Uptime = utilx.TimeNow("atom")
}

func Init() {
	go func() {

		for {

			Refresh()

			if IdentityServiceStatus == IdentityServiceUnRegistered ||
				IdentityServiceStatus == IdentityServiceOK {
				break
			}

			time.Sleep(3e9)
		}
	}()
}

func Refresh() {

	Locker.Lock()
	defer Locker.Unlock()

	// Check if Identity Service Available
	hc := httpclient.Get(idclient.ServiceUrl + "/v1/status/info")
	var rsjson struct {
		Status string `json:"status"`
	}

	err := hc.ReplyJson(&rsjson)
	hc.Close()
	if err != nil || rsjson.Status != "OK" {
		IdentityServiceStatus = IdentityServiceUnavailable
	} else { // Check if this Registered to ID Service

		hc = httpclient.Get(idclient.ServiceUrl +
			"/v1/app-auth/info?instance_id=" + config.Config.InstanceID)

		var info idsapi.AppAuthInfo

		if err := hc.ReplyJson(&info); err == nil && info.Kind == "AppAuthInfo" {
			IdentityServiceStatus = IdentityServiceOK
		} else {
			IdentityServiceStatus = IdentityServiceUnRegistered
		}

		hc.Close()
	}
}
