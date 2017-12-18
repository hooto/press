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
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
)

const (
	fieldStringMaxLen = 102400
)

var (
	timeFormatMap = []string{
		"Y", "2006",
		"y", "06",
		"m", "01",
		"d", "02",
		"H", "15",
		"h", "03",
		"i", "04",
		"s", "05",
		"date", "2006-01-02",
		"datetime", "2006-01-02 15:04:05",
		"atom", time.RFC3339,
	}
	timeFormator = strings.NewReplacer(timeFormatMap...)
)

var (
	regElement    = regexp.MustCompile("\\<[\\S\\s]+?\\>")
	regStyle      = regexp.MustCompile("\\<style[\\S\\s]+?\\</style\\>")
	regScript     = regexp.MustCompile("\\<script[\\S\\s]+?\\</script\\>")
	regMultiLine  = regexp.MustCompile("\\n\\n+")
	regMultiSpace = regexp.MustCompile("\\s{2,}")
	regLineSpace  = regexp.MustCompile("\\n\\s*\\n")
	mkp           = bluemonday.UGCPolicy()
	htmlp         = bluemonday.UGCPolicy()
	shtmlp        = bluemonday.UGCPolicy()
)

func init() {
	mkp.AllowAttrs("class").OnElements("code")

	//
	shtmlp.AllowElements("script")
	shtmlp.AllowElements("button", "style")

	//
	shtmlp.AllowAttrs("class").OnElements("div")
	shtmlp.AllowAttrs("style").OnElements("div")

	//
	shtmlp.AllowAttrs("class").OnElements("button")
	shtmlp.AllowAttrs("onclick").OnElements("button")

	//
	shtmlp.AllowAttrs("href").OnElements("a")
	shtmlp.AllowAttrs("target").OnElements("a")
	shtmlp.AllowAttrs("class").OnElements("a")
	shtmlp.AllowAttrs("class").OnElements("img")
	shtmlp.AllowAttrs("class").OnElements("span")
}

func TimeFormat(timeString, formatFrom, formatTo string) string {

	tp, err := time.ParseInLocation(timeFormator.Replace(formatFrom), timeString, time.Local)
	if err != nil {
		return timeString
	}

	return tp.Format(timeFormator.Replace(formatTo))
}

func FieldTimeFormat(fields []*api.NodeField, colname, format string) string {

	val, _ := fieldValue(fields, colname)

	tp, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
	if err != nil {
		return val
	}

	format = timeFormator.Replace(format)

	return tp.Format(format)
}

func TextHtml2Str(src string) string {

	src = regElement.ReplaceAllStringFunc(src, strings.ToLower)

	src = regStyle.ReplaceAllString(src, "")
	src = regScript.ReplaceAllString(src, "")

	src = regElement.ReplaceAllString(src, "\n")
	src = regMultiSpace.ReplaceAllString(src, "\n")

	src = regMultiLine.ReplaceAllString(src, "\n\n")

	return strings.TrimSpace(src)
}

func fieldValue(fields []*api.NodeField, colname string) (string, map[string]string) {

	var (
		val   = ""
		attrs = map[string]string{}
	)

	for _, field := range fields {

		if field.Name != colname {
			continue
		}

		for _, v := range field.Attrs {
			attrs[v.Key] = v.Value
		}

		val = field.Value

		break
	}

	return val, attrs
}

func fieldValueCache(fields []*api.NodeField, colname string, cache_key string) (string, map[string]string, bool) {

	var (
		val    = ""
		attrs  = map[string]string{}
		cached = false
	)

	for _, field := range fields {

		if field.Name != colname {
			continue
		}

		for _, v := range field.Attrs {
			attrs[v.Key] = v.Value
		}

		if val, cached = field.ValueCaches[cache_key]; !cached {
			val = field.Value
		}

		break
	}

	return val, attrs, cached
}

