package v1

import (
	"github.com/lessos/lessgo/pagelet"
)

type TermModel struct {
	*pagelet.Controller
}

func (c TermModel) ListAction() {

	c.AutoRender = false
}
