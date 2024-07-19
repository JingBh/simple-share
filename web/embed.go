package web

import (
	"embed"
	"github.com/spf13/viper"
	"io/fs"
	"net/http"
	"os"
)

//go:embed dist
var EmbedFs embed.FS

func HttpFs() http.FileSystem {
	if !viper.GetBool("debug") {
		fsys, err := fs.Sub(EmbedFs, "dist")
		if err != nil {
			panic(err)
		}

		return http.FS(fsys)
	}

	return http.FS(os.DirFS("web/dist"))
}
