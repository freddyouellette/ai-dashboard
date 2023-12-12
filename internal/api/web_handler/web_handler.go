package web_handler

import (
	"net/http"
	"os"
	"path/filepath"
)

type WebHandler struct {
	StaticDir string
	IndexFile string
}

func NewWebHandler(staticDir string, indexFile string) *WebHandler {
	return &WebHandler{StaticDir: staticDir, IndexFile: indexFile}
}

func (h *WebHandler) Handler(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	path := filepath.Join(h.StaticDir, r.URL.Path)

	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticDir, h.IndexFile))
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.StaticDir)).ServeHTTP(w, r)
}
