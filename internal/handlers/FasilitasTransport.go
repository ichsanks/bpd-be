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

type FasilitasTransportHandler struct {
	FasilitasTransportService master.FasilitasTransportService
	Config                    *configs.Config
}

func ProvideFasilitasTransportHandler(service master.FasilitasTransportService, config *configs.Config) FasilitasTransportHandler {
	return FasilitasTransportHandler{
		FasilitasTransportService: service,
		Config:                    config,
	}
}

func (h *FasilitasTransportHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/fasilitas-transport", func(r chi.Router) {
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

// ResolveAll list data FasilitasTransport.
// @Summary Get list data FasilitasTransport.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data FasilitasTransport sesuai dengan filter yang dikirimkan.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.FasilitasTransport
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport [get]
func (h *FasilitasTransportHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.FasilitasTransportService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all FasilitasTransport.
// @Summary Get list all FasilitasTransport.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data FasilitasTransport sesuai dengan filter yang dikirimkan.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport/all [get]
func (h *FasilitasTransportHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}
	data, err := h.FasilitasTransportService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createFasilitasTransport adalah untuk menambah data FasilitasTransport.
// @Summary menambahkan data FasilitasTransport.
// @Description Endpoint ini adalah untuk menambahkan data FasilitasTransport.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param FasilitasTransport body master.RequestFasilitasTransport true "FasilitasTransport yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.FasilitasTransport}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport [post]
func (h *FasilitasTransportHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestFasilitasTransport
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

	newData, err := h.FasilitasTransportService.Create(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateFasilitasTransport adalah untuk mengubah data FasilitasTransport.
// @Summary mengubah data FasilitasTransport
// @Description Endpoint ini adalah untuk mengubah data FasilitasTransport.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param FasilitasTransport body master.RequestFasilitasTransport true "FasilitasTransport yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.FasilitasTransport}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport/{id} [put]
func (h *FasilitasTransportHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat master.RequestFasilitasTransport
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

	FasilitasTransport, err := h.FasilitasTransportService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, FasilitasTransport)
}

// ResolveByID adalah untuk mendapatkan satu data FasilitasTransport berdasarkan ID.
// @Summary Mendapatkan satu data FasilitasTransport berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan FasilitasTransport By ID.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.FasilitasTransport}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport/{id} [get]
func (h *FasilitasTransportHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.FasilitasTransportService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data FasilitasTransport.
// @Summary menghapus data FasilitasTransport.
// @Description Endpoint ini adalah untuk menghapus data FasilitasTransport.
// @Tags fasilitas-transport
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/fasilitas-transport/{id} [delete]
func (h *FasilitasTransportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.FasilitasTransportService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
