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
	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/hooto/hootopress/api"
	"code.hooto.com/hooto/hootopress/config"
)

type NodeModel struct {
	*httpsrv.Controller
}

func (c NodeModel) EntryAction() {

	rsp := api.NodeModel{}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.read", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	modname, modelid := c.Params.Get("modname"), c.Params.Get("modelid")

	nmodel, err := config.SpecNodeModel(modname, modelid)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Model Not Found",
		}
		return
	}

	rsp = *nmodel
	rsp.Kind = "NodeModel"
}
