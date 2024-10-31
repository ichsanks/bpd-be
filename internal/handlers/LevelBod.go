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

type LevelBodHandler struct {
	LevelBodService master.LevelBodService
	Config          *configs.Config
}

func ProvideLevelBodHandler(service master.LevelBodService, config *configs.Config) LevelBodHandler {
	return LevelBodHandler{
		LevelBodService: service,
		Config:          config,
	}
}

func (h *LevelBodHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/level-bod", func(r chi.Router) {
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

// ResolveAll list data LevelBod.
// @Summary Get list data LevelBod.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data LevelBod sesuai dengan filter yang dikirimkan.
// @Tags level-bod
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.LevelBod
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod [get]
func (h *LevelBodHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.LevelBodService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all LevelBod.
// @Summary Get list all LevelBod.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data LevelBod sesuai dengan filter yang dikirimkan.
// @Tags level-bod
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod/all [get]
func (h *LevelBodHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	data, err := h.LevelBodService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createLevelBod adalah untuk menambah data LevelBod.
// @Summary menambahkan data LevelBod.
// @Description Endpoint ini adalah untuk menambahkan data LevelBod.
// @Tags level-bod
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param LevelBod body master.RequestLevelBod true "LevelBod yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.LevelBod}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod [post]
func (h *LevelBodHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestLevelBod
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

	newData, err := h.LevelBodService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateLevelBod adalah untuk mengubah data LevelBod.
// @Summary mengubah data LevelBod
// @Description Endpoint ini adalah untuk mengubah data LevelBod.
// @Tags level-bod
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param LevelBod body master.RequestLevelBod true "LevelBod yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.LevelBod}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod/{id} [put]
func (h *LevelBodHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestLevelBod
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

	LevelBod, err := h.LevelBodService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, LevelBod)
}

// ResolveByID adalah untuk mendapatkan satu data LevelBod berdasarkan ID.
// @Summary Mendapatkan satu data LevelBod berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan LevelBod By ID.
// @Tags level-bod
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.LevelBod}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod/{id} [get]
func (h *LevelBodHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.LevelBodService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data LevelBod.
// @Summary menghapus data LevelBod.
// @Description Endpoint ini adalah untuk menghapus data LevelBod.
// @Tags level-bod
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/level-bod/{id} [delete]
func (h *LevelBodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.LevelBodService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
