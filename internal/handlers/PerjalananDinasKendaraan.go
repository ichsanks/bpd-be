package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PerjalananDinasKendaraanHandler struct {
	PerjalananDinasKendaraanService bpd.PerjalananDinasKendaraanService
	Config                          *configs.Config
}

func ProvidePerjalananDinasKendaraanHandler(service bpd.PerjalananDinasKendaraanService, config *configs.Config) PerjalananDinasKendaraanHandler {
	return PerjalananDinasKendaraanHandler{
		PerjalananDinasKendaraanService: service,
		Config:                          config,
	}
}

func (h *PerjalananDinasKendaraanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/kendaraan-perjalanan-dinas", func(r chi.Router) {
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

// ResolveAll list data Perjalanan dinas kendaraan.
// @Summary Get list data Perjalanan dinas kendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan dinas kendaraan sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idPerjalananDinas query string false "Set ID Perjalanan Dinas"
// @Success 200 {object} master.Bidang
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas [get]
func (h *PerjalananDinasKendaraanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	idPerjalananDinas := r.URL.Query().Get("idPerjalananDinas")
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
		Keyword:           keyword,
		PageSize:          pageSize,
		PageNumber:        pageNumber,
		SortBy:            sortBy,
		SortType:          sortType,
		IdPerjalananDinas: idPerjalananDinas,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.PerjalananDinasKendaraanService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all perjalanan dinas kendaraan.
// @Summary Get list all perjalanan dinas kendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data perjalanan dinas kendaraan sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPerjalananDinas query string true "Set ID Perjalanan Dinas"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas/all [get]
func (h *PerjalananDinasKendaraanHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	idPerjalananDinas := r.URL.Query().Get("idPerjalananDinas")
	data, err := h.PerjalananDinasKendaraanService.GetAll(idPerjalananDinas)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data perjalanan dinas kendaraan.
// @Summary menambahkan data perjalanan dinas kendaraan.
// @Description Endpoint ini adalah untuk menambahkan data perjalanan dinas kendaraan.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Data body bpd.RequestPerjalananDinasKendaraan true "Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas [post]
func (h *PerjalananDinasKendaraanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.RequestPerjalananDinasKendaraan
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

	newData, err := h.PerjalananDinasKendaraanService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// Update adalah untuk mengubah data perjalanan dinas kendaraan.
// @Summary mengubah data perjalanan dinas kendaraan
// @Description Endpoint ini adalah untuk mengubah data perjalanan dinas kendaraan.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Param Data body bpd.RequestPerjalananDinasKendaraan true "Perjalanan dinas kendaraan yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas/{id} [put]
func (h *PerjalananDinasKendaraanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))

	var reqFormat bpd.RequestPerjalananDinasKendaraan
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

	data, err := h.PerjalananDinasKendaraanService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID adalah untuk mendapatkan satu data perjalanan dinas kendaraan berdasarkan ID.
// @Summary Mendapatkan satu data perjalanan dinas kendaraan berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan perjalanan dinas kendaraan By ID.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas/{id} [get]
func (h *PerjalananDinasKendaraanHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	data, err := h.PerjalananDinasKendaraanService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data perjalanan dinas kendaraan.
// @Summary menghapus data perjalanan dinas kendaraan.
// @Description Endpoint ini adalah untuk menghapus data perjalanan dinas kendaraan.
// @Tags perjalanan-dinas-kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/kendaraan-perjalanan-dinas/{id} [delete]
func (h *PerjalananDinasKendaraanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.PerjalananDinasKendaraanService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
