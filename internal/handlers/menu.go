package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

// MenuHandler adalah HTTP handler untuk domain Role
type MenuHandler struct {
	MenuService auth.MenuService
}

// ProvideMenuHandler adalah provider untuk handler ini
func ProvideMenuHandler(MenuService auth.MenuService) MenuHandler {
	return MenuHandler{
		MenuService: MenuService,
	}
}

// Router untuk setup dari router untuk domain ini
func (h *MenuHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/menu-user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveMenuByRoleID)
			r.Get("/permission", h.ResolveAllMenuUser)
			r.Post("/", h.CreateMenuUser)
			r.Delete("/{id}", h.DeleteMenuUser)
			r.Put("/sort", h.SortMenu)
		})
	})
	r.Route("/menu", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllMenu)
			r.Get("/{id}", h.ResolveMenuByID)
			r.Post("/", h.CreateMenu)
			r.Put("/{id}", h.UpdateMenu)
			r.Delete("/{id}", h.DeleteMenu)
		})
	})
}

// CreateMenuUser adalah untuk menambah data MenuUser.
// @Summary menambahkan data MenuUser.
// @Description Endpoint ini adalah untuk menambahkan data MenuUser.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param MenuUser body auth.RequestMenuUserFormat true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.MenuUser}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-user [post]
func (h *MenuHandler) CreateMenuUser(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestMenuUserFormat
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.TenantID = &tenantID
	fmt.Println("reqFormat", reqFormat)
	newMenuUser, err := h.MenuService.CreateMenuUser(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenuUser)
}

// DeleteMenuUser adalah untuk menghapus data MenuUser.
// @Summary hapus data MenuUser.
// @Description Endpoint ini adalah untuk menghapus data MenuUser.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-user/{id} [delete]
func (h *MenuHandler) DeleteMenuUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newID, err := uuid.FromString(id)
	fmt.Println("ID:", newID)
	err = h.MenuService.DeleteMenuUserByID(newID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Data Berhasil di Hapus",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// MASTER MENU

// GetAllMenu list all menu.
// @Summary Get list all menu.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data menu sesuai dengan filter yang dikirimkan.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param app query string false "Set App is one of [1=SILATURAHMI | 2=CPL]"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/all [get]
func (h *MenuHandler) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	status, err := h.MenuService.GetAllMenu(app)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// ResolveAll list all Menu.
// @Summary Get list all Menu.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Menu sesuai dengan filter yang dikirimkan.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [id | namaMenu ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param app query string false "Set App is one of [1=SILATURAHMI | 2=CPL]"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu [get]
func (h *MenuHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	app := r.URL.Query().Get("app")
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

	req := model.StandardRequestMenu{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		App:        app,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.MenuService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// CreateMenu adalah untuk menambah data Menu.
// @Summary menambahkan data Menu.
// @Description Endpoint ini adalah untuk menambahkan data Menu.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Menu body auth.RequestMenuFormat true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu [post]
func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestMenuFormat
	fmt.Println("reqFormat", reqFormat)
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	newMenu, err := h.MenuService.CreateMenu(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenu)
}

// ResolveMenuByID adalah untuk mendapatkan satu data Menu berdasarkan ID.
// @Summary Mendapatkan satu data Menu berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Menu ID.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [get]
func (h *MenuHandler) ResolveMenuByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	menu, err := h.MenuService.ResolveMenuByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, menu)
}

// UpdateMenu adakan untuk update data Menu
// @Summary update data Menu
// @Description endpoint ini adalah untuk mengubah data Menu
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Param Menu body auth.RequestMenuFormat true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [put]
func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	fmt.Println("harusnya oke")

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
	}

	var newMenu auth.RequestMenuFormat
	err = json.NewDecoder(r.Body).Decode(&newMenu)

	menu, err := h.MenuService.UpdateMenu(id, newMenu)
	if err != nil {
		log.Info().Msg("Error: " + err.Error())
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, menu)
}

