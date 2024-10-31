package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type JenisApprovalHandler struct {
	JenisApprovalService master.JenisApprovalService
	Config               *configs.Config
}

func ProvideJenisApprovalHandler(service master.JenisApprovalService, config *configs.Config) JenisApprovalHandler {
	return JenisApprovalHandler{
		JenisApprovalService: service,
		Config:               config,
	}
}

func (h *JenisApprovalHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-approval", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllData)
		})
	})
}

// GetDataAll list all jenis approval.
// @Summary Get list all jenis approval.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data jenis approval sesuai dengan filter yang dikirimkan.
// @Tags jenis-approval
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id query string false "Set ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-approval/all [get]
func (h *JenisApprovalHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	status, err := h.JenisApprovalService.GetAll(ids)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
