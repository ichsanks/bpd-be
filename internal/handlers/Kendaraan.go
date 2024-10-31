package handlers

import (
	"encoding/json"
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

type KendaraanHandler struct {
	KendaraanService master.KendaraanService
	Config           *configs.Config
}

func ProvideKendaraanHandler(service master.KendaraanService, config *configs.Config) KendaraanHandler {
	return KendaraanHandler{
		KendaraanService: service,
		Config:           config,
	}
}

func (h *KendaraanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/kendaraan", func(r chi.Router) {
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

// ResolveAll list data Kendaraan.
// @Summary Get list data Kendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Kendaraan sesuai dengan filter yang dikirimkan.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nopol | nama | namaJenisKendaraan ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idJenisKendaraan query string false "id jenis kendaraan"
// @Success 200 {object} master.Kendaraan
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan [get]
func (h *KendaraanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	idJenisKendaraan := r.URL.Query().Get("idJenisKendaraan")

	req := model.StandardRequestKendaraan{
		Keyword:          keyword,
		PageSize:         pageSize,
		PageNumber:       pageNumber,
		SortBy:           sortBy,
		SortType:         sortType,
		IdJenisKendaraan: idJenisKendaraan,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.KendaraanService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Kendaraan.
// @Summary Get list all Kendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Kendaraan sesuai dengan filter yang dikirimkan.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan/all [get]
func (h *KendaraanHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := h.KendaraanService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// createKendaraan adalah untuk menambah data Kendaraan.
// @Summary menambahkan data Kendaraan.
// @Description Endpoint ini adalah untuk menambahkan data Kendaraan.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Kendaraan body master.RequestKendaraan true "Kendaraan yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.Kendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan [post]
func (h *KendaraanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestKendaraan
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

	newData, err := h.KendaraanService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateKendaraan adalah untuk mengubah data Kendaraan.
// @Summary mengubah data Kendaraan
// @Description Endpoint ini adalah untuk mengubah data Kendaraan.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param Kendaraan body master.RequestKendaraan true "Kendaraan yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.Kendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan/{id} [put]
func (h *KendaraanHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RequestKendaraan
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

	Kendaraan, err := h.KendaraanService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, Kendaraan)
}

// ResolveByID adalah untuk mendapatkan satu data Kendaraan berdasarkan ID.
// @Summary Mendapatkan satu data Kendaraan berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Kendaraan By ID.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.Kendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan/{id} [get]
func (h *KendaraanHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.KendaraanService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data Kendaraan.
// @Summary menghapus data Kendaraan.
// @Description Endpoint ini adalah untuk menghapus data Kendaraan.
// @Tags kendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/kendaraan/{id} [delete]
func (h *KendaraanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.KendaraanService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
