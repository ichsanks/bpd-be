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

type UnitKerjaHandler struct {
	UnitKerjaService master.UnitKerjaService
	Config           *configs.Config
}

func ProvideUnitKerjaHandler(service master.UnitKerjaService, config *configs.Config) UnitKerjaHandler {
	return UnitKerjaHandler{
		UnitKerjaService: service,
		Config:           config,
	}
}

func (h *UnitKerjaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/unit-kerja", func(r chi.Router) {
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

// ResolveAll list data unit kerja.
// @Summary Get list data unit kerja.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data unit kerja sesuai dengan filter yang dikirimkan.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idBidang query string false "id bidang"
// @Success 200 {object} master.UnitOrganisasiKerja
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja [get]
func (h *UnitKerjaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "ASC"
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

	idBidang := r.URL.Query().Get("idBidang")

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
		IdBidang:   idBidang,
		IdBranch:   idBranchs,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.UnitKerjaService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all unit kerja.
// @Summary Get list all unit kerja.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data unit kerja sesuai dengan filter yang dikirimkan.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBidang query string false "Set ID Bidang"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja/all [get]
func (h *UnitKerjaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	idBidang := r.URL.Query().Get("idBidang")
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
		IdBidang: idBidang,
	}

	data, err := h.UnitKerjaService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createUnitKerja adalah untuk menambah data Unit Kerja.
// @Summary menambahkan data Unit Kerja.
// @Description Endpoint ini adalah untuk menambahkan data Unit Kerja.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Unor body master.RequestUnor true "Unor yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.UnitOrganisasiKerja}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja [post]
func (h *UnitKerjaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestUnor
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

	newData, err := h.UnitKerjaService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateUnitKerja adalah untuk mengubah data Unit Kerja.
// @Summary mengubah data Unit Kerja
// @Description Endpoint ini adalah untuk mengubah data Unit Kerja.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param Bidang body master.RequestUnor true "Unit Kerja yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.UnitOrganisasiKerja}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja/{id} [put]
func (h *UnitKerjaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestUnor
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

	bidang, err := h.UnitKerjaService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, bidang)
}

// ResolveByID adalah untuk mendapatkan satu data unit kerja berdasarkan ID.
// @Summary Mendapatkan satu data unit kerja berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan unit kerja By ID.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.UnitOrganisasiKerja}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja/{id} [get]
func (h *UnitKerjaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unor, err := h.UnitKerjaService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unor)
}

// delete adalah untuk menghapus data unit kerja.
// @Summary menghapus data unit kerja.
// @Description Endpoint ini adalah untuk menghapus data unit kerja.
// @Tags unitKerja
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/unit-kerja/{id} [delete]
func (h *UnitKerjaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.UnitKerjaService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
