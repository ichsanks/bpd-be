package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

// RoleHandler adalah HTTP handler untuk domain Role
type RoleHandler struct {
	RoleService auth.RoleService
}

// ProvideRoleHandler adalah provider untuk handler ini
func ProvideRoleHandler(roleService auth.RoleService) RoleHandler {
	return RoleHandler{
		RoleService: roleService,
	}
}

// Router untuk setup dari router untuk domain ini
func (h *RoleHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/roles", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.GetData)
			r.Get("/all", h.resolveAll)
			r.Get("/{id}", h.resolveRoleByID)
			r.Get("/menus/{roleId}", h.resolveMenusByRoleID)
			r.Post("/", h.createRole)
			r.Put("/{id}", h.updateRole)
			r.Delete("/{id}", h.deleteRole)
			r.Put("/create-or-update", h.createOrUpdateRole)
		})
	})
}

// ResolveAll list all role.
// @Summary Get list all role.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data role sesuai dengan filter yang dikirimkan.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama | description ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} auth.Role
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles [get]
func (h *RoleHandler) GetData(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "nama"
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

	status, err := h.RoleService.GetData(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// resolveAll adalah untuk mendapatkan semua data Role.
// @Summary Mendapatkan semua data Role.
// @Description Endpoint ini adalah untuk mendapatkan semua data role.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base{data=auth.Role}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/all [get]
func (h *RoleHandler) resolveAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.RoleService.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

// creteRole adalah untuk menambah data Role.
// @Summary menambahkan data Role.
// @Description Endpoint ini adalah untuk menambahkan data role.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param role body auth.RequestRole true "Role yang akan ditambahkan"
// @Success 200 {object} response.Base{data=auth.Role}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles [post]
func (h *RoleHandler) createRole(w http.ResponseWriter, r *http.Request) {
	var newRole auth.Role
	err := json.NewDecoder(r.Body).Decode(&newRole)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	exist, err := h.RoleService.CreateRole(newRole)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	if exist {
		response.WithError(w, failure.Conflict("createROle", "role", err.Error()))
	}
	response.WithJSON(w, http.StatusCreated, newRole)
}

// creteRole adalah untuk mengubah data Role.
// @Summary update data Role.
// @Description Endpoint ini adalah untuk mengubah data role.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Role"
// @Param role body auth.RequestRole true "Role yang akan diedit"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/{id} [put]
func (h *RoleHandler) updateRole(w http.ResponseWriter, r *http.Request) {
	var newRole auth.Role
	err := json.NewDecoder(r.Body).Decode(&newRole)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(newRole)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	id := chi.URLParam(r, "id")

	err = h.RoleService.UpdateRole(id, newRole)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, "success")
}

// deleteRole adalah untuk menghapus data Role.
// @Summary hapus data Role.
// @Description Endpoint ini adalah untuk menghapus data role.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Role"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/{id} [delete]
func (h *RoleHandler) deleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.RoleService.DeleteRole(id)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, "success")
}

// resolveRoleByID adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "RoleID"
// @Success 200 {object} response.Base{data=auth.RoleMenuFormat}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/{id} [get]
func (h *RoleHandler) resolveRoleByID(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	log.Info().Msg("RoleID: " + roleID)
	role, err := h.RoleService.ResolveByRoleID(roleID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, role)
}

// resolveRoleByID adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId path string true "RoleID"
// @Success 200 {object} response.Base{data=[]string}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/menus/{roleId} [get]
func (h *RoleHandler) resolveMenusByRoleID(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleId")
	log.Info().Msg("RoleID: " + roleID)
	menus, err := h.RoleService.ResolveMenuByRoleID(roleID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menus)
}

// createOrUpdateRole adalah untuk menambahkan atau mengubah data Role dan Menu yang diperbolehkan.
// @Summary update data RoleMenu.
// @Description Endpoint ini adalah untuk mengubah data menu berdasarkan role.
// @Tags roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param menus body auth.RoleMenuFormat true "Menu yang di set untuk Role ini"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/roles/create-or-update [put]
func (h *RoleHandler) createOrUpdateRole(w http.ResponseWriter, r *http.Request) {
	var roleMenu auth.RoleMenuFormat
	err := json.NewDecoder(r.Body).Decode(&roleMenu)
	err = h.RoleService.CreateOrUpdateMenuByRoleID(roleMenu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, "success")
}
