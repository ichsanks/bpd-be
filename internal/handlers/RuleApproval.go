package handlers

import (
	"encoding/json"
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

type RuleApprovalHandler struct {
	RuleApprovalService master.RuleApprovalService
}

func ProvideRuleApprovalHandler(service master.RuleApprovalService) RuleApprovalHandler {
	return RuleApprovalHandler{
		RuleApprovalService: service,
	}
}

func (h *RuleApprovalHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/rule-approval", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Put("/", h.Update)
			r.Get("/{id}", h.ResolveByIDDTO)
			r.Delete("/{id}", h.Delete)
			r.Get("/ttd/{id}", h.ResolveTTD)
		})
	})
}

// ResolveAll list all Rule Approval.
// @Summary Get list all Rule Approval.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Rule Approval sesuai dengan filter yang dikirimkan.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param jenis query string false "Set Jenis"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval [get]
func (h *RuleApprovalHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	jenis := r.URL.Query().Get("jenis")
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

	var idBranchs string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranchs = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	req := model.StandardRequestRuleApproval{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		Jenis:      jenis,
		IdBranch:   idBranchs,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.RuleApprovalService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// create adalah untuk menambah data Rule Approval.
// @Summary menambahkan data Rule Approval.
// @Description Endpoint ini adalah untuk menambahkan data Rule Approval.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RuleApproval body master.RuleApprovalRequest true "Rule Approval yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.RuleApproval}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval [post]
func (h *RuleApprovalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RuleApprovalRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		fmt.Print("error jsondecoder")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)

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

	data, err := h.RuleApprovalService.Create(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error create")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// update adalah untuk menambah data Rule Approval.
// @Summary menambahkan data Rule Approval.
// @Description Endpoint ini adalah untuk menambahkan data Rule Approval.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RuleApproval body master.RuleApprovalRequest true "Rule Approval yang akan ditambahkan"
// @Success 200 {object} response.Base{data=master.RuleApproval}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval [put]
func (h *RuleApprovalHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqFormat master.RuleApprovalRequest
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)

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

	data, err := h.RuleApprovalService.Update(reqFormat, userID, tenantID)
	if err != nil {
		fmt.Print("error update")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// ResolveByIDDTO adalah untuk mendapatkan satu data Rule Approval berdasarkan ID.
// @Summary Mendapatkan satu data Rule Approval berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Rule Approval By ID.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval/{id} [get]
func (h *RuleApprovalHandler) ResolveByIDDTO(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	data, err := h.RuleApprovalService.ResolveByIDDTO(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// delete adalah untuk menghapus data Rule Approval.
// @Summary hapus data Rule Approval.
// @Description Endpoint ini adalah untuk menghapus data Rule Approval.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval/{id} [delete]
func (h *RuleApprovalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetClaimsValue(r.Context(), "userId").(string)
	err := h.RuleApprovalService.DeleteByID(id, userID)
	if err != nil {
		fmt.Println(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// GetDataAll list all Rule Approval.
// @Summary Get list all Rule Approval.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Rule Approval sesuai dengan filter yang dikirimkan.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param idFungsionalitas query string false "Set ID Fungsionalitas"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval/all [get]
func (h *RuleApprovalHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	idFungsionalitas := r.URL.Query().Get("idFungsionalitas")
	data, err := h.RuleApprovalService.GetAll(idFungsionalitas)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveTTD adalah untuk mendapatkan satu data Rule Approval berdasarkan ID.
// @Summary Mendapatkan satu data Rule Approval berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Rule Approval By ID.
// @Tags rule-approval
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/rule-approval/ttd/{id} [get]
func (h *RuleApprovalHandler) ResolveTTD(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	data, err := h.RuleApprovalService.ResolveTtd(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}
