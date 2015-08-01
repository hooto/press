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

package api

import (
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessids/idsapi"
)

type SysStatus struct {
	types.TypeMeta  `json:",inline"`
	Meta            types.ObjectMeta `json:"meta,omitempty"`
	InstanceID      string           `json:"instance_id,omitempty"`
	AppVersion      string           `json:"app_version,omitempty"`
	RuntimeVersion  string           `json:"runtime_version,omitempty"`
	Uptime          string           `json:"uptime,omitempty"`
	CoroutineNumber int              `json:"coroutine_number,omitempty"`
	Info            SysStatusInfo    `json:"info"`
	MemStats        SysMemStats      `json:"memstats"`
}

type SysMemStats struct {
	//
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	//
	NextGC       uint64 `json:"next_gc"`
	LastGC       uint64 `json:"last_gc"`
	PauseTotalNs uint64 `json:"pause_total_ns"`
	NumGC        uint32 `json:"num_gc"`
}

type SysStatusInfo struct {
	CpuNum    int       `json:"cpu_num"`
	Uptime    int64     `json:"uptime"`
	Loads     [3]uint64 `json:"loads"`
	MemTotal  uint64    `json:"mem_total"`
	MemFree   uint64    `json:"mem_free"`
	MemShared uint64    `json:"mem_shared"`
	MemBuffer uint64    `json:"mem_buffer"`
	MemUsed   uint64    `json:"mem_used"`
	SwapTotal uint64    `json:"swap_total"`
	SwapFree  uint64    `json:"swap_free"`
	Procs     uint16    `json:"procs"`
}

type SysIdentityStatus struct {
	types.TypeMeta     `json:",inline"`
	ServiceUrl         string             `json:"service_url"`
	InstanceSelf       idsapi.AppInstance `json:"instance_self"`
	InstanceRegistered idsapi.AppInstance `json:"instance_registered"`
}
