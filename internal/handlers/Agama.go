package handlers

import (
	"net/http"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
)

type AgamaHandler struct {
	AgamaService master.AgamaService
	Config       *configs.Config
}

func ProvideAgamaHandler(service master.AgamaService, config *configs.Config) AgamaHandler {
	return AgamaHandler{
		AgamaService: service,
		Config:       config,
	}
}

func (h *AgamaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/agama", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllData)
		})
	})
}

// GetDataAll list all agama.
// @Summary Get list all agama.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data agama.
// @Tags bioData
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/agama/all [get]
func (h *AgamaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	status, err := h.AgamaService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
