package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type SyaratDokumenHandler struct {
	SyaratDokumenService master.SyaratDokumenService
	Config               *configs.Config
}

func ProvideSyaratDokumenHandler(service master.SyaratDokumenService, config *configs.Config) SyaratDokumenHandler {
	return SyaratDokumenHandler{
		SyaratDokumenService: service,
		Config:               config,
	}
}

func (h *SyaratDokumenHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/syarat-dokumen", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Get("/{id}", h.ResolveByID)
		})
	})
}

// ResolveAll list data SyaratDokumen.
// @Summary Get list data SyaratDokumen.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data SyaratDokumen sesuai dengan filter yang dikirimkan.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.SyaratDokumen
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen [get]
func (h *SyaratDokumenHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "DESC"
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		IdBranch:   idBranchs,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.SyaratDokumenService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Syarat Dokumen.
// @Summary Get list all Syarat Dokumen.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data SyaratDokumen sesuai dengan filter yang dikirimkan.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idTransaksi query string false "jenis transaksi"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen/all [get]
func (h *SyaratDokumenHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	idTransaksi := r.URL.Query().Get("idTransaksi")
	req := model.StandardRequest{
		IdBranch:    idBranchs,
		IdTransaksi: idTransaksi,
	}
	data, err := h.SyaratDokumenService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createSyaratDokumen adalah untuk menambah data SyaratDokumen.
// @Summary menambahkan data SyaratDokumen.
// @Description Endpoint ini adalah untuk menambahkan data SyaratDokumen.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SyaratDokumen body master.RequestSyaratDokumen true "Dokumen yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.SyaratDokumen}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen [post]
func (h *SyaratDokumenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestSyaratDokumen
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
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

	newData, err := h.SyaratDokumenService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateDokumen adalah untuk mengubah data SyaratDokumen.
// @Summary mengubah data SyaratDokumen
// @Description Endpoint ini adalah untuk mengubah data SyaratDokumen.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param SyaratDokumen body master.RequestSyaratDokumen true "Dokumen yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.SyaratDokumen}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen/{id} [put]
func (h *SyaratDokumenHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestSyaratDokumen
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.ID = id
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
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

	Dokumen, err := h.SyaratDokumenService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, Dokumen)
}

// ResolveByID adalah untuk mendapatkan satu data SyaratDokumen berdasarkan ID.
// @Summary Mendapatkan satu data SyaratDokumen berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Dokumen By ID.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.SyaratDokumen}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen/{id} [get]
func (h *SyaratDokumenHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.SyaratDokumenService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data SyaratDokumen.
// @Summary menghapus data SyaratDokumen.
// @Description Endpoint ini adalah untuk menghapus data SyaratDokumen.
// @Tags syarat-dokumen
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/syarat-dokumen/{id} [delete]
func (h *SyaratDokumenHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.SyaratDokumenService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
