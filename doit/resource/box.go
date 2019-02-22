package resource

import (
	"github.com/gobuffalo/packr"
)

var (
	FontBox packr.Box
)

func Load() {
	FontBox = packr.NewBox("./fonts")
}
