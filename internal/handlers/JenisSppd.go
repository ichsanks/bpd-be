package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type JenisSppdHandler struct {
	JenisSppdService master.JenisSppdService
	Config           *configs.Config
}

func ProvideJenisSppdHandler(service master.JenisSppdService, config *configs.Config) JenisSppdHandler {
	return JenisSppdHandler{
		JenisSppdService: service,
		Config:           config,
	}
}

func (h *JenisSppdHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-sppd", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllData)
		})
	})
}

// GetDataAll list all jenis sppd.
// @Summary Get list all jenis sppd.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data jenis sppd sesuai dengan filter yang dikirimkan.
// @Tags jenis-sppd
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-sppd/all [get]
func (h *JenisSppdHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	status, err := h.JenisSppdService.GetAll()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
