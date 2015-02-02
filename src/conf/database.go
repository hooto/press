package conf

import (
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
)

func (c *ConfigCommon) DatabaseInstance() (*base.Client, error) {

	dc, err := rdo.NewClient("def", c.Database)
	if err != nil {
		return dc, err
	}

	ds, err := base.LoadDataSetFromString(dsBase)
	err = dc.Dialect.SchemaSync(c.Database.Dbname, ds)

	return dc, err
}
