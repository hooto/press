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

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"

	"../../api"
	"../../config"
	"../../status"
)

type Sys struct {
	*httpsrv.Controller
}

func (c Sys) ConfigListAction() {

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		c.RenderJson(types.TypeMeta{
			Error: &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"},
		})
		return
	}

	c.RenderJson(config.SysConfigList)
}

func (c Sys) ConfigSetAction() {

	var ls api.SysConfigList

	defer c.RenderJson(&ls)

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	err := c.Request.JsonDecode(&ls)
	if err != nil {
		ls.Error = &types.ErrorMeta{api.ErrCodeBadArgument, "Bad Argument " + err.Error()}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		ls.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return
	}

	for _, entry := range ls.Items {

		if prev := config.SysConfigList.Fetch(entry.Key); prev == nil {
			continue
		}

		q := rdobase.NewQuerySet().From("sys_config").Limit(1)
		q.Where.And("key", entry.Key)

		rs, err := dcn.Base.Query(q)
		if err != nil {
			ls.Error = &types.ErrorMeta{
				Code:    api.ErrCodeInternalError,
				Message: "Can not pull database instance",
			}
			return
		}

		set := map[string]interface{}{
			"value": entry.Value,
		}

		if len(rs) > 0 {

			if rs[0].Field("value").String() != entry.Value {

				ft := rdobase.NewFilter()
				ft.And("key", entry.Key)
				_, err = dcn.Base.Update("sys_config", set, ft)
			}

		} else {

			set["key"] = entry.Key

			_, err = dcn.Base.Insert("sys_config", set)
		}

		if err != nil {
			ls.Error = &types.ErrorMeta{
				Code:    api.ErrCodeInternalError,
				Message: err.Error(),
			}
			return
		}

		config.SysConfigList.Insert(entry)
	}

	ls.Kind = "SysConfigList"
}

func (c Sys) StatusAction() {

	set := api.SysStatus{}

	defer c.RenderJson(&set)

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	set.InstanceID = config.Config.InstanceID
	set.AppVersion = config.Version
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

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	host := c.Request.Host
	if i := strings.Index(host, ":"); i > 0 {
		host = host[:i]
	}

	insturl := "http://" + host
	if config.Config.HttpPort != 80 {
		insturl += fmt.Sprintf(":%d", config.Config.HttpPort)
	}

	sets = api.SysIdentityStatus{
		ServiceUrl: idclient.ServiceUrl,
		InstanceSelf: idsapi.AppInstance{
			Meta: types.ObjectMeta{
				ID: config.Config.InstanceID,
			},
			AppID:      "lesscms",
			AppTitle:   config.Config.AppTitle,
			Version:    config.Version,
			Privileges: config.Perms,
			Url:        insturl,
		},
	}

	hc := httpclient.Get(fmt.Sprintf("%s/v1/my-app/inst-entry?instid=%s&%s=%s",
		idclient.ServiceUrl, config.Config.InstanceID,
		idclient.AccessTokenKey, idclient.SessionAccessToken(c.Session)))

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
