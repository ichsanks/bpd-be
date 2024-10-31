package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PerjalananDinasHandler struct {
	PerjalananDinasService      bpd.PerjalananDinasService
	PerjalananDinasBiayaService bpd.PerjalananDinasBiayaService
	PengajuanBpdHistoriService  bpd.PengajuanBpdHistoriService
	PegawaiService              master.PegawaiService
}

func ProvidePerjalananDinasHandler(service bpd.PerjalananDinasService, pengajuan bpd.PengajuanBpdHistoriService, pegawai master.PegawaiService, biaya bpd.PerjalananDinasBiayaService) PerjalananDinasHandler {
	return PerjalananDinasHandler{
		PerjalananDinasService:      service,
		PerjalananDinasBiayaService: biaya,
		PengajuanBpdHistoriService:  pengajuan,
		PegawaiService:              pegawai,
	}
}

func (h *PerjalananDinasHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/perjalanan-dinas", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/list-approval", h.ResolveAllApproval)
			r.Get("/list-penyelesaian", h.ResolveAllPenyelesaian)
			r.Post("/", h.Create)
			r.Put("/", h.Update)
			r.Get("/{id}", h.ResolveByIDDTO)
			r.Get("/detail", h.ResolveByDetailID)
			r.Get("/bpd-pegawai", h.ResolveByBpdPegawaiID)
			r.Delete("/{id}", h.Delete)
			r.Post("/upload", h.UploadFilePerjalananDinas)
			r.Post("/esign", h.CreateEsign)
			r.Post("/verifikasi-esign", h.VerifikasiEsign)
			r.Post("/penyelesaian-esign", h.PenyelesaianEsign)
			r.Post("/upload/tiket", h.UploadFile)
			r.Get("/exist-by-id-sppd", h.ExistByIdSppd)
			r.Get("/exist-reimbursement", h.ExistReimbursement)
			r.Get("/biaya-pegawai", h.GetDetailBiaya)
			r.Delete("/delete-file", h.DeleteFile)
			r.Get("/detail-histori", h.ResolveByDetailHistori)
		})
	})

	r.Route("/bpd/bpd-pegawai", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Put("/", h.UpdateBpdPegawai)
		})
	})
}

