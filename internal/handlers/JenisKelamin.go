package handlers

import (
	"net/http"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
)

type JenisKelaminHandler struct {
	JenisKelaminService master.JenisKelaminService
	Config              *configs.Config
}

func ProvideJenisKelaminHandler(service master.JenisKelaminService, config *configs.Config) JenisKelaminHandler {
	return JenisKelaminHandler{
		JenisKelaminService: service,
		Config:              config,
	}
}

func (h *JenisKelaminHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-kelamin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllData)
		})
	})
}

// GetDataAll list all jenis kelamin.
// @Summary Get list all jenis kelamin.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data jenis kelamin.
// @Tags bioData
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kelamin/all [get]
func (h *JenisKelaminHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	status, err := h.JenisKelaminService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
