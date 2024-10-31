package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PengajuanSppdHistoriHandler struct {
	PengajuanSppdHistoriService bpd.PengajuanSppdHistoriService
}

func ProvidePengajuanSppdHistoriHandler(service bpd.PengajuanSppdHistoriService) PengajuanSppdHistoriHandler {
	return PengajuanSppdHistoriHandler{
		PengajuanSppdHistoriService: service,
	}
}

func (h *PengajuanSppdHistoriHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/pengajuan-sppd-histori", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Post("/penyelesaian", h.CreatePenyelesaian)
			r.Post("/approval", h.Approve)
			r.Post("/batal", h.Batal)
			r.Post("/revisi-biaya", h.RevisiPenyelesaianBiaya)
			r.Get("/timeline", h.GetTimeline)
			r.Post("/pengajuan-revisi", h.PengajuanRevisi)
		})
	})
}

// create adalah untuk menambah data pengajuan sppd histori.
// @Summary menambahkan data pengajuan sppd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan sppd histori.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanSppdHistoriInputRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanSppdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori [post]
func (h *PengajuanSppdHistoriHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanSppdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID := middleware.GetClaimsValue(r.Context(), "tenantId").(string)
	if err != nil {
		fmt.Print("error tenantId")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if reqFormat.IdBranch == nil {
		reqFormat.IdBranch = &idBranch
	}

	if reqFormat.TenantId == nil {
		reqFormat.TenantId = &tenantID
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanSppdHistoriService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// CreatePenyelesaian adalah untuk menambah data penyelesaian sppd histori.
// @Summary menambahkan data penyelesaian sppd histori.
// @Description Endpoint ini adalah untuk menambahkan data penyelesaian sppd histori.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanSppdHistoriInputRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanSppdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/penyelesaian [post]
func (h *PengajuanSppdHistoriHandler) CreatePenyelesaian(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanSppdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TypeApproval = "PENYELESAIAN"
	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanSppdHistoriService.CreatePenyelesaian(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// approve adalah untuk menambah data pengajuan sppd histori.
// @Summary menambahkan data pengajuan sppd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan sppd histori.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanSppdHistoriApproveRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanSppdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/approval [post]
func (h *PengajuanSppdHistoriHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanSppdHistoriApproveRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanSppdHistoriInputRequest{
		ID:           reqFormat.ID,
		IdPegawai:    reqFormat.IdPegawai,
		Catatan:      reqFormat.Catatan,
		Keterangan:   reqFormat.Keterangan,
		Status:       reqFormat.Status,
		TypeApproval: reqFormat.TypeApproval,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanSppdHistoriService.Approve(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// GetTimeline list all timeline pengajuan sppd histori.
// @Summary Get list all timeline pengajuan sppd histori.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data timeline pengajuan sppd histori sesuai dengan filter yang dikirimkan.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idSuratPerjalananDinas query string false "Set ID Perjalanan Dinas"
// @Param idBpdPegawai query string false "Set ID BPD Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/timeline [get]
func (h *PengajuanSppdHistoriHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	IdSuratPerjalananDinas := r.URL.Query().Get("idSuratPerjalananDinas")
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.PengajuanSppdHistoriService.GetTimeline(IdSuratPerjalananDinas, idBpdPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// batal adalah untuk membatalkan data pengajuan bpd.
// @Summary membatalkan data pengajuan bpd.
// @Description Endpoint ini adalah untuk membatalkan data pengajuan bpd.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.BatalBpdRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/batal [post]
func (h *PengajuanSppdHistoriHandler) Batal(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.BatalBpdRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanSppdHistoriInputRequest{
		IdSuratPerjalananDinas: reqFormat.ID,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	err = h.PengajuanSppdHistoriService.Batal(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, "success")
}

// RevisiPenyelesaianBiaya adalah untuk revisi biaya data pengajuan sppd histori.
// @Summary revisi biaya data pengajuan sppd histori.
// @Description Endpoint ini adalah untuk revisi biaya data pengajuan sppd histori.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Data body bpd.PengajuanSppdHistoriApproveRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanSppdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/revisi-biaya [post]
func (h *PengajuanSppdHistoriHandler) RevisiPenyelesaianBiaya(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanSppdHistoriApproveRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanSppdHistoriInputRequest{
		ID:        reqFormat.ID,
		IdPegawai: reqFormat.IdPegawai,
		Status:    reqFormat.Status,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanSppdHistoriService.RevisiPenyelesaianBiaya(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// pengajuanRevisi adalah untuk menambah data pengajuan sppd histori.
// @Summary menambahkan data pengajuan sppd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan sppd histori.
// @Tags pengajuan-sppd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanSppdHistoriInputRequest true "Pengajuan sppd histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanSppdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-sppd-histori/pengajuan-revisi [post]
func (h *PengajuanSppdHistoriHandler) PengajuanRevisi(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanSppdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID := middleware.GetClaimsValue(r.Context(), "tenantId").(string)
	if err != nil {
		fmt.Print("error tenantId")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if reqFormat.IdBranch == nil {
		reqFormat.IdBranch = &idBranch
	}

	if reqFormat.TenantId == nil {
		reqFormat.TenantId = &tenantID
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanSppdHistoriService.RevisiPengajuan(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}
