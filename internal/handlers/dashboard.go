package handlers

import (
	"net/http"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
)

type DashboardHandler struct {
	DashboardService auth.DashboardService
	Config           *configs.Config
}

func ProvideDashboardHandler(service auth.DashboardService, config *configs.Config) DashboardHandler {
	return DashboardHandler{
		DashboardService: service,
		Config:           config,
	}
}

func (h *DashboardHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/dashboard/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/jml-pegawai", h.GetAllData)
			r.Get("/aktif-bpd", h.GetDataAktifBpd)
			r.Get("/bpd", h.GetDataDashboardBpd)
			r.Get("/sppd", h.GetDataDashboardSppd)
			r.Get("/bpd-new", h.GetDataDashboardBpdNew)
			r.Get("/aktif-sppd", h.GetDataAktifSppd)
			r.Get("/aktif-bpd-new", h.GetDataAktifBpdNew)
			r.Get("/jumlah-sppd", h.GetJumlahSppd)
			r.Get("/jumlah-bpd", h.GetJumlahBpd)
		})
	})
}

// GetDataAll list all jumlah pegawai.
// @Summary Get list all jumlah pegawai.
// @Description endpoint ini digunakan untuk mendapatkan jumlah seluruh pegawai aktif.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/jml-pegawai [get]
func (h *DashboardHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	status, err := h.DashboardService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataDashboardBpd jml data belum proses pengajuan atau pengajuan ataun dalam dinas dan penyelesaian.
// @Summary Get jml data belum proses pengajuan atau pengajuan ataun dalam dinas dan penyelesaian.
// @Description endpoint ini digunakan untuk mendapatkan data belum proses pengajuan atau pengajuan ataun dalam dinas dan penyelesaian.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param idBidang query string false "Set Id Bidang"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/bpd [get]
func (h *DashboardHandler) GetDataDashboardBpd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	idBidang := r.URL.Query().Get("idBidang")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		IdBidang:          idBidang,
	}

	result, err := h.DashboardService.GetDataDashboardBpd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataDashboardSppd jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Summary Get jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Description endpoint ini digunakan untuk mendapatkan data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param idBidang query string false "Set Id Bidang"
// @Param startDate query string false "Set Start Date"
// @Param endDate query string false "Set End Date"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/sppd [get]
func (h *DashboardHandler) GetDataDashboardSppd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	idBidang := r.URL.Query().Get("idBidang")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		IdBidang:          idBidang,
		StartDate:         startDate,
		EndDate:           endDate,
	}

	result, err := h.DashboardService.GetDataDashboardSppd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataDashboardBpdNew jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Summary Get jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Description endpoint ini digunakan untuk mendapatkan data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param idBidang query string false "Set Id Bidang"
// @Param startDate query string false "Set Start Date"
// @Param endDate query string false "Set End Date"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/bpd-new [get]
func (h *DashboardHandler) GetDataDashboardBpdNew(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	idBidang := r.URL.Query().Get("idBidang")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		IdBidang:          idBidang,
		StartDate:         startDate,
		EndDate:           endDate,
	}

	result, err := h.DashboardService.GetDataDashboardBpdNew(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataAktifBpd menampilkan dan Aktif BPD.
// @Summary Get data Aktif BPD.
// @Description endpoint ini digunakan untuk mendapatkan data aktif BPD.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/aktif-bpd [get]
func (h *DashboardHandler) GetDataAktifBpd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
	}

	result, err := h.DashboardService.GetDataBpd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataAktifSppd menampilkan dan Aktif SPPD.
// @Summary Get data Aktif SPPD.
// @Description endpoint ini digunakan untuk mendapatkan data aktif SPPD.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param startDate query string false "Set Start Date"
// @Param endDate query string false "Set End Date"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/aktif-bpd [get]
func (h *DashboardHandler) GetDataAktifSppd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		StartDate:         startDate,
		EndDate:           endDate,
	}

	result, err := h.DashboardService.GetDataSppd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataAktifBpdNew menampilkan dan Aktif Bpd New.
// @Summary Get data Aktif Bpd New.
// @Description endpoint ini digunakan untuk mendapatkan data aktif Bpd New.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param startDate query string false "Set Start Date"
// @Param endDate query string false "Set End Date"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/aktif-bpd-new [get]
func (h *DashboardHandler) GetDataAktifBpdNew(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		StartDate:         startDate,
		EndDate:           endDate,
	}

	result, err := h.DashboardService.GetDataBpdNew(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataDashboardSppd jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Summary Get jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Description endpoint ini digunakan untuk mendapatkan data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param idBidang query string false "Set Id Bidang"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/jumlah-sppd [get]
func (h *DashboardHandler) GetJumlahSppd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	idBidang := r.URL.Query().Get("idBidang")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		IdBidang:          idBidang,
	}

	result, err := h.DashboardService.GetJumlahSppd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// GetDataDashboardSppd jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Summary Get jml data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Description endpoint ini digunakan untuk mendapatkan data belum proses pengajuan atau pengajuan ataun dalam dinas.
// @Tags dashboard
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPegawai query string false "Set Id Pegawai"
// @Param idPegawaiApproval query string false "Set Id Pegawai Approval"
// @Param idBidang query string false "Set Id Bidang"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/dashboard/jumlah-bpd [get]
func (h *DashboardHandler) GetJumlahBpd(w http.ResponseWriter, r *http.Request) {
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	idBidang := r.URL.Query().Get("idBidang")
	req := auth.DashboardRequest{
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		IdBidang:          idBidang,
	}

	result, err := h.DashboardService.GetJumlahBpd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}
