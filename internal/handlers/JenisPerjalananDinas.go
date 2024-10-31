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

type JenisPerjalananDinasHandler struct {
	JenisPerjalananDinasService master.JenisPerjalananDinasService
	Config                      *configs.Config
}

func ProvideJenisPerjalananDinasHandler(service master.JenisPerjalananDinasService, config *configs.Config) JenisPerjalananDinasHandler {
	return JenisPerjalananDinasHandler{
		JenisPerjalananDinasService: service,
		Config:                      config,
	}
}

func (h *JenisPerjalananDinasHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-perjalanan-dinas", func(r chi.Router) {
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

// ResolveAll list data JenisPerjalananDinas.
// @Summary Get list data JenisPerjalananDinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisPerjalananDinas sesuai dengan filter yang dikirimkan.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.JenisPerjalananDinas
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas [get]
func (h *JenisPerjalananDinasHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.JenisPerjalananDinasService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all JenisPerjalananDinas.
// @Summary Get list all JenisPerjalananDinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisPerjalananDinas sesuai dengan filter yang dikirimkan.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas/all [get]
func (h *JenisPerjalananDinasHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	status, err := h.JenisPerjalananDinasService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createJenisPerjalananDinas adalah untuk menambah data JenisPerjalananDinas.
// @Summary menambahkan data JenisPerjalananDinas.
// @Description Endpoint ini adalah untuk menambahkan data JenisPerjalananDinas.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisPerjalananDinas body master.JenisPerjalananDinasFormat true "JenisPerjalananDinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.JenisPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas [post]
func (h *JenisPerjalananDinasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisPerjalananDinasFormat
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

	newData, err := h.JenisPerjalananDinasService.Create(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateJenisPerjalananDinas adalah untuk merubah data JenisPerjalananDinas.
// @Summary merubah data JenisPerjalananDinas
// @Description Endpoint ini adalah untuk merubah data JenisPerjalananDinas.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisPerjalananDinas body master.JenisPerjalananDinasFormat true "JenisPerjalananDinas yang akan dirubah"
// @Success 200 {object} response.Base{data=master.JenisPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas [put]
func (h *JenisPerjalananDinasHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisPerjalananDinasFormat
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
	newJenisPerjalananDinas, err := h.JenisPerjalananDinasService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newJenisPerjalananDinas)
}

// delete adalah untuk menghapus data JenisPerjalananDinas.
// @Summary menghapus data JenisPerjalananDinas.
// @Description Endpoint ini adalah untuk menghapus data JenisPerjalananDinas.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas/{id} [delete]
func (h *JenisPerjalananDinasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisPerjalananDinas, _ := uuid.FromString(id)
	err := h.JenisPerjalananDinasService.DeleteByID(idJenisPerjalananDinas)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data JenisPerjalananDinas berdasarkan ID.
// @Summary Mendapatkan satu data JenisPerjalananDinas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan JenisPerjalananDinas By ID.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.JenisPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas/{id} [get]
func (h *JenisPerjalananDinasHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.JenisPerjalananDinasService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data JenisPerjalananDinas.
// @Summary menghapus data JenisPerjalananDinas.
// @Description Endpoint ini adalah untuk menghapus data JenisPerjalananDinas.
// @Tags jenisPerjalananDinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-perjalanan-dinas/soft/{id} [delete]
func (h *JenisPerjalananDinasHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisPerjalananDinas, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.JenisPerjalananDinasService.DeleteSoft(idJenisPerjalananDinas, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
