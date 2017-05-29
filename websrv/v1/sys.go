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

package v1

import (
	"fmt"
	"runtime"
	"syscall"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/hooto/hootopress/api"
	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/status"
)

type Sys struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *Sys) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c Sys) ConfigListAction() {

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		c.RenderJson(types.TypeMeta{
			Error: &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"},
		})
		return
	}

	c.RenderJson(config.SysConfigList)
}

func (c Sys) ConfigSetAction() {

	var ls api.SysConfigList

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
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

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
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

func (c Sys) IamStatusAction() {

	var sets api.SysIamStatus

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		sets.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	inst_url := "://" + c.Request.Host
	if c.Request.TLS != nil {
		inst_url = "https" + inst_url
	} else {
		inst_url = "http" + inst_url
	}

	if len(httpsrv.GlobalService.Config.UrlBasePath) > 0 {
		inst_url += "/" + httpsrv.GlobalService.Config.UrlBasePath
	}

	sets = api.SysIamStatus{
		ServiceUrl: iamclient.ServiceUrl,
		InstanceSelf: iamapi.AppInstance{
			Meta: types.InnerObjectMeta{
				ID: config.Config.InstanceID,
			},
			AppID:      config.AppName,
			AppTitle:   config.Config.AppTitle,
			Version:    config.Version,
			Privileges: config.Perms,
			Url:        inst_url,
		},
	}

	hc := httpclient.Get(fmt.Sprintf("%s/v1/my-app/inst-entry?instid=%s&%s=%s",
		iamclient.ServiceUrl, config.Config.InstanceID,
		iamclient.AccessTokenKey, iamclient.SessionAccessToken(c.Session)))

	var info iamapi.AppInstance

	if err := hc.ReplyJson(&info); err == nil {
		sets.InstanceRegistered = info
	} else {
		sets.InstanceRegistered.Error = &types.ErrorMeta{}
	}

	hc.Close()

	sets.Kind = "SysIamStatus"

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
