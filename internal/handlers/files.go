package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"

	"github.com/go-chi/chi"
)

type FileHandler struct {
	Config *configs.Config
}

func ProvideFileHandler(conf *configs.Config) FileHandler {
	return FileHandler{
		Config: conf,
	}
}

func (h *FileHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/files", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", h.ReadFile)
		})
	})
}

func (h *FileHandler) ReadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("path")
	dir := h.Config.App.File.Dir
	fileLocation := filepath.Join(dir, filename)
	img, err := os.Open(fileLocation)

	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
	}
	defer img.Close()
	// w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
	io.Copy(w, img)
}

// untuk ngambil 1 folder saja
func (h *FileHandler) ReadFileParams(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "path")
	dir := h.Config.App.File.Dir
	fileDir := h.Config.App.File.FotoProfile

	path := filepath.Join(fileDir, filename)
	fileLocation := filepath.Join(dir, path)

	img, err := os.Open(fileLocation)

	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
	}
	defer img.Close()
	io.Copy(w, img)
}