// DeleteMenu adalah untuk menghapus data Menu.
// @Summary hapus data Menu.
// @Description Endpoint ini adalah untuk menghapus data Menu.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [delete]
func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newID, err := uuid.FromString(id)

	err = h.MenuService.DeleteMenuByID(newID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Data Berhasil di Hapus",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// resolveMenuByID adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId query string true "Set RoleID"
// @Param app query string false "Set App is one of [1=SILATURAHMI | 2=CPL]"
// @Param posisi query string false "Set Posisi is one of [1,2,3]"
// @Param posisiSubMenu query string false "Set Posisi sub menu is one of [1,2,3]"
// @Param idBidang query string false "Set IdBidang"
// @Param idBranch query string false "Set IdBranch"
// @Param idTenant query string false "Set idTenant"
// @Success 200 {object} response.Base{data=auth.MenuResponse}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-user [get]
func (h *MenuHandler) ResolveMenuByRoleID(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	app := r.URL.Query().Get("app")
	posisi := r.URL.Query().Get("posisi")
	posisiSubMenu := r.URL.Query().Get("posisiSubMenu")
	idBidang := r.URL.Query().Get("idBidang")

	if posisi == "" {
		posisi = "1"
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

	var idTenant string
	tenantID := middleware.GetClaimsValue(r.Context(), "tenantId")
	if tenantID != nil {
		idTenant = middleware.GetClaimsValue(r.Context(), "tenantId").(string)
	}

	idTenantStr := r.URL.Query().Get("idTenant")
	if idTenantStr != "" {
		idTenant = idTenantStr
	}

	var req = auth.RequestMenuUserFilter{
		IdRole:        roleID,
		App:           app,
		Posisi:        posisi,
		PosisiSubMenu: posisiSubMenu,
		Level:         "1",
		IdBidang:      idBidang,
		IDBranch:      idBranch,
		TenantID:      idTenant,
	}

	menu, err := h.MenuService.ResolveMenuByRoleID(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// ResolveAllMenuUser adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId query string true "Set RoleID"
// @Param app query string false "Set App is one of [1=SILATURAHMI | 2=CPL]"
// @Param posisi query string false "Set Posisi is one of [1,2,3]"
// @Param level query string false "Set Level is one of [1,2,3]"
// @Param linkParent query string false "Set Link parent"
// @Param idBranch query string false "Set Id Branch"
// @Success 200 {object} response.Base{data=auth.MenuResponse}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-user/permission [get]
func (h *MenuHandler) ResolveAllMenuUser(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	app := r.URL.Query().Get("app")
	posisi := r.URL.Query().Get("posisi")
	level := r.URL.Query().Get("level")
	linkParent := r.URL.Query().Get("linkParent")

	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	idBranchStr := r.URL.Query().Get("idBranch")
	if idBranchStr != "" {
		idBranch = idBranchStr
	}

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var req = auth.RequestMenuUserFilter{
		IdRole:     roleID,
		App:        app,
		Posisi:     posisi,
		Level:      level,
		LinkParent: linkParent,
		IDBranch:   idBranch,
		TenantID:   tenantID.String(),
	}

	menu, err := h.MenuService.ResolveAllMenuUser(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// SortingMenu adakan untuk mengurutkan Menu
// @Summary Sort data menu
// @Description endpoint ini adalah untuk mengurutkan Menu
// @Tags menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Menu body auth.RequestMenuSortFormat true "Menu yang akan diurutkan UP||DOWN"
// @Success 200 {object} response.Base{data=[]auth.MenuUser}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-user/sort [put]
func (h *MenuHandler) SortMenu(w http.ResponseWriter, r *http.Request) {
	var sortMenu auth.RequestMenuSortFormat
	err := json.NewDecoder(r.Body).Decode(&sortMenu)

	menu, err := h.MenuService.SortMenu(sortMenu)
	if err != nil {
		log.Info().Msg("Error: " + err.Error())
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, menu)
}
