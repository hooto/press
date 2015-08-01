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

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessids/idclient"

	"./conf"
	"./datax"
	"./status"
)

import (
	ext_comment "../modules/comment/websrv"
	cdef "./websrv/frontend"
	cmgr "./websrv/mgr"
	capi "./websrv/v1"
	ext_captcha "github.com/eryx/hcaptcha/captcha"
)

var (
	flagPrefix     = flag.String("prefix", "", "the prefix folder path")
	flagCpuprofile = flag.String("pprof", "", "the pprof path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// render functions
	httpsrv.GlobalService.Config.TemplateFuncRegister("TimeFormat", datax.TimeFormat)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldDebug", datax.FieldDebug)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldString", datax.FieldString)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldSubString", datax.FieldSubString)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldHtml", datax.FieldHtml)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FieldSubHtml", datax.FieldSubHtml)
	httpsrv.GlobalService.Config.TemplateFuncRegister("pagelet", datax.Pagelet)
	httpsrv.GlobalService.Config.TemplateFuncRegister("FilterUri", datax.FilterUri)
}

func main() {

	//
	flag.Parse()
	if err := conf.Initialize(*flagPrefix); err != nil {
		fmt.Println("Error on conf.Initialize", err)
		logger.Printf("error", "conf.Initialize error: %v", err)
		os.Exit(1)
	}

	idclient.ServiceUrl = conf.Config.IdentityServiceUrl

	// httpsrv.Config.UrlBasePath = "cmf"
	httpsrv.GlobalService.Config.HttpPort = conf.Config.HttpPort
	// httpsrv.Config.LessIdsServiceUrl = conf.Config.LessIdsUrl

	// status
	status.Init()
	// for i := 0; i < 3; i++ {

	// 	status.Refresh()

	// 	if status.IdentityServiceStatus == status.IdentityServiceUnRegistered ||
	// 		status.IdentityServiceStatus == status.IdentityServiceOK {
	// 		break
	// 	}

	// 	time.Sleep(3e9)
	// }

	// conf.SpecRefresh("c8f0ltxp")

	//
	// httpsrv.Config.I18n(conf.Config.Prefix + "/src/i18n/en.json")
	// httpsrv.Config.I18n(conf.Config.Prefix + "/src/i18n/zh_CN.json")

	httpsrv.GlobalService.ModuleRegister("/+/comment", ext_comment.NewModule())
	httpsrv.GlobalService.ModuleRegister("/+/hcaptcha", ext_captcha.WebServerModule())

	//
	httpsrv.GlobalService.ModuleRegister("/v1", capi.NewModule())
	httpsrv.GlobalService.ModuleRegister("/mgr", cmgr.NewModule())
	httpsrv.GlobalService.ModuleRegister("/", cdef.NewModule())

	//
	fmt.Println("Running")
	httpsrv.GlobalService.Start()
}
