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

package v1

import (
	"fmt"
	"runtime"
	"strings"
	"syscall"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"

	"../../api"
	"../../conf"
	"../../status"
)

type Sys struct {
	*httpsrv.Controller
}

func (c Sys) StatusAction() {

	set := api.SysStatus{}

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("sys.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	set.InstanceID = conf.Config.InstanceID
	set.AppVersion = conf.Version
	set.RuntimeVersion = runtime.Version()
	set.Uptime = status.Uptime
	set.CoroutineNumber = runtime.NumGoroutine()

	ms := memStatsFetch()

	set.MemStats.Alloc = ms.Alloc
	set.MemStats.TotalAlloc = ms.TotalAlloc
	set.MemStats.Sys = ms.Sys

	set.MemStats.NextGC = ms.NextGC
	set.MemStats.LastGC = ms.LastGC
	set.MemStats.PauseTotalNs = ms.TotalAlloc
	set.MemStats.NumGC = ms.NumGC

	set.Info = sysinfoFetch()

	set.Kind = "SysStatus"
}

func (c Sys) IdentityStatusAction() {

	var sets api.SysIdentityStatus

	if !c.Session.AccessAllowed("sys.admin") {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	host := c.Request.Host
	if i := strings.Index(host, ":"); i > 0 {
		host = host[:i]
	}

	insturl := "http://" + host
	if conf.Config.HttpPort != 80 {
		insturl += fmt.Sprintf(":%d", conf.Config.HttpPort)
	}

	sets = api.SysIdentityStatus{
		ServiceUrl: idclient.ServiceUrl,
		InstanceSelf: idsapi.AppInstance{
			Meta: types.ObjectMeta{
				ID: conf.Config.InstanceID,
			},
			AppID:      "lesscms",
			AppTitle:   conf.Config.AppTitle,
			Version:    conf.Version,
			Privileges: conf.Perms,
			Url:        insturl,
		},
	}

	hc := httpclient.Get(idclient.ServiceUrl +
		"/v1/my-app/inst-entry?instid=" + conf.Config.InstanceID +
		"&access_token=" + c.Session.AccessToken)

	var info idsapi.AppInstance

	if err := hc.ReplyJson(&info); err == nil {
		sets.InstanceRegistered = info
	} else {
		sets.InstanceRegistered.Error = &types.ErrorMeta{}
	}

	sets.Kind = "SysIdentityStatus"

	c.RenderJson(sets)
}

func memStatsFetch() runtime.MemStats {

	var ms runtime.MemStats

	runtime.ReadMemStats(&ms)

	return ms
}

func sysinfoFetch() api.SysStatusInfo {

	var si syscall.Sysinfo_t
	syscall.Sysinfo(&si)

	return api.SysStatusInfo{
		CpuNum:    runtime.NumCPU(),
		Uptime:    si.Uptime,
		Loads:     si.Loads,
		MemTotal:  si.Totalram,
		MemFree:   si.Freeram,
		MemShared: si.Sharedram,
		MemBuffer: si.Bufferram,
		MemUsed:   si.Totalram - si.Freeram,
		SwapTotal: si.Totalswap,
		SwapFree:  si.Freeswap,
		Procs:     si.Procs,
		// TimeNow: time.Now().Format(time.RFC3339),
	}
}
