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
	"github.com/lessos/lessgo/httpsrv"

	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/websrv/mgr/controllers"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("htp_mgr")

	// module.RouteSet(httpsrv.Route{
	// 	Type:       httpsrv.RouteTypeStatic,
	// 	Path:       "~",
	// 	StaticPath: config.Prefix + "/webui/",
	// })

	// module.RouteSet(httpsrv.Route{
	// 	Type:       httpsrv.RouteTypeStatic,
	// 	Path:       "-",
	// 	StaticPath: config.Prefix + "/webui/htpm/",
	// })

	module.TemplatePathSet(config.Prefix + "/websrv/mgr/views")

	module.ControllerRegister(new(controllers.Index))
	module.ControllerRegister(new(controllers.Setup))

	return module
}
