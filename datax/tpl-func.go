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
	httpsrv.DefaultService.Config.TemplateFuncRegister("TimeFormat", TimeFormat)
	httpsrv.DefaultService.Config.TemplateFuncRegister("UnixtimeFormat", UnixtimeFormat)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldHtmlPrint", FieldHtmlPrint)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldStringPrint", FieldStringPrint)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldHtmlSubPrint", FieldHtmlSubPrint)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldDebug", FieldDebug)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldString", FieldString)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldSubString", FieldSubString)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldHtml", FieldHtml)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FieldSubHtml", FieldSubHtml)
	httpsrv.DefaultService.Config.TemplateFuncRegister("pagelet", Pagelet)
	httpsrv.DefaultService.Config.TemplateFuncRegister("FilterUri", FilterUri)
	httpsrv.DefaultService.Config.TemplateFuncRegister("T", hlang.StdLangFeed.Translate)
}
