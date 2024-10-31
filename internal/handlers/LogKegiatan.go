package handlers

import (
	"fmt"
	"image/jpeg"
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/bpd"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type LogKegiatanHandler struct {
	LogKegiatanService bpd.LogKegiatanService
	Config             *configs.Config
}

func ProvideLogKegiatanHandler(service bpd.LogKegiatanService, config *configs.Config) LogKegiatanHandler {
	return LogKegiatanHandler{
		LogKegiatanService: service,
		Config:             config,
	}
}

func (h *LogKegiatanHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/bpd/log-kegiatan", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllData)
			r.Get("/generate", h.GenerateTextToImage)
			r.Post("/", h.Create)
			r.Delete("/{id}", h.Delete)
		})
	})

	// Perjalanan dinas dokumen
	r.Route("/bpd/dokumen", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/all", h.GetAllDataDokumen)
			r.Post("/", h.CreateDokumen)
			r.Delete("/{id}", h.DeleteDokumen)
		})
	})
}

// GetDataAll list all log kegiatan.
// @Summary Get list all log kegiatan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data log kegiatan sesuai dengan filter yang dikirimkan.
// @Tags log-kegiatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPerjalananDinas query string false "id perjalanan dinas"
// @Param idBpdPegawai query string false "id bpd pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/log-kegiatan/all [get]
func (h *LogKegiatanHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	idPerjalananDinas := r.URL.Query().Get("idPerjalananDinas")
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.LogKegiatanService.GetAll(idPerjalananDinas, idBpdPegawai)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data log kegiatan.
// @Summary menambahkan data log kegiatan.
// @Description Endpoint ini adalah untuk menambahkan data log kegiatan.
// @Tags log-kegiatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idPerjalananDinas formData string true "ID Perjalanan Dinas"
// @Param idBpdPegawai formData string true "ID Bpd Pegawai"
// @Param tanggal formData string false "Tanggal"
// @Param keterangan formData string false "Keterangan"
// @Param file formData file true "Foto"
// @Success 200 {object} response.Base{data=bpd.LogKegiatan}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/log-kegiatan [post]
func (h *LogKegiatanHandler) Create(w http.ResponseWriter, r *http.Request) {
	idPerjalanananDinas := r.FormValue("idPerjalananDinas")
	idBpdPegawai := r.FormValue("idBpdPegawai")
	tanggal := r.FormValue("tanggal")
	keterangan := r.FormValue("keterangan")
	lat := r.FormValue("lat")
	long := r.FormValue("long")
	address := r.FormValue("address")

	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.LogKegiatanService.UploadFile(w, r, "file", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.LogKegiatanRequest{
		Tanggal:           tanggal,
		IdPerjalananDinas: idPerjalanananDinas,
		IdBpdPegawai:      &idBpdPegawai,
		Foto:              path,
	}

	if keterangan != "" {
		reqFormat.Keterangan = &keterangan
	}

	if lat != "" {
		reqFormat.Lat = &lat
	}

	if long != "" {
		reqFormat.Long = &long
	}

	if address != "" {
		reqFormat.Address = &address
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.LogKegiatanService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// delete adalah untuk menghapus data log kegiatan.
// @Summary menghapus data log kegiatan.
// @Description Endpoint ini adalah untuk menghapus data log kegiatan.
// @Tags log-kegiatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/log-kegiatan/{id} [delete]
func (h *LogKegiatanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.LogKegiatanService.DeleteByID(id)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// GetDataAll list all dokumen perjalanan dinas.
// @Summary Get list all dokumen perjalanan dinas.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data dokumen perjalanan dinas sesuai dengan filter yang dikirimkan.
// @Tags dokumen-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai query string false "id bpd pegawai"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/dokumen/all [get]
func (h *LogKegiatanHandler) GetAllDataDokumen(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.URL.Query().Get("idBpdPegawai")
	data, err := h.LogKegiatanService.GetAllDokumen(idBpdPegawai)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data perjalanan dinas dokumen.
// @Summary menambahkan data perjalanan dinas dokumen.
// @Description Endpoint ini adalah untuk menambahkan data perjalanan dinas dokumen.
// @Tags dokumen-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idBpdPegawai formData string true "ID Perjalanan Dinas"
// @Param keterangan formData string false "Keterangan"
// @Param file formData file true "File"
// @Success 200 {object} response.Base{data=bpd.PerjalananDinasDokumen}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/dokumen [post]
func (h *LogKegiatanHandler) CreateDokumen(w http.ResponseWriter, r *http.Request) {
	idBpdPegawai := r.FormValue("idBpdPegawai")
	keterangan := r.FormValue("keterangan")

	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.LogKegiatanService.UploadFileDokumen(w, r, "file", "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	var reqFormat = bpd.PerjalananDinasDokumenRequest{
		IdBpdPegawai: idBpdPegawai,
		File:         path,
	}

	if keterangan != "" {
		reqFormat.Keterangan = &keterangan
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	newData, err := h.LogKegiatanService.CreateDokumen(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newData)
}

// delete adalah untuk menghapus data perjalanan dinas dokumen.
// @Summary menghapus data perjalanan dinas dokumen.
// @Description Endpoint ini adalah untuk menghapus data perjalanan dinas dokumen.
// @Tags dokumen-perjalanan-dinas
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/dokumen/{id} [delete]
func (h *LogKegiatanHandler) DeleteDokumen(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete")
	id := chi.URLParam(r, "id")
	err := h.LogKegiatanService.DeleteDokumen(id)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// GetDataAll list all log kegiatan.
// @Summary Get list all log kegiatan.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data log kegiatan sesuai dengan filter yang dikirimkan.
// @Tags log-kegiatan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/bpd/log-kegiatan/generate [get]
func (h *LogKegiatanHandler) GenerateTextToImage(w http.ResponseWriter, r *http.Request) {
	lat := "-7.3531392"
	long := "112.7907328"
	address := "Pondok Tjandra Indah, Tambak Rejo, Waru, Sidoarjo, East Java, Java, 60294, Indonesia"
	var reqFormat = bpd.LogKegiatanRequest{
		Tanggal: "28-11-2023 16:34:00",
		Lat:     &lat,
		Long:    &long,
		Address: &address,
		Foto:    "/files/foto_dinas/foto_dinas_1f6eaed7-55f4-469b-905e-ff4f1477b1d0.jpg",
	}
	data, err := h.LogKegiatanService.GenerateTextToImage(reqFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}
	jpeg.Encode(w, data, nil)
}
