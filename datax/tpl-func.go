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

package datax

import (
	"github.com/hooto/hlang4g/hlang"
	"github.com/hooto/httpsrv"
)

func init() {
	httpsrv.DefaultService.Config.RegisterTemplateFunc("TimeFormat", TimeFormat)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("UnixtimeFormat", UnixtimeFormat)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldHtmlPrint", FieldHtmlPrint)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldStringPrint", FieldStringPrint)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldHtmlSubPrint", FieldHtmlSubPrint)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldDebug", FieldDebug)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldString", FieldString)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldSubString", FieldSubString)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldHtml", FieldHtml)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FieldSubHtml", FieldSubHtml)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("pagelet", Pagelet)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("FilterUri", FilterUri)
	httpsrv.DefaultService.Config.RegisterTemplateFunc("T", hlang.StdLangFeed.Translate)
}
