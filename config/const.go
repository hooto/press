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

package config

import (
	"github.com/lessos/iam/iamapi"
)

var (
	coreModules = []string{"core/general", "core/comment", "core/blog", "ruilog/notebook"}

	Perms = []iamapi.AppPrivilege{
		{
			Privilege: "frontend.list",
			Desc:      "Frontend - List",
			Roles:     []uint32{100, 1000},
		},
		{
			Privilege: "frontend.read",
			Desc:      "Frontend - Read",
			Roles:     []uint32{100, 1000},
		},
		{
			Privilege: "editor.list",
			Desc:      "Editor - List",
			Roles:     []uint32{100},
		},
		{
			Privilege: "editor.write",
			Desc:      "Editor - Write",
			Roles:     []uint32{100},
		},
		{
			Privilege: "editor.read",
			Desc:      "Editor - Read",
			Roles:     []uint32{100},
		},
		{
			Privilege: "sys.admin",
			Desc:      "System Admin",
			Roles:     []uint32{1},
		},
	}
)
