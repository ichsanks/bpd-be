package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type SppdDokumenHandler struct {
	SppdDokumenService bpd.SppdDokumenService
	Config             *configs.Config
}

func ProvideSppdDokumenHandler(service bpd.SppdDokumenService, config *configs.Config) SppdDokumenHandler {
	return SppdDokumenHandler{
		SppdDokumenService: service,
		Config:             config,
	}
}

func (h *SppdDokumenHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/sppd-dokumen", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/upload", h.CreateDokumen)
			r.Post("/", h.Create)
		})
	})
}

// create adalah untuk menambah data sppd dokumen.
// @Summary menambahkan data sppd dokumen.
// @Description Endpoint ini adalah untuk menambahkan data sppd dokumen.
// @Tags sppd-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/sppd-dokumen/upload [post]
func (h *SppdDokumenHandler) CreateDokumen(w http.ResponseWriter, r *http.Request) {
	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.SppdDokumenService.UploadFileDokumen(w, r, "file", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	response.WithJSON(w, http.StatusCreated, path)
}

// create adalah untuk menambah data Dokumen Surat Perjalanan Dinas.
// @Summary menambahkan data Dokumen Surat Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Dokumen Surat Perjalanan Dinas.
// @Tags sppd-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param sppdDokumen body bpd.SppdDokumenRequest true "Dokumen Surat Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.SppdDokumen}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/sppd-dokumen [post]
func (h *SppdDokumenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.SppdDokumenRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	fmt.Println("UserID", userID)
	data, err := h.SppdDokumenService.Create(reqFormat, userID)
	if err != nil {
		fmt.Print("error create")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}
