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

type KategoriBiayaHandler struct {
	KategoriBiayaService master.KategoriBiayaService
	Config               *configs.Config
}

func ProvideKategoriBiayaHandler(service master.KategoriBiayaService, config *configs.Config) KategoriBiayaHandler {
	return KategoriBiayaHandler{
		KategoriBiayaService: service,
		Config:               config,
	}
}

func (h *KategoriBiayaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/kategori-biaya", func(r chi.Router) {
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

// ResolveAll list data KategoriBiaya.
// @Summary Get list data KategoriBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data KategoriBiaya sesuai dengan filter yang dikirimkan.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.KategoriBiaya
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya [get]
func (h *KategoriBiayaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.KategoriBiayaService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all KategoriBiaya.
// @Summary Get list all KategoriBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data KategoriBiaya sesuai dengan filter yang dikirimkan.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya/all [get]
func (h *KategoriBiayaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	data, err := h.KategoriBiayaService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createKategoriBiaya adalah untuk menambah data KategoriBiaya.
// @Summary menambahkan data KategoriBiaya.
// @Description Endpoint ini adalah untuk menambahkan data KategoriBiaya.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param KategoriBiaya body master.RequestKategoriBiaya true "KategoriBiaya yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.KategoriBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya [post]
func (h *KategoriBiayaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestKategoriBiaya
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

	newData, err := h.KategoriBiayaService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateKategoriBiaya adalah untuk mengubah data KategoriBiaya.
// @Summary mengubah data KategoriBiaya
// @Description Endpoint ini adalah untuk mengubah data KategoriBiaya.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param KategoriBiaya body master.RequestKategoriBiaya true "KategoriBiaya yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.KategoriBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya/{id} [put]
func (h *KategoriBiayaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestKategoriBiaya
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

	KategoriBiaya, err := h.KategoriBiayaService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, KategoriBiaya)
}

// ResolveByID adalah untuk mendapatkan satu data KategoriBiaya berdasarkan ID.
// @Summary Mendapatkan satu data KategoriBiaya berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan KategoriBiaya By ID.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.KategoriBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya/{id} [get]
func (h *KategoriBiayaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.KategoriBiayaService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data KategoriBiaya.
// @Summary menghapus data KategoriBiaya.
// @Description Endpoint ini adalah untuk menghapus data KategoriBiaya.
// @Tags kategori-biaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kategori-biaya/{id} [delete]
func (h *KategoriBiayaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.KategoriBiayaService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
