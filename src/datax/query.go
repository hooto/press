package datax

import (
	// "fmt"
	// "../api"
	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
)

type QuerySet struct {
	SpecID  string
	ModelID string
	cols    string
	Table   string
	order   string
	limit   int64
	offset  int64
	filter  rdobase.Filter
	Pager   bool
}

func NewQuery(specid, table string) *QuerySet {
	return &QuerySet{
		SpecID:  specid,
		ModelID: table,
		Table:   table,
		cols:    "*",
		limit:   1,
		offset:  0,
	}
}

func (q *QuerySet) Query() []rdobase.Entry {

	rs := []rdobase.Entry{}

	//
	dc, err := rdo.ClientPull("def")
	if err != nil {
		return rs
	}

	qs := rdobase.NewQuerySet().
		Select(q.cols).
		From(q.Table).
		Order(q.order).
		Limit(q.limit).
		Offset(q.offset)

	qs.Where = q.filter

	rs, err = dc.Base.Query(qs)
	if err != nil {
		return rs
	}

	return rs
}

func (q *QuerySet) QueryEntry() *rdobase.Entry {

	q.limit = 1
	if ls := q.Query(); len(ls) > 0 {
		return &ls[0]
	}

	return nil
}

func (q *QuerySet) Select(s string) string {
	q.cols = s
	return ""
}

func (q *QuerySet) From(s string) string {
	q.Table = s
	return ""
}

func (q *QuerySet) Order(s string) string {
	q.order = s
	return ""
}

func (q *QuerySet) Limit(num int64) string {
	q.limit = num
	return ""
}

func (q *QuerySet) Offset(num int64) string {
	q.offset = num
	return ""
}

func (q *QuerySet) Filter(expr string, args ...interface{}) string {
	q.filter.And(expr, args...)
	return ""
}

func (q *QuerySet) FilterOr(expr string, args ...interface{}) string {
	q.filter.Or(expr, args...)
	return ""
}
