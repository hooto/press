package datax

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strings"

	"../api"

	"github.com/lessos/lessgo/data/rdo"
	rdobase "github.com/lessos/lessgo/data/rdo/base"
)

var (
	spaceReg = regexp.MustCompile(" +")
)

type TermList api.TermList

func (t *TermList) Index() string {

	if len(t.Items) < 1 {
		return ""
	}

	idxs := []string{}
	for _, v := range t.Items {
		idxs = append(idxs, fmt.Sprintf("%v", v.ID))
	}

	return strings.Join(idxs, ",")
}

func (t *TermList) Content() string {

	if len(t.Items) < 1 {
		return ""
	}

	ts := []string{}
	for _, v := range t.Items {
		ts = append(ts, v.Title)
	}

	return strings.Join(ts, ",")
}

func NodeTermsQuery(model *api.NodeModel, terms []api.NodeTerm) []api.NodeTerm {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return terms
	}

	for _, modTerm := range model.Terms {

		for k, term := range terms {

			if modTerm.Metadata.Name != term.Name {
				continue
			}

			table := fmt.Sprintf("tx%s_%s", model.SpecID, modTerm.Metadata.Name)

			q := rdobase.NewQuerySet().From(table)

			switch modTerm.Type {

			case api.TermTag:

				// TODO

			case api.TermTaxonomy:

				q.Limit(1)
				q.Where.And("id", term.Value)

				if rs, err := dcn.Base.Query(q); err == nil && len(rs) > 0 {

					terms[k].Items = append(terms[k].Items, api.Term{
						ID:    rs[0].Field("id").Uint32(),
						Title: rs[0].Field("title").String(),
					})
				}
			}

			break
		}
	}

	return terms
}

func TermSync(specid, modelid, terms string) (TermList, error) {

	ls := TermList{}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return ls, err
	}

	terms = spaceReg.ReplaceAllString(terms, " ")

	tars := strings.Split(terms, ",")

	ids := []interface{}{}

	for _, term := range tars {

		tag := api.Term{
			Title: strings.TrimSpace(term),
		}

		if len(tag.Title) < 1 {
			continue
		}

		h := md5.New()
		io.WriteString(h, strings.ToLower(tag.Title))
		tag.UID = fmt.Sprintf("%x", h.Sum(nil))[:16]

		exist := false
		for _, prev := range ids {
			if prev.(string) == tag.UID {
				exist = true
				break
			}
		}
		if exist {
			continue
		}

		ls.Items = append(ls.Items, tag)

		ids = append(ids, tag.UID)
	}

	table := fmt.Sprintf("tx%s_%s", specid, modelid)

	if len(ids) > 0 {

		q := rdobase.NewQuerySet().From(table).Limit(int64(len(ids)))
		q.Where.And("uid.in", ids...)

		if rs, err := dcn.Base.Query(q); err == nil {
			for _, v := range rs {

				for tk, tv := range ls.Items {
					if v.Field("uid").String() == tv.UID {
						ls.Items[tk].ID = v.Field("id").Uint32()
						break
					}
				}
			}
		}
	}

	timenow := rdobase.TimeNow("datetime")

	for tk, tv := range ls.Items {

		if tv.ID > 0 {
			continue
		}

		rs, err := dcn.Base.Insert(table, map[string]interface{}{
			"uid":     tv.UID,
			"title":   tv.Title,
			"userid":  "sysadmin",
			"state":   1,
			"created": timenow,
			"updated": timenow,
		})

		if err == nil {
			if incrid, err := rs.LastInsertId(); err == nil && incrid > 0 {
				ls.Items[tk].ID = uint32(incrid)
			}
		}
	}

	return ls, nil
}
