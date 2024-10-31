package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type SuratPerjalananDinasHandler struct {
	SuratPerjalananDinasService bpd.SuratPerjalananDinasService
	SppdDokumenService          bpd.SppdDokumenService
	PengajuanSppdHistoriService bpd.PengajuanSppdHistoriService
	PegawaiService              master.PegawaiService
	Config                      *configs.Config
}

func ProvideSuratPerjalananDinasHandler(service bpd.SuratPerjalananDinasService, service1 bpd.SppdDokumenService, pengajuan bpd.PengajuanSppdHistoriService, pegawai master.PegawaiService, config *configs.Config) SuratPerjalananDinasHandler {
	return SuratPerjalananDinasHandler{
		SuratPerjalananDinasService: service,
		SppdDokumenService:          service1,
		PengajuanSppdHistoriService: pengajuan,
		PegawaiService:              pegawai,
		Config:                      config,
	}
}

func (h *SuratPerjalananDinasHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/surat-perjalanan-dinas", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Put("/", h.Update)
			r.Get("/", h.ResolveAll)
			r.Get("/{id}", h.ResolveByID)
			r.Delete("/{id}", h.Delete)
			r.Get("/detail", h.ResolveByDetailID)
			r.Get("/list-approval", h.ResolveAllApproval)
			r.Post("/verifikasi-esign", h.VerifikasiEsign)
			r.Get("/biaya-pegawai", h.BiayaPegawai)
			r.Post("/update-file", h.UpdateFile)
			r.Post("/update-filependukung", h.UpdateDokumenPendukung)
		})
	})
}

