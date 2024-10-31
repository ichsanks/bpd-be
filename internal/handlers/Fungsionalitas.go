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

type FungsionalitasHandler struct {
	FungsionalitasService master.FungsionalitasService
	Config                *configs.Config
}

func ProvideFungsionalitasHandler(service master.FungsionalitasService, config *configs.Config) FungsionalitasHandler {
	return FungsionalitasHandler{
		FungsionalitasService: service,
		Config:                config,
	}
}

func (h *FungsionalitasHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/fungsionalitas", func(r chi.Router) {
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

// ResolveAll list data Fungsionalitas.
// @Summary Get list data Fungsionalitas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Fungsionalitas sesuai dengan filter yang dikirimkan.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.Fungsionalitas
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas [get]
func (h *FungsionalitasHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.FungsionalitasService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all fungsionalitas.
// @Summary Get list all fungsionalitas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data fungsionalitas sesuai dengan filter yang dikirimkan.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas/all [get]
func (h *FungsionalitasHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := h.FungsionalitasService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createFungsionalitas adalah untuk menambah data fungsionalitas.
// @Summary menambahkan data fungsionalitas.
// @Description Endpoint ini adalah untuk menambahkan data fungsionalitas.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Fungsionalitas body master.RequestFungsionalitas true "fungsionalitas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.Fungsionalitas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas [post]
func (h *FungsionalitasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestFungsionalitas
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

	newData, err := h.FungsionalitasService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateFungsionalitas adalah untuk mengubah data fungsionalitas.
// @Summary mengubah data fungsionalitas
// @Description Endpoint ini adalah untuk mengubah data fungsionalitas.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param fungsionalitas body master.RequestFungsionalitas true "Fungsionalitas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.Fungsionalitas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas/{id} [put]
func (h *FungsionalitasHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestFungsionalitas
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

	data, err := h.FungsionalitasService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID adalah untuk mendapatkan satu data fungsionalitas berdasarkan ID.
// @Summary Mendapatkan satu data fungsionalitas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan fungsionalitas By ID.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.Fungsionalitas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas/{id} [get]
func (h *FungsionalitasHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.FungsionalitasService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data fungsionalitas.
// @Summary menghapus data fungsionalitas.
// @Description Endpoint ini adalah untuk menghapus data fungsionalitas.
// @Tags fungsionalitas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fungsionalitas/{id} [delete]
func (h *FungsionalitasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.FungsionalitasService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
