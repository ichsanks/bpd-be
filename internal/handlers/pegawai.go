package handlers

import (
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

type PegawaiHandler struct {
	PegawaiService master.PegawaiService
	Config         *configs.Config
}

func ProvidePegawaiHandler(service master.PegawaiService, config *configs.Config) PegawaiHandler {
	return PegawaiHandler{
		PegawaiService: service,
		Config:         config,
	}
}

func (h *PegawaiHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/pegawai", func(r chi.Router) {
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

// ResolveAll list data Pegawai.
// @Summary Get list data Pegawai.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Pegawai sesuai dengan filter yang dikirimkan.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idBidang query string false "Set Id Bidang"
// @Param idUnor query string false "Set Id Unor"
// @Param idFungsionalitas query string false "Set Id Fungsionalitas"
// @Success 200 {object} master.Pegawai
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai [get]
func (h *PegawaiHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	idBidang := r.URL.Query().Get("idBidang")
	idUnor := r.URL.Query().Get("idUnor")
	idFungsionalitas := r.URL.Query().Get("idFungsionalitas")
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
		Keyword:          keyword,
		PageSize:         pageSize,
		PageNumber:       pageNumber,
		SortBy:           sortBy,
		SortType:         sortType,
		IdBidang:         idBidang,
		IdUnor:           idUnor,
		IdFungsionalitas: idFungsionalitas,
		IdBranch:         idBranchs,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.PegawaiService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all Pegawai.
// @Summary Get list all Pegawai.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Pegawai sesuai dengan filter yang dikirimkan.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai/all [get]
func (h *PegawaiHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	req := model.StandardRequest{
		IdBranch: idBranchs,
	}

	status, err := h.PegawaiService.GetAll(req)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// createPegawai adalah untuk menambah data Pegawai.
// @Summary menambahkan data Pegawai.
// @Description Endpoint ini adalah untuk menambahkan data Pegawai.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param nip formData string true "nip"
// @Param nama formData string true "nama"
// @Param kodeJk formData string false "kodeJk"
// @Param kodeAgama formData string true "kodeAgama"
// @Param alamat formData string false "alamat"
// @Param noHp formData string false "noHp"
// @Param email formData string false "email"
// @Param idUnor formData string false "idUnor"
// @Param idJabatan formData string false "idJabatan"
// @Param idGolongan formData string false "idGolongan"
// @Param idFungsionalitas formData string false "idFungsionalitas"
// @Param nik formData string false "nik"
// @Param idBranch formData string false "idBranch"
// @Param idStatusPegawai formData string false "idStatusPegawai"
// @Param idJobGrade formData string false "idJobGrade"
// @Param idPersonGrade formData string false "idPersonGrade"
// @Param idLevelBod formData string false "idLevelBod"
// @Param idStatusKontrak formData string false "idStatusKontrak"
// @Param file formData file false "File ttd"
// @Success 200 {object} response.Base{data=master.Pegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai [post]
func (h *PegawaiHandler) Create(w http.ResponseWriter, r *http.Request) {
	nip := r.FormValue("nip")
	nama := r.FormValue("nama")
	kodeJk := r.FormValue("kodeJk")
	kodeAgama := r.FormValue("kodeAgama")
	alamat := r.FormValue("alamat")
	noHp := r.FormValue("noHp")
	email := r.FormValue("email")
	idUnor := r.FormValue("idUnor")
	idJabatan := r.FormValue("idJabatan")
	idGolongan := r.FormValue("idGolongan")
	idFungsionalitas := r.FormValue("idFungsionalitas")
	idBidang := r.FormValue("idBidang")
	nik := r.FormValue("nik")
	idBranch := r.FormValue("idBranch")
	idStatusPegawai := r.FormValue("idStatusPegawai")
	idJobGrade := r.FormValue("idJobGrade")
	idPersonGrade := r.FormValue("idPersonGrade")
	idLevelBod := r.FormValue("idLevelBod")
	idStatusKontrak := r.FormValue("idStatusKontrak")
	kodeVendor := r.FormValue("kodeVendor")
	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.PegawaiService.UploadFile(w, r, "file", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	reqFormat := master.PegawaiFormat{
		Nip:              nip,
		Nama:             nama,
		IdJabatan:        &idJabatan,
		IdFungsionalitas: &idFungsionalitas,
		IdBidang:         idBidang,
	}

	if nik != "" {
		reqFormat.Nik = &nik
	}

	if path != "" {
		reqFormat.FotoTtd = &path
	}

	if kodeJk != "" {
		reqFormat.KodeJk = &kodeJk
	}

	if kodeAgama != "" {
		reqFormat.KodeAgama = &kodeAgama
	}

	if alamat != "" {
		reqFormat.Alamat = &alamat
	}

	if noHp != "" {
		reqFormat.NoHp = &noHp
	}

	if email != "" {
		reqFormat.Email = &email
	}

	if idUnor != "" {
		reqFormat.IdUnor = &idUnor
	}

	if idGolongan != "" {
		reqFormat.IdGolongan = &idGolongan
	}

	if idStatusPegawai != "" {
		reqFormat.IdStatusPegawai = &idStatusPegawai
	}

	if idJobGrade != "" {
		reqFormat.IdJobGrade = &idJobGrade
	}

	if idPersonGrade != "" {
		reqFormat.IdPersonGrade = &idPersonGrade
	}

	if idLevelBod != "" {
		reqFormat.IdLevelBod = &idLevelBod
	}

	if idStatusKontrak != "" {
		reqFormat.IdStatusKontrak = &idStatusKontrak
	}

	if kodeVendor != "" {
		reqFormat.KodeVendor = &kodeVendor
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if idBranch == "" {
		reqFormat.IdBranch = &idBranchs
	}

	newData, err := h.PegawaiService.Create(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// UpdatePegawai adalah untuk merubah data Pegawai.
// @Summary merubah data Pegawai
// @Description Endpoint ini adalah untuk merubah data Pegawai.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "id"
// @Param nip formData string true "nip"
// @Param nama formData string true "nama"
// @Param kodeJk formData string false "kodeJk"
// @Param kodeAgama formData string true "kodeAgama"
// @Param alamat formData string false "alamat"
// @Param noHp formData string false "noHp"
// @Param email formData string false "email"
// @Param idUnor formData string false "idUnor"
// @Param idJabatan formData string false "idJabatan"
// @Param idGolongan formData string false "idGolongan"
// @Param idFungsionalitas formData string false "idFungsionalitas"
// @Param nik formData string false "nik"
// @Param idBranch formData string false "idBranch"
// @Param idStatusPegawai formData string false "idStatusPegawai"
// @Param idJobGrade formData string false "idJobGrade"
// @Param idPersonGrade formData string false "idPersonGrade"
// @Param idLevelBod formData string false "idLevelBod"
// @Param idStatusKontrak formData string false "idStatusKontrak"
// @Param file formData file false "File ttd"
// @Success 200 {object} response.Base{data=master.Pegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai [put]
func (h *PegawaiHandler) Update(w http.ResponseWriter, r *http.Request) {

	idPegawai := r.FormValue("id")
	id, err := uuid.FromString(idPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}
	pg, err := h.PegawaiService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	nip := r.FormValue("nip")
	nama := r.FormValue("nama")
	kodeJk := r.FormValue("kodeJk")
	kodeAgama := r.FormValue("kodeAgama")
	alamat := r.FormValue("alamat")
	noHp := r.FormValue("noHp")
	email := r.FormValue("email")
	idUnor := r.FormValue("idUnor")
	idJabatan := r.FormValue("idJabatan")
	idGolongan := r.FormValue("idGolongan")
	idFungsionalitas := r.FormValue("idFungsionalitas")
	idBidang := r.FormValue("idBidang")
	nik := r.FormValue("nik")
	idBranch := r.FormValue("idBranch")
	idStatusPegawai := r.FormValue("idStatusPegawai")
	idJobGrade := r.FormValue("idJobGrade")
	idPersonGrade := r.FormValue("idPersonGrade")
	idLevelBod := r.FormValue("idLevelBod")
	idStatusKontrak := r.FormValue("idStatusKontrak")
	kodeVendor := r.FormValue("kodeVendor")
	uploadedFile, _, _ := r.FormFile("file")
	var path string
	ttd := pg.FotoTtd
	if uploadedFile != nil {
		filepath, err := h.PegawaiService.UploadFile(w, r, "file", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
		// remove file
		if ttd != nil {
			err = h.PegawaiService.DeleteFile(*ttd)
			if err != nil {
				println("delete file", err)
			}
		}
	} else {
		if ttd != nil {
			path = *ttd
		}
	}

	reqFormat := master.PegawaiFormat{
		ID:               id,
		Nip:              nip,
		Nama:             nama,
		IdJabatan:        &idJabatan,
		IdFungsionalitas: &idFungsionalitas,
		IdBidang:         idBidang,
	}

	if nik != "" {
		reqFormat.Nik = &nik
	}

	if path != "" {
		reqFormat.FotoTtd = &path
	}

	if kodeJk != "" {
		reqFormat.KodeJk = &kodeJk
	}

	if kodeAgama != "" {
		reqFormat.KodeAgama = &kodeAgama
	}

	if alamat != "" {
		reqFormat.Alamat = &alamat
	}

	if noHp != "" {
		reqFormat.NoHp = &noHp
	}

	if email != "" {
		reqFormat.Email = &email
	}

	if idUnor != "" {
		reqFormat.IdUnor = &idUnor
	}

	if idGolongan != "" {
		reqFormat.IdGolongan = &idGolongan
	}
	if idStatusPegawai != "" {
		reqFormat.IdStatusPegawai = &idStatusPegawai
	}

	if idJobGrade != "" {
		reqFormat.IdJobGrade = &idJobGrade
	}

	if idPersonGrade != "" {
		reqFormat.IdPersonGrade = &idPersonGrade
	}

	if kodeVendor != "" {
		reqFormat.KodeVendor = &kodeVendor
	}

	if idLevelBod != "" {
		reqFormat.IdLevelBod = &idLevelBod
	}

	if idStatusKontrak != "" {
		reqFormat.IdStatusKontrak = &idStatusKontrak
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if idBranch == "" {
		reqFormat.IdBranch = &idBranchs
	}
	newPegawai, err := h.PegawaiService.Update(reqFormat, userID, tenantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newPegawai)
}

// delete adalah untuk menghapus data Pegawai.
// @Summary menghapus data Pegawai.
// @Description Endpoint ini adalah untuk menghapus data Pegawai.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai/{id} [delete]
func (h *PegawaiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idPegawai, _ := uuid.FromString(id)
	err := h.PegawaiService.DeleteByID(idPegawai)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByID adalah untuk mendapatkan satu data Pegawai berdasarkan ID.
// @Summary Mendapatkan satu data Pegawai berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Pegawai By ID.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.Pegawai}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai/{id} [get]
func (h *PegawaiHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	unorLokasi, err := h.PegawaiService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, unorLokasi)
}

// delete adalah untuk menghapus data Pegawai.
// @Summary menghapus data Pegawai.
// @Description Endpoint ini adalah untuk menghapus data Pegawai.
// @Tags pegawai
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/pegawai/soft/{id} [delete]
func (h *PegawaiHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idPegawai, _ := uuid.FromString(id)
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}
	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}
	err = h.PegawaiService.DeleteSoft(idPegawai, userID, idBranchs)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
