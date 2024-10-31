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

type JabatanHandler struct {
	JabatanService master.JabatanService
	Config         *configs.Config
}

func ProvideJabatanHandler(service master.JabatanService, config *configs.Config) JabatanHandler {
	return JabatanHandler{
		JabatanService: service,
		Config:         config,
	}
}

func (h *JabatanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jabatan", func(r chi.Router) {
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

// ResolveAll list data Jabatan.
// @Summary Get list data Jabatan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Jabatan sesuai dengan filter yang dikirimkan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.Jabatan
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan [get]
func (h *JabatanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.JabatanService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Jabatan.
// @Summary Get list all Jabatan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Jabatan sesuai dengan filter yang dikirimkan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan/all [get]
func (h *JabatanHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	status, err := h.JabatanService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createJabatan adalah untuk menambah data Jabatan.
// @Summary menambahkan data Jabatan.
// @Description Endpoint ini adalah untuk menambahkan data Jabatan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Jabatan body master.JabatanFormat true "Jabatan yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.Jabatan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan [post]
func (h *JabatanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JabatanFormat
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

	newData, err := h.JabatanService.Create(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateJabatan adalah untuk merubah data Jabatan.
// @Summary merubah data Jabatan
// @Description Endpoint ini adalah untuk merubah data Jabatan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Jabatan body master.JabatanFormat true "Jabatan yang akan dirubah"
// @Success 200 {object} response.Base{data=master.Jabatan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan [put]
func (h *JabatanHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JabatanFormat
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
	newJabatan, err := h.JabatanService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newJabatan)
}

// delete adalah untuk menghapus data Jabatan.
// @Summary menghapus data Jabatan.
// @Description Endpoint ini adalah untuk menghapus data Jabatan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan/{id} [delete]
func (h *JabatanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJabatan, _ := uuid.FromString(id)
	err := h.JabatanService.DeleteByID(idJabatan)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data Jabatan berdasarkan ID.
// @Summary Mendapatkan satu data Jabatan berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Jabatan By ID.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.Jabatan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan/{id} [get]
func (h *JabatanHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.JabatanService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data Jabatan.
// @Summary menghapus data Jabatan.
// @Description Endpoint ini adalah untuk menghapus data Jabatan.
// @Tags jabatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jabatan/soft/{id} [delete]
func (h *JabatanHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJabatan, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.JabatanService.DeleteSoft(idJabatan, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
