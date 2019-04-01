package resource

import (
	"github.com/gobuffalo/packr"
	"fmt"
)

var (
	FontBox      packr.Box
)

func Load() {
	FontBox = packr.NewBox("./fonts")
	fmt.Println(FontBox.Path)
}
