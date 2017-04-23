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

package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"code.hooto.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"

	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/datax"
	"code.hooto.com/hooto/hootopress/status"
	"code.hooto.com/hooto/hootopress/store"

	cdef "code.hooto.com/hooto/hootopress/websrv/frontend"
	cmgr "code.hooto.com/hooto/hootopress/websrv/mgr"
	capi "code.hooto.com/hooto/hootopress/websrv/v1"

	ext_comment "code.hooto.com/hooto/hootopress/modules/core/comment/websrv"
	ext_captcha "github.com/eryx/hcaptcha/captcha"
)

var (
	flagPrefix = flag.String("prefix", "", "the prefix folder path")
)

func init() {
	//
	runtime.GOMAXPROCS(runtime.NumCPU())
	//
	flag.Parse()
}

func main() {

	//
	if err := config.Initialize(*flagPrefix); err != nil {
		fmt.Println("Error on config.Initialize", err)
		logger.Printf("error", "config.Initialize error: %v", err)
		os.Exit(1)
	}

	if err := store.Init(config.Config.IoConnectors); err != nil {
		logger.Printf("error", "store.Init error: %v", err)
		fmt.Println("error", "store.Init error ", err)
		os.Exit(1)
	}

	ext_captcha.DataConnector = store.LocalCache
	if err := ext_captcha.Config(config.CaptchaConfig); err != nil {
		logger.Printf("error", "ext_captcha.Config error: %v", err)
		fmt.Println("ext_captcha.Config error", err)
		os.Exit(1)
	}

	iamclient.ServiceUrl = config.Config.IamServiceUrl
	iamclient.InstanceID = config.Config.InstanceID
	iamclient.InstanceOwner = config.Config.AppInstance.Meta.UserID

	httpsrv.GlobalService.Config.UrlBasePath = config.Config.UrlBasePath
	httpsrv.GlobalService.Config.HttpPort = config.Config.HttpPort

	// status
	status.Init()
	datax.Worker()

	//
	// httpsrv.Config.I18n(config.Prefix + "/src/i18n/en.json")
	// httpsrv.Config.I18n(config.Prefix + "/src/i18n/zh_CN.json")

	httpsrv.GlobalService.ModuleRegister("/+/comment", ext_comment.NewModule())
	httpsrv.GlobalService.ModuleRegister("/+/hcaptcha", ext_captcha.WebServerModule())

	//
	httpsrv.GlobalService.ModuleRegister("/v1", capi.NewModule())
	httpsrv.GlobalService.ModuleRegister("/mgr", cmgr.NewModule())
	httpsrv.GlobalService.ModuleRegister("/", cdef.NewModule())

	//
	go http.ListenAndServe(":60001", nil)

	fmt.Println("Running")
	if err := httpsrv.GlobalService.Start(); err != nil {
		fmt.Println("httpsrv.GlobalService.Start error", err)
		os.Exit(1)
	}

	select {}
}
