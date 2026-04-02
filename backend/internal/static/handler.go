package static

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist dist/*
var assets embed.FS

func NewHandler() (http.Handler, error) {
	distFS, err := fs.Sub(assets, "dist")
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServerFS(distFS)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if requestPath == "." || requestPath == "" {
			http.ServeFileFS(w, r, distFS, "index.html")
			return
		}

		if _, err := fs.Stat(distFS, requestPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, distFS, "index.html")
	}), nil
}

