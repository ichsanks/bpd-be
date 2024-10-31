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

type StatusPegawaiHandler struct {
	StatusPegawaiService master.StatusPegawaiService
	Config               *configs.Config
}

func ProvideStatusPegawaiHandler(service master.StatusPegawaiService, config *configs.Config) StatusPegawaiHandler {
	return StatusPegawaiHandler{
		StatusPegawaiService: service,
		Config:               config,
	}
}

func (h *StatusPegawaiHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/status-pegawai", func(r chi.Router) {
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

// ResolveAll list data StatusPegawai.
// @Summary Get list data StatusPegawai.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data StatusPegawai sesuai dengan filter yang dikirimkan.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.StatusPegawai
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai [get]
func (h *StatusPegawaiHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.StatusPegawaiService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all StatusPegawai.
// @Summary Get list all StatusPegawai.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data StatusPegawai sesuai dengan filter yang dikirimkan.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai/all [get]
func (h *StatusPegawaiHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	data, err := h.StatusPegawaiService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createStatusPegawai adalah untuk menambah data StatusPegawai.
// @Summary menambahkan data StatusPegawai.
// @Description Endpoint ini adalah untuk menambahkan data StatusPegawai.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param StatusPegawai body master.RequestStatusPegawai true "StatusPegawai yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.StatusPegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai [post]
func (h *StatusPegawaiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestStatusPegawai
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

	newData, err := h.StatusPegawaiService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateStatusPegawai adalah untuk mengubah data StatusPegawai.
// @Summary mengubah data StatusPegawai
// @Description Endpoint ini adalah untuk mengubah data StatusPegawai.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param StatusPegawai body master.RequestStatusPegawai true "StatusPegawai yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.StatusPegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai/{id} [put]
func (h *StatusPegawaiHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestStatusPegawai
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

	StatusPegawai, err := h.StatusPegawaiService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, StatusPegawai)
}

// ResolveByID adalah untuk mendapatkan satu data StatusPegawai berdasarkan ID.
// @Summary Mendapatkan satu data StatusPegawai berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan StatusPegawai By ID.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.StatusPegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai/{id} [get]
func (h *StatusPegawaiHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.StatusPegawaiService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data StatusPegawai.
// @Summary menghapus data StatusPegawai.
// @Description Endpoint ini adalah untuk menghapus data StatusPegawai.
// @Tags status-pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/status-pegawai/{id} [delete]
func (h *StatusPegawaiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.StatusPegawaiService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
