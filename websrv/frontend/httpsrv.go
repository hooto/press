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

package frontend

import (
	"code.hooto.com/lessos/iam/iamclient"
	"github.com/hooto/httpsrv"

	"github.com/hooto/hpress/config"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("default")

	module.ControllerRegister(new(Index))
	module.ControllerRegister(new(Error))
	module.ControllerRegister(new(S2))

	return module
}

func NewHtpModule() httpsrv.Module {

	module := httpsrv.NewModule("default_hpress")

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~/hchart",
		StaticPath: config.Prefix + "/vendor/github.com/hooto/hchart/webui",
	})

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~",
		StaticPath: config.Prefix + "/webui/",
	})

	module.ControllerRegister(new(iamclient.Auth))

	return module
}
