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

package websrv

import (
	"github.com/lessos/lessgo/types"
)

type TypeComment struct {
	types.TypeMeta  `json:",inline"`
	Meta            types.ObjectMeta `json:"meta,omitempty"`
	PID             string           `json:"pid,omitempty"`
	ReferID         string           `json:"refer_id,omitempty"`
	ReferModName    string           `json:"refer_modname,omitempty"`
	ReferDataxTable string           `json:"refer_datax_table,omitempty"`
	Content         string           `json:"content,omitempty"`
	Author          string           `json:"author,omitempty"`
	CaptchaToken    string           `json:"captcha_token,omitempty"`
	CaptchaWord     string           `json:"captcha_word,omitempty"`
}
