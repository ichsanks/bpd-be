package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PengajuanBpdHistoriHandler struct {
	PengajuanBpdHistoriService bpd.PengajuanBpdHistoriService
}

func ProvidePengajuanBpdHistoriHandler(service bpd.PengajuanBpdHistoriService) PengajuanBpdHistoriHandler {
	return PengajuanBpdHistoriHandler{
		PengajuanBpdHistoriService: service,
	}
}

func (h *PengajuanBpdHistoriHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/pengajuan-bpd-histori", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Post("/penyelesaian", h.CreatePenyelesaian)
			r.Post("/approval", h.Approve)
			r.Post("/batal", h.Batal)
			r.Post("/revisi-biaya", h.RevisiPenyelesaianBiaya)
			r.Get("/timeline", h.GetTimeline)
			r.Get("/timeline-ttd", h.GetTimelineTtd)
			r.Post("/pengajuan-revisi", h.PengajuanRevisi)
		})
	})
}

// create adalah untuk menambah data pengajuan bpd histori.
// @Summary menambahkan data pengajuan bpd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan bpd histori.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanBpdHistoriInputRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanBpdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori [post]
func (h *PengajuanBpdHistoriHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanBpdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanBpdHistoriService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// CreatePenyelesaian adalah untuk menambah data penyelesaian bpd histori.
// @Summary menambahkan data penyelesaian bpd histori.
// @Description Endpoint ini adalah untuk menambahkan data penyelesaian bpd histori.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanBpdHistoriInputRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanBpdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/penyelesaian [post]
func (h *PengajuanBpdHistoriHandler) CreatePenyelesaian(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanBpdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TypeApproval = "PENYELESAIAN"
	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanBpdHistoriService.CreatePenyelesaian(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// approve adalah untuk menambah data pengajuan bpd histori.
// @Summary menambahkan data pengajuan bpd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan bpd histori.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanBpdHistoriApproveRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanBpdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/approval [post]
func (h *PengajuanBpdHistoriHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanBpdHistoriApproveRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanBpdHistoriInputRequest{
		ID:           reqFormat.ID,
		IdPegawai:    reqFormat.IdPegawai,
		Catatan:      reqFormat.Catatan,
		Keterangan:   reqFormat.Keterangan,
		Status:       reqFormat.Status,
		TypeApproval: reqFormat.TypeApproval,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanBpdHistoriService.Approve(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// GetTimeline list all timeline pengajuan bpd histori.
// @Summary Get list all timeline pengajuan bpd histori.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data timeline pengajuan bpd histori sesuai dengan filter yang dikirimkan.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPerjalananDinas query string false "Set ID Perjalanan Dinas"
// @Param idBpdPegawai query string false "Set ID BPD Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/timeline [get]
func (h *PengajuanBpdHistoriHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	idPerjalananDinas := r.URL.Query().Get("idPerjalananDinas")
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.PengajuanBpdHistoriService.GetTimeline(idPerjalananDinas, idBpdPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// GetTimelineTtd list all timeline pengajuan bpd histori.
// @Summary Get list all timelinettd pengajuan bpd histori.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data timelinettd pengajuan bpd histori sesuai dengan filter yang dikirimkan.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPerjalananDinas query string false "Set ID Perjalanan Dinas"
// @Param idBpdPegawai query string false "Set ID BPD Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/timeline-ttd [get]
func (h *PengajuanBpdHistoriHandler) GetTimelineTtd(w http.ResponseWriter, r *http.Request) {
	idPerjalananDinas := r.URL.Query().Get("idPerjalananDinas")
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.PengajuanBpdHistoriService.GetTimelineTtd(idPerjalananDinas, idBpdPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// batal adalah untuk membatalkan data pengajuan bpd.
// @Summary membatalkan data pengajuan bpd.
// @Description Endpoint ini adalah untuk membatalkan data pengajuan bpd.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.BatalBpdRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/batal [post]
func (h *PengajuanBpdHistoriHandler) Batal(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.BatalBpdRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanBpdHistoriInputRequest{
		IdPerjalananDinas: reqFormat.ID,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	err = h.PengajuanBpdHistoriService.Batal(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, "success")
}

// RevisiPenyelesaianBiaya adalah untuk revisi biaya data pengajuan bpd histori.
// @Summary revisi biaya data pengajuan bpd histori.
// @Description Endpoint ini adalah untuk revisi biaya data pengajuan bpd histori.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Data body bpd.PengajuanBpdHistoriApproveRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanBpdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/revisi-biaya [post]
func (h *PengajuanBpdHistoriHandler) RevisiPenyelesaianBiaya(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanBpdHistoriApproveRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := bpd.PengajuanBpdHistoriInputRequest{
		ID:        reqFormat.ID,
		IdPegawai: reqFormat.IdPegawai,
		Status:    reqFormat.Status,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanBpdHistoriService.RevisiPenyelesaianBiaya(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// createPengajuanRevisi adalah untuk menambah data pengajuan bpd histori.
// @Summary menambahkan data pengajuan bpd histori.
// @Description Endpoint ini adalah untuk menambahkan data pengajuan bpd histori.
// @Tags pengajuan-bpd-histori
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PengajuanBpdHistoriInputRequest true "Pengajuan Bpd Histori yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PengajuanBpdHistori}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/pengajuan-bpd-histori/pengajuan-revisi [post]
func (h *PengajuanBpdHistoriHandler) PengajuanRevisi(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PengajuanBpdHistoriInputRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PengajuanBpdHistoriService.PengajuanRevisi(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}
