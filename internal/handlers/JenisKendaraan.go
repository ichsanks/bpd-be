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

type JenisKendaraanHandler struct {
	JenisKendaraanService master.JenisKendaraanService
	Config                *configs.Config
}

func ProvideJenisKendaraanHandler(service master.JenisKendaraanService, config *configs.Config) JenisKendaraanHandler {
	return JenisKendaraanHandler{
		JenisKendaraanService: service,
		Config:                config,
	}
}

func (h *JenisKendaraanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-kendaraan", func(r chi.Router) {
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

// ResolveAll list data JenisKendaraan.
// @Summary Get list data JenisKendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisKendaraan sesuai dengan filter yang dikirimkan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} master.JenisKendaraan
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan [get]
func (h *JenisKendaraanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.JenisKendaraanService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all JenisKendaraan.
// @Summary Get list all JenisKendaraan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisKendaraan sesuai dengan filter yang dikirimkan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan/all [get]
func (h *JenisKendaraanHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	status, err := h.JenisKendaraanService.GetAll()

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createJenisKendaraan adalah untuk menambah data JenisKendaraan.
// @Summary menambahkan data JenisKendaraan.
// @Description Endpoint ini adalah untuk menambahkan data JenisKendaraan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisKendaraan body master.JenisKendaraanFormat true "JenisKendaraan yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.JenisKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan [post]
func (h *JenisKendaraanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisKendaraanFormat
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

	newData, err := h.JenisKendaraanService.Create(reqFormat, userID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateJenisKendaraan adalah untuk merubah data JenisKendaraan.
// @Summary merubah data JenisKendaraan
// @Description Endpoint ini adalah untuk merubah data JenisKendaraan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisKendaraan body master.JenisKendaraanFormat true "JenisKendaraan yang akan dirubah"
// @Success 200 {object} response.Base{data=master.JenisKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan [put]
func (h *JenisKendaraanHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisKendaraanFormat
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
	newJenisKendaraan, err := h.JenisKendaraanService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newJenisKendaraan)
}

// delete adalah untuk menghapus data JenisKendaraan.
// @Summary menghapus data JenisKendaraan.
// @Description Endpoint ini adalah untuk menghapus data JenisKendaraan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan/{id} [delete]
func (h *JenisKendaraanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisKendaraan, _ := uuid.FromString(id)
	err := h.JenisKendaraanService.DeleteByID(idJenisKendaraan)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data JenisKendaraan berdasarkan ID.
// @Summary Mendapatkan satu data JenisKendaraan berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan JenisKendaraan By ID.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.JenisKendaraan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan/{id} [get]
func (h *JenisKendaraanHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.JenisKendaraanService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data JenisKendaraan.
// @Summary menghapus data JenisKendaraan.
// @Description Endpoint ini adalah untuk menghapus data JenisKendaraan.
// @Tags jenisKendaraan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-kendaraan/soft/{id} [delete]
func (h *JenisKendaraanHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisKendaraan, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.JenisKendaraanService.DeleteSoft(idJenisKendaraan, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
