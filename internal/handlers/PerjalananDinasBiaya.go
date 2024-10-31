package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PerjalananDinasBiayaHandler struct {
	PerjalananDinasBiayaService bpd.PerjalananDinasBiayaService
	PerjalananDinasService      bpd.PerjalananDinasService
}

func ProvidePerjalananDinasBiayaHandler(service bpd.PerjalananDinasBiayaService, pdService bpd.PerjalananDinasService) PerjalananDinasBiayaHandler {
	return PerjalananDinasBiayaHandler{
		PerjalananDinasBiayaService: service,
		PerjalananDinasService:      pdService,
	}
}

func (h *PerjalananDinasBiayaHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/biaya-perjalanan-dinas", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.CreateBulk)
			r.Get("/all", h.GetAllData)
			r.Post("/um", h.CreateBulkUm)
			r.Get("/all-um", h.GetAllData)
			r.Post("/create", h.Create)
			r.Put("/update", h.Update)
			r.Get("/dto", h.GetAllDto)
			r.Delete("/{id}", h.Delete)
			r.Get("/histori-biaya", h.GetHistoriBiaya)
		})
	})

	r.Route("/bpd/penyelesaian", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/upload", h.UploadDocPenyelesaianBpd)
		})
	})
}

// create adalah untuk menambah data biaya perjalanan dinas.
// @Summary menambahkan data biaya perjalanan dinas.
// @Description Endpoint ini adalah untuk menambahkan data biaya perjalanan dinas.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RuleApproval body bpd.RequestPerjalananDinasBiaya true "Biaya perjalanan dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas [post]
func (h *PerjalananDinasBiayaHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.RequestPerjalananDinasBiaya
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

	data, err := h.PerjalananDinasBiayaService.CreateBulk(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// create adalah untuk menambah data biaya uang muka perjalanan dinas.
// @Summary menambahkan data biaya uang muka perjalanan dinas.
// @Description Endpoint ini adalah untuk menambahkan data biaya uang muka perjalanan dinas.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RuleApproval body bpd.RequestPerjalananDinasBiaya true "Biaya perjalanan dinas yang akan ditambahkan"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasBiaya}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/um [post]
func (h *PerjalananDinasBiayaHandler) CreateBulkUm(w http.ResponseWriter, r *http.Request) {
	var reqFormat bpd.RequestPerjalananDinasBiaya
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

	data, err := h.PerjalananDinasBiayaService.CreateBulkUm(reqFormat, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// GetDataAll list all biaya perjalanan dinas.
// @Summary Get list all biaya perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data biaya perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai query string false "Set ID Perjalanan Dinas Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/all [get]
func (h *PerjalananDinasBiayaHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.PerjalananDinasBiayaService.GetAllData(idBpdPegawai)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// GetDataAll list all biaya uang muka perjalanan dinas.
// @Summary Get list all biaya uang muka perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data biaya uang muka perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai query string false "Set ID Perjalanan Dinas Pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/all-um [get]
func (h *PerjalananDinasBiayaHandler) GetAllDataUm(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.PerjalananDinasBiayaService.GetAllDataUm(idBpdPegawai)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data log kegiatan.
// @Summary menambahkan data log kegiatan.
// @Description Endpoint ini adalah untuk menambahkan data log kegiatan.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai formData string true "ID BPD Pegawai"
// @Param file formData file true "Foto"
// @Success 200 {object} response.Base{data=bpd.LogKegiatan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/penyelesaian/upload [post]
func (h *PerjalananDinasBiayaHandler) UploadDocPenyelesaianBpd(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.FormValue("idBpdPegawai")
	uploadedFile, _, _ := r.FormFile("file")

	// get bpd pegawai id
	pd, err := h.PerjalananDinasService.ResolveBpdPegawaiByID(idBpdPegawai)
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
		ID:   idBpdPegawai,
		File: path,
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.PerjalananDinasBiayaService.UploadDocPenyelesaianBpd(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// create adalah untuk menambah data biaya perjalanan dinas.
// @Summary menambahkan data biaya perjalanan dinas.
// @Description Endpoint ini adalah untuk menambahkan data biaya perjalanan dinas.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "ID"
// @Param idBpdPegawai formData string true "ID BPD Pegawai"
// @Param idJenisBiaya formData string true "ID Jenis Biaya"
// @Param idKomponenBiaya formData string false "Id komponen Biaya"
// @Param nominal formData string false "Nominal"
// @Param idPegawai formData string false "Id Pegawai"
// @Param isReimbursement formData string false "Is reimbursement"
// @Param keterangan formData string false "keterangan"
// @Param tanggal formData string true "tanggal"
// @Param jenis formData string false "jenis"
// @Param file formData file true "Foto"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasBiayaDetail}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/create [post]
func (h *PerjalananDinasBiayaHandler) Create(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idBpdPegawai := r.FormValue("idBpdPegawai")
	idJenisBiaya := r.FormValue("idJenisBiaya")
	idKomponenBiaya := r.FormValue("idKomponenBiaya")
	nominal := r.FormValue("nominal")
	idPegawai := r.FormValue("idPegawai")
	isReimbursement := r.FormValue("isReimbursement")
	boolValue, err := strconv.ParseBool(isReimbursement)
	if err != nil {
		log.Fatal(err)
	}
	keterangan := r.FormValue("keterangan")
	tanggal := r.FormValue("tanggal")
	jenis := r.FormValue("jenis")

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

	nom, err := decimal.NewFromString(nominal)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var reqFormat = bpd.PerjalananDinasBiayaDetail{
		ID:              id,
		IDBpdPegawai:    idBpdPegawai,
		IDJenisBiaya:    idJenisBiaya,
		IDKomponenBiaya: idKomponenBiaya,
		Nominal:         nom,
		IdPegawai:       idPegawai,
		IsReimbursement: boolValue,
		Keterangan:      keterangan,
		Tanggal:         &tanggal,
		File:            path,
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, err)
		return
	}
	newData, err := h.PerjalananDinasBiayaService.Create(reqFormat, userID, jenis)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// GetDataAll list all biaya perjalanan dinas.
// @Summary Get list all biaya perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data biaya perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai query string false "Set ID Perjalanan Dinas Pegawai"
// @Param idPegawai query string false "Set ID Pegawai"
// @Param isReimbursement query string false "Set is reimbursement"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/dto [get]
func (h *PerjalananDinasBiayaHandler) GetAllDto(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	idPegawai := r.URL.Query().Get("idPegawai")
	isReimbursement := r.URL.Query().Get("isReimbursement")
	data, err := h.PerjalananDinasBiayaService.GetBiayaDto(idBpdPegawai, idPegawai, isReimbursement)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data perjalanan dinas biaya.
// @Summary menghapus data perjalanan dinas biaya.
// @Description Endpoint ini adalah untuk menghapus data perjalanan dinas biaya.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/{id} [delete]
func (h *PerjalananDinasBiayaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, err)
		return
	}
	id, _ := uuid.FromString(chi.URLParam(r, "id"))
	err = h.PerjalananDinasBiayaService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, "success")
}

// create adalah untuk menambah data biaya perjalanan dinas.
// @Summary menambahkan data biaya perjalanan dinas.
// @Description Endpoint ini adalah untuk menambahkan data biaya perjalanan dinas.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string true "ID"
// @Param idBpdPegawai formData string true "ID BPD Pegawai"
// @Param idJenisBiaya formData string true "ID Jenis Biaya"
// @Param idKomponenBiaya formData string false "Id komponen Biaya"
// @Param nominal formData string false "Nominal"
// @Param idPegawai formData string false "Id Pegawai"
// @Param isReimbursement formData string false "Is reimbursement"
// @Param keterangan formData string false "keterangan"
// @Param jenis formData string false "jenis"
// @Param file formData file true "Foto"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasBiayaDetail}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/update [put]
func (h *PerjalananDinasBiayaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idBpdPegawai := r.FormValue("idBpdPegawai")
	idJenisBiaya := r.FormValue("idJenisBiaya")
	idKomponenBiaya := r.FormValue("idKomponenBiaya")
	nominal := r.FormValue("nominal")
	idPegawai := r.FormValue("idPegawai")
	isReimbursement := r.FormValue("isReimbursement")
	boolValue, err := strconv.ParseBool(isReimbursement)
	if err != nil {
		log.Fatal(err)
	}
	keterangan := r.FormValue("keterangan")
	tanggal := r.FormValue("tanggal")
	jenis := r.FormValue("jenis")

	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.PerjalananDinasService.UploadFile(w, r, "file", "", id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		ids, err := uuid.FromString(id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		pBiaya, err := h.PerjalananDinasBiayaService.ResolveByID(ids)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = pBiaya.File
	}

	nom, err := decimal.NewFromString(nominal)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var reqFormat = bpd.PerjalananDinasBiayaDetail{
		ID:              id,
		IDBpdPegawai:    idBpdPegawai,
		IDJenisBiaya:    idJenisBiaya,
		IDKomponenBiaya: idKomponenBiaya,
		Nominal:         nom,
		IdPegawai:       idPegawai,
		IsReimbursement: boolValue,
		Keterangan:      keterangan,
		Tanggal:         &tanggal,
		File:            path,
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, err)
		return
	}
	newData, err := h.PerjalananDinasBiayaService.Update(reqFormat, userID, jenis)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// GetDataAll list all biaya perjalanan dinas.
// @Summary Get list all biaya perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data biaya perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags biaya-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai query string false "Set ID Perjalanan Dinas Pegawai"
// @Param idPegawai query string false "Set ID Pegawai"
// @Param idJenisBiaya query string false "Set ID Jenis Biaya"
// @Param isReimbursement query string false "Set is reimbursement"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/biaya-perjalanan-dinas/histori-biaya [get]
func (h *PerjalananDinasBiayaHandler) GetHistoriBiaya(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	idPegawai := r.URL.Query().Get("idPegawai")
	idJenisBiaya := r.URL.Query().Get("idJenisBiaya")
	isReimbursement := r.URL.Query().Get("isReimbursement")
	data, err := h.PerjalananDinasBiayaService.GetHistoriBiaya(idBpdPegawai, idPegawai, idJenisBiaya, isReimbursement)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}
