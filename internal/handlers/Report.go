package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/report"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type ReportHandler struct {
	ReportService report.ReportService
	Config        *configs.Config
}

func ProvideReportHandler(service report.ReportService, config *configs.Config) ReportHandler {
	return ReportHandler{
		ReportService: service,
		Config:        config,
	}
}

func (h *ReportHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/report", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/rekap-bpd", h.RptRekapBpd)
			r.Get("/rekap-bpd-bagian", h.RptRekapBpdBagian)
			r.Get("/rekap-total-bpd", h.RptRekapTotalBpd)
			r.Get("/rekap-reimbursement", h.RptRekapReimbursement)
			r.Get("/rekap-akomodasi", h.RptRekapAkomodasi)
			// Export
			r.Get("/export/rekap-bpd", h.ExportRekapBpd)
			r.Get("/export/rekap-bpd-bagian", h.ExportRekapBpdBagian)
			r.Get("/export/rekap-akomodasi-detail", h.ExportRekapAkomodasiDetail)
			r.Get("/export/rekap-reimbursement", h.ExportRekapReimbusment)
			r.Get("/export/rekap-akomodasi", h.ExportRekapAkomodasi)
		})
	})
}

// RptRekapBpd report rekap BPD.
// @Summary Get report rekap BPD.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data report rekap BPD.
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param idUnor query string false "Set ID Unor"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/rekap-bpd [get]
func (h *ReportHandler) RptRekapBpd(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	idUnor := r.URL.Query().Get("idUnor")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
		IdUnor:   idUnor,
	}
	data, err := h.ReportService.RptRekapBpd(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// RptRekapBpd report rekap BPD bagian.
// @Summary Get report rekap BPD bagian.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data report rekap BPD bagian.
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/rekap-bpd-bagian [get]
func (h *ReportHandler) RptRekapBpdBagian(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}
	data, err := h.ReportService.RptRekapBpdBagian(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// RptRekapTotalBpd report rekap BPD total.
// @Summary Get report rekap BPD total.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data report rekap BPD total.
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/rekap-total-bpd [get]
func (h *ReportHandler) RptRekapTotalBpd(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}
	data, err := h.ReportService.RptRekapTotalBpd(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// RptRekapReimbusment report rekap Reimbusment.
// @Summary Get report rekap Reimbusment.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data report rekap BPD.
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/rekap-reimbursement [get]
func (h *ReportHandler) RptRekapReimbursement(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}
	data, err := h.ReportService.RptRekapReimbusment(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// RptRekapAkomodasi report rekap Akomodasi.
// @Summary Get report rekap Akomodasi.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data report rekap BPD.
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/rekap-akomodasi [get]
func (h *ReportHandler) RptRekapAkomodasi(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}
	data, err := h.ReportService.RptRekapAkomodasi(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ExportRekapBpd untuk export Rekap BPD
// @Summary Export Rekap BPD
// @Description End point ini digunakan untuk export Rekap BPD
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/export/rekap-bpd [get]
func (h *ReportHandler) ExportRekapBpd(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}

	file, err := h.ReportService.ExportRekapBpd(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fileName := fmt.Sprintf(`%v.xlsx`, "rekap_bpd")
	contentDisposition := fmt.Sprintf(`attachment; filename=%v`, fileName)

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}

// ExportRekapBpdBagian untuk export Rekap BPD Bagian
// @Summary Export Rekap BPD Bagian
// @Description End point ini digunakan untuk export Rekap BPD Bagian
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/export/rekap-bpd-bagian [get]
func (h *ReportHandler) ExportRekapBpdBagian(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}

	file, err := h.ReportService.ExportRekapBpdBagian(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fileName := fmt.Sprintf(`%v.xlsx`, "rekap_bpd_bagian")
	contentDisposition := fmt.Sprintf(`attachment; filename=%v`, fileName)

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}

// ExportRekapAkomodasiDetail untuk export Rekap Akomodasi Detail
// @Summary Export Rekap Akomodasi Detail
// @Description End point ini digunakan untuk export Rekap Akomodasi Detail
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/export/rekap-akomodasi-detail [get]
func (h *ReportHandler) ExportRekapAkomodasiDetail(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
	}

	file, err := h.ReportService.ExportBiayaAkomodasiDetail(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fileName := fmt.Sprintf(`%v.xlsx`, "biaya_akomodasi_detail")
	contentDisposition := fmt.Sprintf(`attachment; filename=%v`, fileName)

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}

// ExportRekapRe untuk export Rekap BPD
// @Summary Export Rekap BPD
// @Description End point ini digunakan untuk export Rekap BPD
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Param type query string false "Set Type"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/export/rekap-reimbursement [get]
func (h *ReportHandler) ExportRekapReimbusment(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
		Type:     "REIMBURSEMENT",
	}

	file, err := h.ReportService.ExportRekapReimAkm(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fileName := fmt.Sprintf(`%v.xlsx`, "rekap_reimbusment")
	contentDisposition := fmt.Sprintf(`attachment; filename=%v`, fileName)

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}

// ExportRekapRe untuk export Rekap BPD
// @Summary Export Rekap BPD
// @Description End point ini digunakan untuk export Rekap BPD
// @Tags report
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Param tglAwal query string false "Set Tanggal Awal"
// @Param tglAkhir query string false "Set Tanggal Akhir"
// @Param type query string false "Set Type"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/report/export/rekap-akomodasi [get]
func (h *ReportHandler) ExportRekapAkomodasi(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")

	req := report.FilterReport{
		TglAwal:  tglAwal,
		TglAkhir: tglAkhir,
		IdBidang: idBidang,
		Type:     "AKOMODASI",
	}

	file, err := h.ReportService.ExportRekapReimAkm(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	fileName := fmt.Sprintf(`%v.xlsx`, "rekap_reimbusment")
	contentDisposition := fmt.Sprintf(`attachment; filename=%v`, fileName)

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}