// create adalah untuk menambah data Surat Perjalanan Dinas.
// @Summary menambahkan data Surat Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Surat Perjalanan Dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SuratPerjalananDinas body bpd.SuratPerjalananDinasRequest true "Surat Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.SuratPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas [post]
func (h *SuratPerjalananDinasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.SuratPerjalananDinasRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	fmt.Println("UserID", userID)
	data, err := h.SuratPerjalananDinasService.Create(reqFormat, userID)
	if err != nil {
		fmt.Print("error create")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	for _, d := range reqFormat.DetailDokumen {
		var detID uuid.UUID
		detID, _ = uuid.NewV4()
		sppdDokumen := bpd.SppdDokumenRequest{
			ID:         detID,
			IdSppd:     data.ID.String(),
			IdDokumen:  d.ID.String(),
			File:       d.File,
			Keterangan: d.Keterangan,
		}
		_, err = h.SppdDokumenService.Create(sppdDokumen, userID)
		if err != nil {
			fmt.Print("error create")
			response.WithError(w, failure.BadRequest(err))
			return
		}

	}

	response.WithJSON(w, http.StatusCreated, data)
}

// update adalah untuk menambah data Surat Perjalanan Dinas.
// @Summary menambahkan data Surat Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Surat Perjalanan Dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SuratPerjalananDinas body bpd.SuratPerjalananDinasRequest true "Surat Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.SuratPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas [put]
func (h *SuratPerjalananDinasHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.SuratPerjalananDinasRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.SuratPerjalananDinasService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	for _, d := range reqFormat.DetailDokumen {
		sppdDokumen := bpd.SppdDokumenRequest{
			ID:         d.ID,
			IdSppd:     d.IdSppd,
			IdDokumen:  d.IdDokumen,
			File:       d.File,
			Keterangan: d.Keterangan,
		}
		_, err = h.SppdDokumenService.Update(sppdDokumen, userID)
		if err != nil {
			fmt.Print("error create")
			response.WithError(w, failure.BadRequest(err))
			return
		}

	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveAll list surat perjalanan dinas.
// @Summary Get list data surat perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data surat perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of "
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idBranch query string false "Set id branch"
// @Param idPegawai query string false "Set id Pegawai"
// @Param status query string false "set status"
// @Success 200 {object} bpd.SuratPerjalananDinasDto
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas [get]
func (h *SuratPerjalananDinasHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	idPegawai := r.URL.Query().Get("idPegawai")
	statuss := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "desc"
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

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	idBranchStr := r.URL.Query().Get("idBranch")
	if idBranchStr != "" {
		idBranch = idBranchStr
	}

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		IdBranch:   idBranch,
		IdPegawai:  idPegawai,
		Status:     statuss,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.SuratPerjalananDinasService.ResolveAllDto(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// ResolveByID adalah untuk mendapatkan satu data surat perjalanan dinas berdasarkan ID.
// @Summary Mendapatkan satu data surat perjalanan dinas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan surat perjalanan dinas By ID.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=bpd.SuratPerjalananDinasDto}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/{id} [get]
func (h *SuratPerjalananDinasHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}

	data, err := h.SuratPerjalananDinasService.ResolveDtoByID(ID.String())
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data Surat Perjalanan Dinas.
// @Summary hapus data Surat Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menghapus data Surat Perjalanan Dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/{id} [delete]
func (h *SuratPerjalananDinasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	err := h.SuratPerjalananDinasService.DeleteByID(id, userID)
	if err != nil {
		fmt.Println(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveByDetailID adalah untuk mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Summary Mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Perjalanan Dinas By ID.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id query string true "Set ID"
// @Param idPgw query string false "Set ID Pegawai"
// @Param idPegawai query string false "Set ID Pegawai Approval"
// @Param typeApproval query string false "Set Type Approval"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/detail [get]
func (h *SuratPerjalananDinasHandler) ResolveByDetailID(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	idPegawai := r.URL.Query().Get("idPgw")
	idPegawaiApproval := r.URL.Query().Get("idPegawai")
	typeApproval := r.URL.Query().Get("typeApproval")

	req := bpd.FilterDetailSPPD{
		ID:                ID,
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		TypeApproval:      typeApproval,
	}

	data, err := h.SuratPerjalananDinasService.ResolveByDetailID(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveAllApproval list all Surat Perjalanan Dinas Approval.
// @Summary Get list all Surat Perjalanan Dinas Approval.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data surat Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idFungsionalitas query string false "Set ID Fungsionalitas"
// @Param idUnor query string false "Set ID Unor"
// @Param typeApproval query string false "Set TypeApproval"
// @Param idPegawaiApproval query string false "Set ID Pegawai Approval"
// @Param status query string false "Set status BPD"
// @Param tglBerangkat query string false "Set tgl Berangkat"
// @Param tglKembali query string false "Set tgl Kembali"
// @Param jenisSppd query string false "Set jenisSppd"
// @Param idBidang query string false "Set idBidang"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/list-approval [get]
func (h *SuratPerjalananDinasHandler) ResolveAllApproval(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	idFungsionalitas := r.URL.Query().Get("idFungsionalitas")
	idUnor := r.URL.Query().Get("idUnor")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	typeApproval := r.URL.Query().Get("typeApproval")
	status := r.URL.Query().Get("status")
	tglBerangkat := r.URL.Query().Get("tglBerangkat")
	tglKembali := r.URL.Query().Get("tglKembali")
	sortBy := r.URL.Query().Get("sortBy")
	jenisSppd := r.URL.Query().Get("jenisSppd")
	idBidang := r.URL.Query().Get("idBidang")
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

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	idBranchStr := r.URL.Query().Get("idBranch")
	if idBranchStr != "" {
		idBranch = idBranchStr
	}

	req := model.StandardRequest{
		Keyword:           keyword,
		PageSize:          pageSize,
		PageNumber:        pageNumber,
		SortBy:            sortBy,
		SortType:          sortType,
		IdFungsionalitas:  idFungsionalitas,
		IdUnor:            idUnor,
		TypeApproval:      typeApproval,
		IdPegawaiApproval: idPegawaiApproval,
		Status:            status,
		StartDate:         tglBerangkat,
		EndDate:           tglKembali,
		IdBranch:          idBranch,
		IdTransaksi:       jenisSppd,
		IdBidang:          idBidang,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.SuratPerjalananDinasService.ResolveAllApproval(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk upload file surat perjalanan dinas.
// @Summary upload file surat perjalanan dinas.
// @Description Endpoint ini adalah untuk upload file surat perjalanan dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "ID Perjalanan Dinas"
// @Param idPengajuan formData string true "ID Pengajuan Dinas"
// @Param idPegawaiApproval formData string false "ID Pegawai Approval"
// @Param passphrase formData string false "Pashparase"
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/verifikasi-esign [post]
func (h *SuratPerjalananDinasHandler) VerifikasiEsign(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idPengajuan := r.FormValue("idPengajuan")
	idPegawaiApproval := r.FormValue("idPegawaiApproval")
	idPgw, _ := uuid.FromString(idPegawaiApproval)
	passphrase := r.FormValue("passphrase")
	uploadedFile, _, _ := r.FormFile("file")

	// Get By ID
	pd, err := h.SuratPerjalananDinasService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	filename := ""
	if pd.File != nil {
		filename = *pd.File
	}

	var path string
	if uploadedFile != nil {
		filepath, err := h.SuratPerjalananDinasService.UploadFile(w, r, "file", filename, id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.FilesSuratPerjalananDinas{
		ID:   id,
		File: path,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.SuratPerjalananDinasService.UploadFileSuratPerjalananDinas(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var reqFormat2 = bpd.LinkFilesSuratPerjalananDinas{
		ID:       id,
		LinkFile: path,
	}

	_, err = h.SuratPerjalananDinasService.UploadLinkFileSuratPerjalananDinas(reqFormat2, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	pg, err := h.PegawaiService.ResolveByID(idPgw)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var reqEsingFormat = bpd.EsignRequest{
		Nik:        *pg.Nik,
		Passphrase: passphrase,
		Pdf:        newData.File,
		Ttd:        pg.FotoTtd,
	}

	if passphrase != "" {
		err = h.SuratPerjalananDinasService.BsreResolveSignPdf(reqEsingFormat)
		if err != nil {
			response.WithError(w, err)
			err = h.SuratPerjalananDinasService.UpdateFile(id, userID)
			return
		}
	}
	payload := bpd.PengajuanSppdHistoriInputRequest{
		ID:           idPengajuan,
		IdPegawai:    &idPegawaiApproval,
		Status:       "2",
		TypeApproval: "PENGAJUAN_SPPD",
	}

	_, err = h.PengajuanSppdHistoriService.Approve(payload, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// ResolveAll list surat perjalanan dinas.
// @Summary Get list data surat perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data surat perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param tglAwal query string false "Tanggal Awal Dinas"
// @Param tglAkhir query string true "Tanggal Akhir Dinas"
// @Param idSppd query string true "Id SPPD"
// @Param idBodLevel query string true "Id BOD Level"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/biaya-pegawai [get]
func (h *SuratPerjalananDinasHandler) BiayaPegawai(w http.ResponseWriter, r *http.Request) {
	tglAwal := r.URL.Query().Get("tglAwal")
	tglAkhir := r.URL.Query().Get("tglAkhir")
	idSppd := r.URL.Query().Get("idSppd")
	idBodLevel := r.URL.Query().Get("idBodLevel")

	status, err := h.SuratPerjalananDinasService.GetDetailBiaya(tglAwal, tglAkhir, idSppd, idBodLevel)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// create adalah untuk upload file surat perjalanan dinas.
// @Summary upload file surat perjalanan dinas.
// @Description Endpoint ini adalah untuk upload file surat perjalanan dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "ID Perjalanan Dinas"
// @Param path formData string true "path"
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/update-file [post]
func (h *SuratPerjalananDinasHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	// pathFile := r.FormValue("path")
	uploadedFile, _, _ := r.FormFile("file")

	// Get By ID
	pd, err := h.SuratPerjalananDinasService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	filename := ""
	if pd.File != nil {
		filename = *pd.File
	}

	// dir := h.Config.App.File.Dir
	// DokumenDir := pathFile
	// fileLocation := filepath.Join(dir, DokumenDir)
	// if pathFile != "" {
	// 	err = os.Remove(fileLocation)
	// }

	if err != nil {
		response.WithError(w, err)
		return
	}

	var path string
	if uploadedFile != nil {
		filepath, err := h.SuratPerjalananDinasService.UploadFile(w, r, "file", filename, id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.FilesSuratPerjalananDinas{
		ID:   id,
		File: path,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.SuratPerjalananDinasService.UploadFileSuratPerjalananDinas(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var reqFormat2 = bpd.LinkFilesSuratPerjalananDinas{
		ID:       id,
		LinkFile: path,
	}

	_, err = h.SuratPerjalananDinasService.UploadLinkFileSuratPerjalananDinas(reqFormat2, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// updateDokumenPendukung adalah untuk menambah data Surat Perjalanan Dinas.
// @Summary menambahkan data Surat Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Surat Perjalanan Dinas.
// @Tags surat-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param SuratPerjalananDinas body bpd.UpdateDokumenPendukung true "Surat Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.SuratPerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/surat-perjalanan-dinas/update-filependukung [post]
func (h *SuratPerjalananDinasHandler) UpdateDokumenPendukung(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.UpdateDokumenPendukung
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	fmt.Println("UserID", userID)

	for _, d := range reqFormat.DetailDokumen {
		sppdDokumen := bpd.SppdDokumenRequest{
			ID:         d.ID,
			IdSppd:     d.IdSppd,
			IdDokumen:  d.IdDokumen,
			File:       d.File,
			Keterangan: d.Keterangan,
		}
		_, err = h.SppdDokumenService.Update(sppdDokumen, userID)
		if err != nil {
			fmt.Print("error create")
			response.WithError(w, failure.BadRequest(err))
			return
		}

	}

	response.WithJSON(w, http.StatusCreated, "success")
}
