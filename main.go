// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamclient"

	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/datax"
	"github.com/hooto/hpress/status"
	"github.com/hooto/hpress/store"

	cdef "github.com/hooto/hpress/websrv/frontend"
	cmgr "github.com/hooto/hpress/websrv/mgr"
	cmod "github.com/hooto/hpress/websrv/module"
	capi "github.com/hooto/hpress/websrv/v1"

	ext_captcha "github.com/hooto/hcaptcha/captcha4g"
	ext_comment "github.com/hooto/hpress/modules/core/comment/websrv"
)

var (
	version    = ""
	release    = ""
	flagPrefix = flag.String("prefix", "", "the prefix folder path")
)

func init() {
	//
	runtime.GOMAXPROCS(runtime.NumCPU())
	//
	flag.Parse()
}

func main() {

	if version != "" {
		config.Version = version
	}
	if release != "" {
		config.Release = release
	}

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Version: %s, Release: %s\n", config.Version, config.Release)
		return
	}

	//
	for {

		err := config.Initialize(*flagPrefix)
		if err == nil {
			break
		}

		// fmt.Println("Error on config.Initialize", err)
		hlog.Printf("error", "config.Initialize error: %v", err)
		time.Sleep(3e9)
	}

	ext_captcha.DataConnector = store.LocalCache
	if err := ext_captcha.Config(config.CaptchaConfig); err != nil {
		hlog.Printf("error", "ext_captcha.Config error: %v", err)
		fmt.Println("ext_captcha.Config error", err)
		os.Exit(1)
	}

	iamclient.ServiceUrl = config.Config.IamServiceUrl
	iamclient.ServiceUrlFrontend = config.Config.IamServiceUrlFrontend

	iamclient.InstanceID = config.Config.InstanceID
	iamclient.InstanceOwner = config.Config.AppInstance.Meta.User

	httpsrv.GlobalService.Config.UrlBasePath = config.Config.UrlBasePath
	httpsrv.GlobalService.Config.HttpPort = config.Config.HttpPort

	// status
	status.Init()
	datax.Worker()

	//
	// httpsrv.Config.I18n(config.Prefix + "/src/i18n/en.json")
	// httpsrv.Config.I18n(config.Prefix + "/src/i18n/zh_CN.json")

	httpsrv.GlobalService.ModuleRegister("/hpress/+/comment", ext_comment.NewModule())
	httpsrv.GlobalService.ModuleRegister("/hpress/+/hcaptcha", ext_captcha.WebServerModule())
	httpsrv.GlobalService.ModuleRegister("/hpress/+", cmod.NewModule())

	//
	httpsrv.GlobalService.ModuleRegister("/hpress/v1", capi.NewModule())
	httpsrv.GlobalService.ModuleRegister("/hpress/mgr", cmgr.NewModule())
	httpsrv.GlobalService.ModuleRegister("/hpress", cdef.NewHtpModule())
	httpsrv.GlobalService.ModuleRegister("/", cdef.NewModule())

	//
	if config.Config.HttpPortPprof > 0 {
		go http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HttpPortPprof), nil)
	}

	if err := httpsrv.GlobalService.Start(); err != nil {
		fmt.Println("httpsrv.GlobalService.Start error", err)
		os.Exit(1)
	}

	select {}
}
