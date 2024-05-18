package serve

import (
	"embed"
	"io/fs"
)

//go:embed ui
var ui embed.FS

func getUiWeb() fs.FS{
	web ,_:= fs.Sub(ui,"ui")
	return web
}