// ResolveAll list all Perjalanan Dinas.
// @Summary Get list all Perjalanan Dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idPegawai query string false "Set ID Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas [get]
func (h *PerjalananDinasHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	idPegawai := r.URL.Query().Get("idPegawai")
	sortBy := r.URL.Query().Get("sortBy")
	status := r.URL.Query().Get("status")
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
		IdPegawai:  idPegawai,
		Status:     status,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.PerjalananDinasService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveAllApproval list all Perjalanan Dinas Approval.
// @Summary Get list all Perjalanan Dinas Approval.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
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
// @Param idBidang query string false "Set id bidang"
// @Param jenisSppd query string false "Set jenis Sppd"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/list-approval [get]
func (h *PerjalananDinasHandler) ResolveAllApproval(w http.ResponseWriter, r *http.Request) {
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
	idBidang := r.URL.Query().Get("idBidang")
	jenisSppd := r.URL.Query().Get("jenisSppd")
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
		IdFungsionalitas:  idFungsionalitas,
		IdUnor:            idUnor,
		TypeApproval:      typeApproval,
		IdPegawaiApproval: idPegawaiApproval,
		Status:            status,
		StartDate:         tglBerangkat,
		EndDate:           tglKembali,
		IdBidang:          idBidang,
		IdTransaksi:       jenisSppd,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.PerjalananDinasService.ResolveAllApproval(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveAllPenyelesaian list all Perjalanan Dinas Approval penyelesaian.
// @Summary Get list all Perjalanan Dinas Approval penyelesaian.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas penyelesaian sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
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
// @Param idPegawai query string false "Set ID Pegawai"
// @Param idPegawaiApproval query string false "Set ID Pegawai Approval"
// @Param status query string false "Set Status"
// @Param isSppb query string false "Set Status"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/list-penyelesaian [get]
func (h *PerjalananDinasHandler) ResolveAllPenyelesaian(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	idFungsionalitas := r.URL.Query().Get("idFungsionalitas")
	idUnor := r.URL.Query().Get("idUnor")
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	status := r.URL.Query().Get("status")
	sppb := r.URL.Query().Get("isSppb")
	typeApproval := r.URL.Query().Get("typeApproval")
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
		IdFungsionalitas:  idFungsionalitas,
		IdUnor:            idUnor,
		TypeApproval:      typeApproval,
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		Status:            status,
		IsSppb:            sppb,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.PerjalananDinasService.ResolveAllPenyelesaian(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data Perjalanan Dinas.
// @Summary menambahkan data Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Perjalanan Dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PerjalananDinasRequest true "Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas [post]
func (h *PerjalananDinasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PerjalananDinasRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	fmt.Println("UserID", userID)
	data, err := h.PerjalananDinasService.Create(reqFormat, userID)
	if err != nil {
		fmt.Print("error create")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// update adalah untuk menambah data Perjalanan Dinas.
// @Summary menambahkan data Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menambahkan data Perjalanan Dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.PerjalananDinasRequest true "Perjalanan Dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinas}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas [put]
func (h *PerjalananDinasHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.PerjalananDinasRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PerjalananDinasService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// updateBpdPegawai adalah untuk update data Perjalanan Dinas Pegawai.
// @Summary update data Perjalanan Dinas Pegawai.
// @Description Endpoint ini adalah untuk update data Perjalanan Dinas Pegawai.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Data body bpd.BpdPegawaiRequest true "Bpd Pegawai yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasPegawaiDetail}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/bpd-pegawai [put]
func (h *PerjalananDinasHandler) UpdateBpdPegawai(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.BpdPegawaiRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	data, err := h.PerjalananDinasService.UpdateBpdPegawai(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// ResolveByIDDTO adalah untuk mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Summary Mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Perjalanan Dinas By ID.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/{id} [get]
func (h *PerjalananDinasHandler) ResolveByIDDTO(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	data, err := h.PerjalananDinasService.ResolveByIDDTO(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data Perjalanan Dinas.
// @Summary hapus data Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menghapus data Perjalanan Dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/{id} [delete]
func (h *PerjalananDinasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	err := h.PerjalananDinasService.DeleteByID(id, userID)
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
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id query string true "Set ID"
// @Param idPgw query string false "Set ID Pegawai"
// @Param idPegawai query string false "Set ID Pegawai Approval"
// @Param typeApproval query string false "Set Type Approval"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/detail [get]
func (h *PerjalananDinasHandler) ResolveByDetailID(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	idPegawai := r.URL.Query().Get("idPgw")
	idPegawaiApproval := r.URL.Query().Get("idPegawai")
	typeApproval := r.URL.Query().Get("typeApproval")

	req := bpd.FilterDetailBPD{
		ID:                ID,
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		TypeApproval:      typeApproval,
	}

	data, err := h.PerjalananDinasService.ResolveByDetailID(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk upload file perjalanan dinas.
// @Summary upload file perjalanan dinas.
// @Description Endpoint ini adalah untuk upload file perjalanan dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "ID Perjalanan Dinas"
// @Param idPegawaiApproval formData string false "ID Pegawai Approval"
// @Param passphrase formData string false "Pashparase"
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/upload [post]
func (h *PerjalananDinasHandler) UploadFilePerjalananDinas(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idPegawaiApproval := r.FormValue("idPegawaiApproval")
	idPgw, _ := uuid.FromString(idPegawaiApproval)
	passphrase := r.FormValue("passphrase")
	uploadedFile, _, _ := r.FormFile("file")

	// Get By ID
	pd, err := h.PerjalananDinasService.ResolveByID(id)
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
		filepath, err := h.PerjalananDinasService.UploadFile(w, r, "file", filename, id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.FilesPerjalananDinas{
		ID:   id,
		File: path,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.PerjalananDinasService.UploadFilePerjalananDinas(reqFormat, userID)
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
	err = h.PerjalananDinasService.BsreResolveSignPdf(reqEsingFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// CreateEsign adalah untuk menambah data sign PDF.
// @Summary menambahkan data SIGN PDG.
// @Description Endpoint ini adalah untuk menambahkan data ESIGN.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param PerjalananDinas body bpd.EsignRequest true "esign yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.EsignRequest}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/esign [post]
func (h *PerjalananDinasHandler) CreateEsign(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.EsignRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	fmt.Println("UserID", userID)
	err = h.PerjalananDinasService.BsreResolveSignPdf(reqFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, err)
}

// create adalah untuk upload file perjalanan dinas.
// @Summary upload file perjalanan dinas.
// @Description Endpoint ini adalah untuk upload file perjalanan dinas.
// @Tags perjalanan-dinas
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
// @Router /v1/bpd/perjalanan-dinas/verifikasi-esign [post]
func (h *PerjalananDinasHandler) VerifikasiEsign(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idPengajuan := r.FormValue("idPengajuan")
	idPegawaiApproval := r.FormValue("idPegawaiApproval")
	idPgw, _ := uuid.FromString(idPegawaiApproval)
	passphrase := r.FormValue("passphrase")
	uploadedFile, _, _ := r.FormFile("file")

	// Get By ID
	pd, err := h.PerjalananDinasService.ResolveByID(id)
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
		filepath, err := h.PerjalananDinasService.UploadFile(w, r, "file", filename, id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.FilesPerjalananDinas{
		ID:                     id,
		File:                   path,
		IdJenisPerjalananDinas: "true",
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.PerjalananDinasService.UploadFilePerjalananDinas(reqFormat, userID)
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
		err = h.PerjalananDinasService.BsreResolveSignPdf(reqEsingFormat)
		if err != nil {
			response.WithError(w, err)
			err = h.PerjalananDinasService.UpdateFile(id, userID)
			return
		}
	}
	payload := bpd.PengajuanBpdHistoriInputRequest{
		ID:           idPengajuan,
		IdPegawai:    &idPegawaiApproval,
		Status:       "2",
		TypeApproval: "PENGAJUAN_BPD",
	}

	_, err = h.PengajuanBpdHistoriService.Approve(payload, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// PenyelesaianEsign adalah untuk verifikasi penyelesaian esign.
// @Summary verifikasi penyelesaian esign.
// @Description Endpoint ini adalah untuk verifikasi penyelesaian esign.
// @Tags perjalanan-dinas
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
// @Router /v1/bpd/perjalanan-dinas/penyelesaian-esign [post]
func (h *PerjalananDinasHandler) PenyelesaianEsign(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idPengajuan := r.FormValue("idPengajuan")
	idPegawaiApproval := r.FormValue("idPegawaiApproval")
	idPgw, _ := uuid.FromString(idPegawaiApproval)
	passphrase := r.FormValue("passphrase")
	uploadedFile, _, _ := r.FormFile("file")

	// Get By ID
	pd, err := h.PerjalananDinasService.ResolveBpdPegawaiByID(id)
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
		filepath, err := h.PerjalananDinasBiayaService.UploadFile(w, r, "file", filename)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.DocPenyelesaianBpdPegawai{
		ID:   id,
		File: path,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.PerjalananDinasBiayaService.UploadDocPenyelesaianBpd(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	pg, err := h.PegawaiService.ResolveByID(idPgw)
	if err != nil {
		response.WithError(w, err)
		return
	}

	// ttd := "/files/ttd.png"
	var reqEsingFormat = bpd.EsignRequest{
		Nik:        *pg.Nik,
		Passphrase: passphrase,
		Pdf:        newData.File,
		Ttd:        pg.FotoTtd,
		// Ttd: &ttd,
	}
	err = h.PerjalananDinasService.BsreResolveSignPdf(reqEsingFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}

	payload := bpd.PengajuanBpdHistoriInputRequest{
		ID:           idPengajuan,
		IdPegawai:    &idPegawaiApproval,
		Status:       "2",
		TypeApproval: "PENYELESAIAN",
	}

	_, err = h.PengajuanBpdHistoriService.Approve(payload, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, "success")
}

// ResolveByBpdPegawaiID adalah untuk mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Summary Mendapatkan satu data Perjalanan Dinas berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Perjalanan Dinas By ID.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id query string true "Set ID"
// @Param idPegawai query string false "Set ID Pegawai"
// @Param idPegawaiApproval query string false "Set ID Pegawai Approval"
// @Param typeApproval query string false "Set Type Approval"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/bpd-pegawai [get]
func (h *PerjalananDinasHandler) ResolveByBpdPegawaiID(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	idPegawai := r.URL.Query().Get("idPegawai")
	idPegawaiApproval := r.URL.Query().Get("idPegawaiApproval")
	typeApproval := r.URL.Query().Get("typeApproval")
	req := bpd.FilterDetailBPD{
		ID:                ID,
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		TypeApproval:      typeApproval,
	}

	data, err := h.PerjalananDinasService.ResolveBpdPegawaiByIDDTO(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk upload file perjalanan dinas.
// @Summary upload file perjalanan dinas.
// @Description Endpoint ini adalah untuk upload file perjalanan dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/upload/tiket [post]
func (h *PerjalananDinasHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	uploadedFile, _, _ := r.FormFile("file")

	var path string
	if uploadedFile != nil {
		filepath, err := h.PerjalananDinasService.UploadFile(w, r, "file", "", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	response.WithJSON(w, http.StatusCreated, path)
}

// ResolveAll list all Perjalanan Dinas.
// @Summary Get list all Perjalanan Dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param idSppd query string false "Id Sppd"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/exist-by-id-sppd [get]
func (h *PerjalananDinasHandler) ExistByIdSppd(w http.ResponseWriter, r *http.Request) {
	idSppd := r.URL.Query().Get("idSppd")
	data, err := h.PerjalananDinasService.ExistByIdSppd(idSppd)
	param := map[string]interface{}{
		"success": data,
		"id":      err,
	}
	response.WithJSON(w, http.StatusOK, param)
}

// ResolveAll list all Perjalanan Dinas.
// @Summary Get list all Perjalanan Dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param idBpdPegawai query string false "Id Bpd Pegawai"
// @Param idPegawai query string false "Id Bpd Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/exist-reimbursement [get]
func (h *PerjalananDinasHandler) ExistReimbursement(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	idBpd, _ := uuid.FromString(idBpdPegawai)
	idPegawai := r.URL.Query().Get("idPegawai")
	idPeg, _ := uuid.FromString(idPegawai)
	data, err := h.PerjalananDinasService.ExistReimbursement(idBpd, idPeg)
	param := map[string]interface{}{
		"success": data,
		"id":      err,
	}
	response.WithJSON(w, http.StatusOK, param)
}

// ResolveAll list all Perjalanan Dinas.
// @Summary Get list all Perjalanan Dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Perjalanan Dinas sesuai dengan filter yang dikirimkan.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param idBpd query string false "Id Bpd Pegawai"
// @Param idPegawai query string false "Id Bpd Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/biaya-pegawai [get]
func (h *PerjalananDinasHandler) GetDetailBiaya(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpd")
	idPegawai := r.URL.Query().Get("idPegawai")
	data, err := h.PerjalananDinasService.GetDetailBiaya(idBpdPegawai, idPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data Perjalanan Dinas.
// @Summary hapus data Perjalanan Dinas.
// @Description Endpoint ini adalah untuk menghapus data Perjalanan Dinas.
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param pathFile query string false "path file"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/delete-file [delete]
func (h *PerjalananDinasHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	pathFile := r.URL.Query().Get("pathFile")
	err := h.PerjalananDinasService.DeleteFile(pathFile)
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
// @Tags perjalanan-dinas
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id query string true "Set ID"
// @Param idPgw query string false "Set ID Pegawai"
// @Param idPegawai query string false "Set ID Pegawai Approval"
// @Param typeApproval query string false "Set Type Approval"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/perjalanan-dinas/detail-histori [get]
func (h *PerjalananDinasHandler) ResolveByDetailHistori(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")
	idPegawai := r.URL.Query().Get("idPgw")
	idPegawaiApproval := r.URL.Query().Get("idPegawai")
	typeApproval := r.URL.Query().Get("typeApproval")

	req := bpd.FilterDetailBPD{
		ID:                ID,
		IdPegawai:         idPegawai,
		IdPegawaiApproval: idPegawaiApproval,
		TypeApproval:      typeApproval,
	}

	data, err := h.PerjalananDinasService.ResolveByDetailHistori(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}