func fieldValueCacheSet(fields []*api.NodeField, colname, value, cache_key string) {

	for _, field := range fields {

		if field.Name == colname {
			field.ValueCaches[cache_key] = value
			break
		}
	}
}

func FieldString(fields []*api.NodeField, colname string) string {

	val, _ := fieldValue(fields, colname)

	return val
}

func FieldSubString(fields []*api.NodeField, colname string, length int) string {

	if length < 1 {
		length = fieldStringMaxLen
	}

	val, _, cached := fieldValueCache(fields, colname, fmt.Sprintf("FieldSubString_%d", length))
	if cached {
		return val
	}

	ustr := []rune(val)

	if len(ustr) > length {
		val = string(ustr[0:length]) + "..."
	}

	fieldValueCacheSet(fields, colname, val, fmt.Sprintf("FieldSubString_%d", length))

	return val
}

func FieldDebug(fields []*api.NodeField, colname string, length int) template.HTML {

	val, attrs := fieldValue(fields, colname)

	if v, ok := attrs["format"]; ok {

		if v == "md" {
			unsafe := blackfriday.MarkdownCommon([]byte(val))
			val = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
			v = "html"
		}

		if v == "html" {
			val = TextHtml2Str(val)
		}
	}

	if len(val) > length {
		val = val[:length] + "..."
	}

	return template.HTML(val)
}

func FieldHtml(fields []*api.NodeField, colname string) template.HTML {

	val, attrs, cached := fieldValueCache(fields, colname, "FieldHtml")
	if cached {
		return template.HTML(val)
	}

	fm, ok := attrs["format"]
	if !ok {
		fm = "text"
	}

	val = strings.TrimSpace(strings.Replace(val, "\r\n", "\n", -1))
	val = regMultiLine.ReplaceAllString(val, "\n\n")

	val = strings.Replace(val, "{{lessos_storage_service_uri}}",
		config.SysConfigList.FetchString("ls2_uri"), -1)

	switch fm {

	case "md":
		unsafe := blackfriday.MarkdownCommon([]byte(val))
		val = string(mkp.SanitizeBytes(unsafe))

	case "html":
		val = shtmlp.Sanitize(val)

	case "text":
		if lines := strings.Split(val, "\n\n"); len(lines) > 1 {
			val = "<p>" + strings.Join(lines, "</p><p>") + "</p>"
			val = strings.Replace(val, "\n", "<br>", -1)
		}
		fallthrough

	default:
		val = htmlp.Sanitize(val)
	}

	fieldValueCacheSet(fields, colname, val, "FieldHtml")

	return template.HTML(val)
}

func StringSub(s string, start, length int) string {

	bt := []rune(s)

	if start < 0 {
		start = 0
	}

	if length < 1 {
		length = 1
	}

	end := start + length

	if end >= len(bt) {
		end = len(bt)
	}

	if end <= start {
		return ""
	}

	return string(bt[start:end])
}

func FieldSubHtml(fields []*api.NodeField, colname string, length int) template.HTML {

	if length < 1 {
		length = fieldStringMaxLen
	}

	val, attrs, cached := fieldValueCache(fields, colname, fmt.Sprintf("FieldSubHtml_%d", length))
	if cached {
		return template.HTML(val)
	}

	if v, ok := attrs["format"]; ok {

		if v == "md" {
			unsafe := blackfriday.MarkdownCommon([]byte(val))
			val = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		}
	}

	ustr := []rune(TextHtml2Str(val))

	if len(ustr) > length {
		val = string(ustr[0:length]) + "..."
	} else {
		val = string(ustr)
	}

	val = strings.Replace(val, "\n", "<br>", -1)

	if lines := strings.Split(val, "\n\n"); len(lines) > 1 {
		val = "<p>" + strings.Join(lines, "</p><p>") + "</p>"
	}

	fieldValueCacheSet(fields, colname, val, fmt.Sprintf("FieldSubHtml_%d", length))

	return template.HTML(val)
}
