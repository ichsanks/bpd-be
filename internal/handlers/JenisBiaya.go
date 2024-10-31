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

type JenisBiayaHandler struct {
	JenisBiayaService master.JenisBiayaService
	Config            *configs.Config
}

func ProvideJenisBiayaHandler(service master.JenisBiayaService, config *configs.Config) JenisBiayaHandler {
	return JenisBiayaHandler{
		JenisBiayaService: service,
		Config:            config,
	}
}

func (h *JenisBiayaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/jenis-biaya", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Put("/", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Get("/{id}", h.ResolveByID)
			r.Delete("/soft/{id}", h.DeleteSoft)
			r.Get("/jumlah/{id}/{ket}", h.GetJumlahBiayaByIdBod)
			r.Get("/all/dto", h.GetAllDataDto)
			r.Get("/all-header", h.GetAllHeader)
		})
	})
}

// ResolveAll list data JenisBiaya.
// @Summary Get list data JenisBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisBiaya sesuai dengan filter yang dikirimkan.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idKategoriBiaya query string false "id Kategori Biaya"
// @Success 200 {object} master.JenisBiaya
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya [get]
func (h *JenisBiayaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	idKategoriBiaya := r.URL.Query().Get("idKategoriBiaya")
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
		Keyword:     keyword,
		PageSize:    pageSize,
		PageNumber:  pageNumber,
		SortBy:      sortBy,
		SortType:    sortType,
		IdBranch:    idBranchs,
		IdTransaksi: idKategoriBiaya,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.JenisBiayaService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all JenisBiaya.
// @Summary Get list all JenisBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisBiaya sesuai dengan filter yang dikirimkan.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param kelompokBiaya query string false "Kelompok Biaya"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/all [get]
func (h *JenisBiayaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	kelompokBiaya := r.URL.Query().Get("kelompokBiaya")
	req := model.StandardRequest{
		IdBranch:    idBranchs,
		IdTransaksi: kelompokBiaya,
	}

	status, err := h.JenisBiayaService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createJenisBiaya adalah untuk menambah data JenisBiaya.
// @Summary menambahkan data JenisBiaya.
// @Description Endpoint ini adalah untuk menambahkan data JenisBiaya.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisBiaya body master.JenisBiayaFormat true "JenisBiaya yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.JenisBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya [post]
func (h *JenisBiayaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisBiayaFormat
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

	newData, err := h.JenisBiayaService.Create(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdateJenisBiaya adalah untuk merubah data JenisBiaya.
// @Summary merubah data JenisBiaya
// @Description Endpoint ini adalah untuk merubah data JenisBiaya.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param JenisBiaya body master.JenisBiayaFormat true "JenisBiaya yang akan dirubah"
// @Success 200 {object} response.Base{data=master.JenisBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya [put]
func (h *JenisBiayaHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.JenisBiayaFormat
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
	newJenisBiaya, err := h.JenisBiayaService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newJenisBiaya)
}

// delete adalah untuk menghapus data JenisBiaya.
// @Summary menghapus data JenisBiaya.
// @Description Endpoint ini adalah untuk menghapus data JenisBiaya.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/{id} [delete]
func (h *JenisBiayaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisBiaya, _ := uuid.FromString(id)
	err := h.JenisBiayaService.DeleteByID(idJenisBiaya)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data JenisBiaya berdasarkan ID.
// @Summary Mendapatkan satu data JenisBiaya berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan JenisBiaya By ID.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.JenisBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/{id} [get]
func (h *JenisBiayaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.JenisBiayaService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data JenisBiaya.
// @Summary menghapus data JenisBiaya.
// @Description Endpoint ini adalah untuk menghapus data JenisBiaya.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/soft/{id} [delete]
func (h *JenisBiayaHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idJenisBiaya, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.JenisBiayaService.DeleteSoft(idJenisBiaya, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data JenisBiaya berdasarkan ID.
// @Summary Mendapatkan satu data JenisBiaya berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan JenisBiaya By ID.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID BOD LEVEL"
// @Success 200 {object} response.Base{data=master.JumlahBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/jumlah/{id}/{ket} [get]
func (h *JenisBiayaHandler) GetJumlahBiayaByIdBod(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	ket := chi.URLParam(r, "ket")
	unorLokasi, err := h.JenisBiayaService.GetJumlahBiayaByIdBod(ID, ket)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// GetDataAll list all JenisBiaya.
// @Summary Get list all JenisBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisBiaya sesuai dengan filter yang dikirimkan.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBranch query string false "IdBranch"
// @Param idBodLevel query string false "IdBodLevel"
// @Param idJenisTujuan query string false "id jenis tujuan"
// @Param kelompokBiaya query string false "Kelompok Biaya"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/all/dto [get]
func (h *JenisBiayaHandler) GetAllDataDto(w http.ResponseWriter, r *http.Request) {
	idBranch := r.URL.Query().Get("idBranch")
	idBodLevel := r.URL.Query().Get("idBodLevel")
	idJenisTujuan := r.URL.Query().Get("idJenisTujuan")
	kelompokBiaya := r.URL.Query().Get("kelompokBiaya")

	req := model.StandardRequest{
		IdBranch:      idBranch,
		IdBodLevel:    idBodLevel,
		IdJenisTujuan: idJenisTujuan,
		IdTransaksi:   kelompokBiaya,
	}

	status, err := h.JenisBiayaService.GetAllDto(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all JenisBiaya.
// @Summary Get list all JenisBiaya.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data JenisBiaya sesuai dengan filter yang dikirimkan.
// @Tags jenisBiaya
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/jenis-biaya/all-header [get]
func (h *JenisBiayaHandler) GetAllHeader(w http.ResponseWriter, r *http.Request) {
	status, err := h.JenisBiayaService.GetAllHeader()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
