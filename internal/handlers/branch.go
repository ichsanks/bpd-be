package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type BranchHandler struct {
	BranchService master.BranchService
}

func ProvideBranchHandler(service master.BranchService) BranchHandler {
	return BranchHandler{
		BranchService: service,
	}
}

func (h *BranchHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/branch", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Post("/update", h.Update)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Delete("/{id}", h.Delete)
			r.Get("/{id}", h.ResolveByID)
		})
	})

}

// create adalah untuk menambah data branch.
// @Summary menambahkan data branch.
// @Description Endpoint ini adalah untuk menambahkan data branch.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param kode formData string false "kode"
// @Param nama formData string false "nama"
// @Param email formData string false "email"
// @Param address formData string false "address"
// @Param city formData string false "city"
// @Param contact formData string false "contact"
// @Param phone formData string false "phone"
// @Param website formData string false "website"
// @Param color formData string false "color"
// @Param umk formData string false "color"
// @Param isDark formData boolean false "color"
// @Param file formData file true "File Dokumen"
// @Success 200 {object} response.Base{data=master.Branch}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch [post]
func (h *BranchHandler) Create(w http.ResponseWriter, r *http.Request) {
	kode := r.FormValue("kode")
	nama := r.FormValue("nama")
	email := r.FormValue("email")
	address := r.FormValue("address")
	city := r.FormValue("city")
	contact := r.FormValue("contact")
	phone := r.FormValue("phone")
	website := r.FormValue("website")
	color := r.FormValue("color")
	isDark, _ := strconv.ParseBool(r.FormValue("isDark"))
	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := h.BranchService.UploadFile(w, r, "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}
	reqFormat := master.RequestBranchFormat{
		Kode:    kode,
		Nama:    nama,
		Email:   &email,
		Address: &address,
		City:    &city,
		Contact: &contact,
		Phone:   &phone,
		Website: &website,
		Color:   &color,
		IsDark:  isDark,
		Image:   &path,
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, err)
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, err)
		return
	}

	reqFormat.TenantID = tenantID
	data, err := h.BranchService.Create(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// Update adalah untuk mengubah data branch.
// @Summary mengubah data branch
// @Description Endpoint ini adalah untuk mengubah data branch.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id formData string false "ID BRANCH"
// @Param kode formData string false "kode"
// @Param nama formData string false "nama"
// @Param email formData string false "email"
// @Param address formData string false "address"
// @Param city formData string false "city"
// @Param contact formData string false "contact"
// @Param phone formData string false "phone"
// @Param website formData string false "website"
// @Param color formData string false "color"
// @Param isDark formData boolean false "color"
// @Param file formData file true "File Dokumen"
// @Success 200 {object} response.Base{data=[]master.Branch}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch/update [post]
func (h *BranchHandler) Update(w http.ResponseWriter, r *http.Request) {
	kode := r.FormValue("kode")
	nama := r.FormValue("nama")
	email := r.FormValue("email")
	address := r.FormValue("address")
	city := r.FormValue("city")
	contact := r.FormValue("contact")
	phone := r.FormValue("phone")
	website := r.FormValue("website")
	color := r.FormValue("color")
	isDark, _ := strconv.ParseBool(r.FormValue("isDark"))
	IDStr := r.FormValue("id")
	ID, err := uuid.FromString(IDStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	branch, err := h.BranchService.ResolveByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	uploadedFile, _, _ := r.FormFile("file")
	var path string
	images := branch.Image
	if uploadedFile != nil {
		filepath, err := h.BranchService.UploadFile(w, r, "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath

		// remove file
		if images != nil {
			err = h.BranchService.DeleteFile(*images)
			if err != nil {
				println("delete file", err)
			}
		}
	} else {
		if images != nil {
			path = *images
		}
	}

	reqFormat := master.RequestBranchFormat{
		Id:      ID,
		Kode:    kode,
		Nama:    nama,
		Email:   &email,
		Address: &address,
		City:    &city,
		Contact: &contact,
		Phone:   &phone,
		Website: &website,
		Color:   &color,
		IsDark:  isDark,
		Image:   &path,
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TenantID = tenantID
	data, err := h.BranchService.Update(reqFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveAll list all branch.
// @Summary Get list all branch.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data branch sesuai dengan filter yang dikirimkan.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch [get]
func (h *BranchHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	fmt.Println("pageSizeStr", pageSizeStr)
	fmt.Println("pageNumberStr", pageNumberStr)
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

	var idTenant string
	tenantID := middleware.GetClaimsValue(r.Context(), "tenantId")
	if tenantID != nil {
		idTenant = middleware.GetClaimsValue(r.Context(), "tenantId").(string)
	}

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		TenantID:   idTenant,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.BranchService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetDataAll list all branch.
// @Summary Get list all branch.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data branch sesuai dengan filter yang dikirimkan.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch/all [get]
func (h *BranchHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	// var idTenant string
	// tenantID := middleware.GetClaimsValue(r.Context(), "tenantId")
	// if tenantID != nil {
	// 	idTenant = middleware.GetClaimsValue(r.Context(), "tenantId").(string)
	// }

	status, err := h.BranchService.GetAllData()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// ResolveByID adalah untuk mendapatkan satu data branch berdasarkan ID.
// @Summary Mendapatkan satu data branch berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan branch By ID.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=master.Branch}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch/{id} [get]
func (h *BranchHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}

	data, err := h.BranchService.ResolveByIDDTO(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data branch.
// @Summary hapus data branch.
// @Description Endpoint ini adalah untuk menghapus data branch.
// @Tags branch
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/branch/{id} [delete]
func (h *BranchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newID, err := uuid.FromString(id)

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = h.BranchService.DeleteByID(newID, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}
