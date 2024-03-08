package static

import "embed"

//go:embed *.html
var StaticDir embed.FS
