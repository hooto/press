package datax

import (
	// "fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"../api"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
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
	regMultiLine  = regexp.MustCompile("\\n+")
	regMultiSpace = regexp.MustCompile("\\s{2,}")
	regLineSpace  = regexp.MustCompile("\\n\\s*\\n")
)

func TimeFormat(timeString, formatFrom, formatTo string) string {

	tp, err := time.ParseInLocation(timeFormator.Replace(formatFrom), timeString, time.Local)
	if err != nil {
		return timeString
	}

	return tp.Format(timeFormator.Replace(formatTo))
}

func FieldTimeFormat(fields []api.NodeField, colname, format string) string {

	val, _ := fieldValue(fields, colname)

	tp, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
	if err != nil {
		return val
	}

	format = timeFormator.Replace(format)

	return tp.Format(format)
}

func textHtml2str(src string) string {

	src = regElement.ReplaceAllStringFunc(src, strings.ToLower)

	src = regStyle.ReplaceAllString(src, "")
	src = regScript.ReplaceAllString(src, "")

	src = regElement.ReplaceAllString(src, "\n")
	src = regMultiSpace.ReplaceAllString(src, "\n")

	src = regMultiLine.ReplaceAllString(src, "\n\n")

	return strings.TrimSpace(src)
}

// func textAutoParagraph(src string) string {
// 	lines := strings.Split(src, "\n\n")
// 	return "<p>" + strings.Join(lines, "</p>\n<p>") + "</p>"
// }

func fieldValue(fields []api.NodeField, colname string) (string, map[string]string) {

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

func FieldString(fields []api.NodeField, colname string) string {

	val, _ := fieldValue(fields, colname)

	return val
}

func FieldSubString(fields []api.NodeField, colname string, length int) string {

	if length < 1 {
		length = fieldStringMaxLen
	}

	val, _ := fieldValue(fields, colname)

	if len(val) > length {
		return val[:length] + "..."
	}

	return val
}

func FieldDebug(fields []api.NodeField, colname string, length int) template.HTML {

	val, attrs := fieldValue(fields, colname)

	if v, ok := attrs["format"]; ok {

		if v == "md" {
			unsafe := blackfriday.MarkdownCommon([]byte(val))
			val = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
			v = "html"
		}

		if v == "html" {
			val = textHtml2str(val)
		}
	}

	if len(val) > length {
		val = val[:length] + "..."
	}

	return template.HTML(val)
}

func FieldHtml(fields []api.NodeField, colname string) template.HTML {

	val, attrs := fieldValue(fields, colname)

	if v, ok := attrs["format"]; ok && v == "md" {
		unsafe := blackfriday.MarkdownCommon([]byte(val))
		val = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	}

	return template.HTML(val)
}

func FieldSubHtml(fields []api.NodeField, colname string, length int) template.HTML {

	if length < 1 {
		length = fieldStringMaxLen
	}

	val, attrs := fieldValue(fields, colname)

	if v, ok := attrs["format"]; ok {

		if v == "md" {
			unsafe := blackfriday.MarkdownCommon([]byte(val))
			val = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
			v = "html"
		}

		if v == "html" {
			val = textHtml2str(val)
		}
	}

	if len(val) > length {
		val = val[:length] + "..."
	}

	lines := strings.Split(val, "\n\n")
	val = "<p>" + strings.Join(lines, "</p>\n<p>") + "</p>"

	return template.HTML(val)
}
