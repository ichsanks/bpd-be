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

type SettingBiayaHandler struct {
	SettingBiayaService master.SettingBiayaService
	Config              *configs.Config
}

func ProvideSettingBiayaHandler(service master.SettingBiayaService, config *configs.Config) SettingBiayaHandler {
	return SettingBiayaHandler{
		SettingBiayaService: service,
		Config:              config,
	}
}

func (h *SettingBiayaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/setting-biaya", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Put("/", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Get("/{id}", h.ResolveByID)
			r.Delete("/soft/{id}", h.DeleteSoft)
		})
	})
}

// ResolveAll list data SettingBiaya.
// @Summary Get list data SettingBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data SettingBiaya sesuai dengan filter yang dikirimkan.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idBodLevel query string false "idBodLevel"
// @Param idJenisTujuan query string false "idJenisTujuan"
// @Param idKategoriBiaya query string false "idKategoriBiaya"
// @Success 200 {object} master.SettingBiaya
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya [get]
func (h *SettingBiayaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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
	idKategoriBiaya := r.URL.Query().Get("idKategoriBiaya")
	idJenisTujuan := r.URL.Query().Get("idJenisTujuan")
	idBodLevel := r.URL.Query().Get("idBodLevel")

	req := model.StandardRequest{
		Keyword:       keyword,
		PageSize:      pageSize,
		PageNumber:    pageNumber,
		SortBy:        sortBy,
		SortType:      sortType,
		IdBranch:      idBranchs,
		IdTransaksi:   idKategoriBiaya,
		IdBodLevel:    idBodLevel,
		IdJenisTujuan: idJenisTujuan,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.SettingBiayaService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all SettingBiaya.
// @Summary Get list all SettingBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data SettingBiaya sesuai dengan filter yang dikirimkan.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya/all [get]
func (h *SettingBiayaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	status, err := h.SettingBiayaService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createSettingBiaya adalah untuk menambah data SettingBiaya.
// @Summary menambahkan data SettingBiaya.
// @Description Endpoint ini adalah untuk menambahkan data SettingBiaya.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SettingBiaya body master.SettingBiayaFormat true "SettingBiaya yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.SettingBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya [post]
func (h *SettingBiayaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.SettingBiayaFormat
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TenantID = &tenantID

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if reqFormat.IdBranch == nil {
		reqFormat.IdBranch = &idBranch
	}

	newData, err := h.SettingBiayaService.Create(reqFormat, userID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateSettingBiaya adalah untuk merubah data SettingBiaya.
// @Summary merubah data SettingBiaya
// @Description Endpoint ini adalah untuk merubah data SettingBiaya.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SettingBiaya body master.SettingBiayaUpdateFormat true "SettingBiaya yang akan dirubah"
// @Success 200 {object} response.Base{data=master.SettingBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya [put]
func (h *SettingBiayaHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.SettingBiayaUpdateFormat
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TenantID = &tenantID

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if reqFormat.IdBranch == nil {
		reqFormat.IdBranch = &idBranch
	}

	newSettingBiaya, err := h.SettingBiayaService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newSettingBiaya)
}

// delete adalah untuk menghapus data SettingBiaya.
// @Summary menghapus data SettingBiaya.
// @Description Endpoint ini adalah untuk menghapus data SettingBiaya.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya/{id} [delete]
func (h *SettingBiayaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idSettingBiaya, _ := uuid.FromString(id)
	err := h.SettingBiayaService.DeleteByID(idSettingBiaya)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data SettingBiaya berdasarkan ID.
// @Summary Mendapatkan satu data SettingBiaya berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan SettingBiaya By ID.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.SettingBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya/{id} [get]
func (h *SettingBiayaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.SettingBiayaService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data SettingBiaya.
// @Summary menghapus data SettingBiaya.
// @Description Endpoint ini adalah untuk menghapus data SettingBiaya.
// @Tags settingBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/setting-biaya/soft/{id} [delete]
func (h *SettingBiayaHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idSettingBiaya, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.SettingBiayaService.DeleteSoft(idSettingBiaya, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
