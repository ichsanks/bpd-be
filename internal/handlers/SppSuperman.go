package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type SppSupermanHandler struct {
	SppSupermanService     bpd.SppSupermanService
	PerjalananDinasService bpd.PerjalananDinasService
}

func ProvideSppSupermanHandler(service bpd.SppSupermanService, perjalananDinas bpd.PerjalananDinasService) SppSupermanHandler {
	return SppSupermanHandler{
		SppSupermanService:     service,
		PerjalananDinasService: perjalananDinas,
	}
}

func (h *SppSupermanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/superman/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/get-vendor", h.GetAllDataVendor)
			r.Get("/get-gl", h.GetAllDataGL)
			r.Get("/get-customer", h.GetAllDataCustomer)
			r.Get("/get-cost-center", h.GetAllDataCostCenter)
			r.Get("/get-profit-center", h.GetAllDataProfitCenter)
			r.Get("/get-cash-flow", h.GetAllDataCashFlow)
			r.Get("/get-sumber-dana", h.GetAllDataSumberDana)
			r.Get("/get-bagian", h.GetAllDataBagian)
			r.Get("/get-ino-karyawan", h.GetAllDataInoKaryawan)
			r.Get("/get-nomor-urut", h.GetAllDataNomorUrut)
			r.Get("/get-spp", h.GetDataListSpp)
			r.Post("/create-supermen", h.Create)
			r.Post("/create-supermen-new", h.CreateNew)
			r.Get("/get-rekam-jejak", h.GetDataRekamJejak)
			r.Get("/get-detail-spp", h.GetDetailSpp)
			r.Get("/get-spp-id", h.GetSppId)
		})
	})
}

// GetDataAll list all vendor.
// @Summary Get list all vendor.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-vendor [get]
func (h *SppSupermanHandler) GetAllDataVendor(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataVendor()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all GL.
// @Summary Get list all GL.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-gl [get]
func (h *SppSupermanHandler) GetAllDataGL(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataGL()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Customer.
// @Summary Get list all Customer.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-customer [get]
func (h *SppSupermanHandler) GetAllDataCustomer(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataCustomer()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all CostCenter.
// @Summary Get list all CostCenter.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-cost-center [get]
func (h *SppSupermanHandler) GetAllDataCostCenter(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataCostCenter()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all profit center.
// @Summary Get list all profit center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-profit-center [get]
func (h *SppSupermanHandler) GetAllDataProfitCenter(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataProfitCenter()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all cashflow center.
// @Summary Get list all cashflow center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-cash-flow [get]
func (h *SppSupermanHandler) GetAllDataCashFlow(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataCashFlow()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all sumberdana center.
// @Summary Get list all sumberdana center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-sumber-dana [get]
func (h *SppSupermanHandler) GetAllDataSumberDana(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataSumberDana()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all bagian center.
// @Summary Get list all bagian center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-bagian [get]
func (h *SppSupermanHandler) GetAllDataBagian(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataBagian()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Ino Karyawan center.
// @Summary Get list all Ino Karyawan center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param nik query string false "Set nik"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-ino-karyawan [get]
func (h *SppSupermanHandler) GetAllDataInoKaryawan(w http.ResponseWriter, r *http.Request) {
	nik := r.URL.Query().Get("nik")
	status, err := h.SppSupermanService.GetDataInoKaryawan(nik)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list nomor urut center.
// @Summary Get list nomor urut center.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-nomor-urut [get]
func (h *SppSupermanHandler) GetAllDataNomorUrut(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataNomorUrut()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createSPPB adalah untuk menambah data SPP.
// @Summary menambahkan data SPP.
// @Description Endpoint ini adalah untuk menambahkan data SPP.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Golongan body bpd.RequestSppb true "sppb yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.DataSppb}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/create-supermen [post]
func (h *SppSupermanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.RequestSppb
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newData, err := h.SppSupermanService.Create(reqFormat)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// GetDataListSpp list all list spp.
// @Summary Get list all list spp.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data list spp filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-spp [get]
func (h *SppSupermanHandler) GetDataListSpp(w http.ResponseWriter, r *http.Request) {
	data, err := h.SppSupermanService.GetDataListSpp()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createSPPB adalah untuk menambah data SPP.
// @Summary menambahkan data SPP.
// @Description Endpoint ini adalah untuk menambahkan data SPP.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SupermanParam body bpd.RequestPayloadUraian true "request yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.ResponseDataSuperman}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/create-supermen-new [post]
func (h *SppSupermanHandler) CreateNew(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.RequestPayloadUraian
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newData, err := h.SppSupermanService.CreateNew(reqFormat)

	if err != nil {
		response.WithError(w, err)
		return
	}

	req := bpd.ResponseSpp{
		ID:    *reqFormat.IdBpd,
		SppId: newData.DataSpp.SppID,
	}
	_, err = h.PerjalananDinasService.UpdateSppBpd(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// GetDataRekamJejak list rekam jejak.
// @Summary Get rekam jejak
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-rekam-jejak [get]
func (h *SppSupermanHandler) GetDataRekamJejak(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDataRekamJejak()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDetailSpp list nomor urut center.
// @Summary Get detail spp
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-detail-spp [get]
func (h *SppSupermanHandler) GetDetailSpp(w http.ResponseWriter, r *http.Request) {
	status, err := h.SppSupermanService.GetDetailSpp()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataSppId list Spp.
// @Summary Get list Spp.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data superman filter yang dikirimkan.
// @Tags api-superman
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param sppId query string false "Set SPPid"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/superman/get-spp-id [get]
func (h *SppSupermanHandler) GetSppId(w http.ResponseWriter, r *http.Request) {
	sppid := r.URL.Query().Get("sppId")
	id, err := strconv.Atoi(sppid)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	status, err := h.SppSupermanService.GetSppId(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
