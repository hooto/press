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
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"
	"github.com/lessos/lessids/idclient"

	"./config"
	"./datax"
	"./status"
	"./store"

	cdef "./websrv/frontend"
	cmgr "./websrv/mgr"
	capi "./websrv/v1"

	ext_comment "./modules/core/comment/websrv"
	ext_captcha "github.com/eryx/hcaptcha/captcha"
)

var (
	flagPrefix = flag.String("prefix", "", "the prefix folder path")
)

func init() {
	//
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
	httpsrv.GlobalService.Config.TemplateFuncRegister("SysConfig", config.SysConfigList.FetchString)
	httpsrv.GlobalService.Config.TemplateFuncRegister("HttpSrvBasePath", config.HttpSrvBasePath)
}

func main() {

	//
	flag.Parse()

	//
	if err := config.Initialize(*flagPrefix); err != nil {
		fmt.Println("Error on config.Initialize", err)
		logger.Printf("error", "config.Initialize error: %v", err)
		os.Exit(1)
	}

	if err := store.Init(config.Config.CacheDB); err != nil {
		logger.Printf("error", "store.Init error: %v", err)
		os.Exit(1)
	}

	ext_captcha.DataConnector = store.CacheDB
	if err := ext_captcha.Config(config.CaptchaConfig); err != nil {
		logger.Printf("error", "ext_captcha.Config error: %v", err)
		os.Exit(1)
	}

	idclient.ServiceUrl = config.Config.IdentityServiceUrl

	httpsrv.GlobalService.Config.UrlBasePath = "ap"
	httpsrv.GlobalService.Config.HttpPort = config.Config.HttpPort

	// status
	status.Init()
	datax.Worker()

	//
	// httpsrv.Config.I18n(config.Config.Prefix + "/src/i18n/en.json")
	// httpsrv.Config.I18n(config.Config.Prefix + "/src/i18n/zh_CN.json")

	httpsrv.GlobalService.ModuleRegister("/+/comment", ext_comment.NewModule())
	httpsrv.GlobalService.ModuleRegister("/+/hcaptcha", ext_captcha.WebServerModule())

	//
	httpsrv.GlobalService.ModuleRegister("/v1", capi.NewModule())
	httpsrv.GlobalService.ModuleRegister("/mgr", cmgr.NewModule())
	httpsrv.GlobalService.ModuleRegister("/", cdef.NewModule())

	//
	go http.ListenAndServe(":60001", nil)

	fmt.Println("Running")
	httpsrv.GlobalService.Start()

	select {}
}
