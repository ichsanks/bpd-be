package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/gofrs/uuid"
)

// UserHandler the HTTP handler for User domain.
type UserHandler struct {
	UserService auth.UserService
	Config      *configs.Config
}

// ProvideUserHandler is the provider for this handler.
func ProvideUserHandler(userService auth.UserService, config *configs.Config) UserHandler {
	return UserHandler{
		UserService: userService,
		Config:      config,
	}
}

// Router sets up the router for this domain.
func (u *UserHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", u.Login)
		r.Post("/validasi-login", u.ValidasiLogin)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", u.CreateUser)
			r.Put("/{id}", u.UpdateUser)
			r.Put("/fcm-token/{id}", u.UpdateUserFcmToken)
			r.Delete("/{id}", u.DeleteUser)
			r.Get("/", u.ResolveAll)
			r.Get("/{id}", u.ResolveUserById)
			r.Put("/password/{id}", u.ChangePassword)
			r.Put("/password/pw/{id}", u.ChangePassword)
			r.Put("/password/reset/{id}", u.ResetPassword)
		})
	})
}

// ValidasiLogin sign in a user
// @Summary sign in a user
// @Description This endpoint sign in a user
// @Tags users
// @Param users body auth.InputLogin true "The User to be sign in."
// @Produce json
// @Success 201 {object} response.Base{auth.ResponseLogin}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/validasi-login [post]
func (u *UserHandler) ValidasiLogin(w http.ResponseWriter, r *http.Request) {
	var input auth.InputLogin
	fmt.Println("INPUT:", input)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	resp, exist, err := u.UserService.ValidasiLogin(input)
	if !exist {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

// Login sign in a user
// @Summary sign in a user
// @Description This endpoint sign in a user
// @Tags users
// @Param users body auth.InputLogin true "The User to be sign in."
// @Produce json
// @Success 201 {object} response.Base{auth.ResponseLogin}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/login [post]
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.InputLogin
	fmt.Println("INPUT:", input)

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	fmt.Println("ip address : ", r.Header.Get("x-forwarded-for"))

	// if input.CapchaTOken == "" {
	// 	return
	// }

	// if err := u.verifyRecaptcha(input.CapchaTOken); err != nil {
	// 	response.WithError(w, failure.BadRequest(err))
	// 	return
	// }

	resp, exist, _, err := u.UserService.Login(input, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if !exist {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

func (u *UserHandler) verifyRecaptcha(response string) error {

	ScreetKey := u.Config.App.APIExternal.RecaptchaSecretKey
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"secret":   ScreetKey,
			"response": response,
		}).
		Get("https://www.google.com/recaptcha/api/siteverify")

	if err != nil {
		return err
	}

	var recaptchaResponse auth.RecaptchaResponse
	err = json.Unmarshal(resp.Body(), &recaptchaResponse)
	if err != nil {
		return err
	}

	if !recaptchaResponse.Success {
		return fmt.Errorf("reCAPTCHA verification failed: %v", recaptchaResponse.ErrorCodes)
	}

	return nil
}

// CreateUser creates a new user
// @Summary Create a new User.
// @Description This endpoint creates a new User.
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param users body auth.InputUser true "The User to be created."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user [post]
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var input auth.InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = input.Validate()
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))

	if err != nil {
		fmt.Print("error userId")
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

	if input.IDBranch == nil {
		input.IDBranch = &idBranch
	}

	input.TenantID = &tenantID
	exist, err := u.UserService.CreateUser(input, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if exist {
		response.WithError(w, failure.Conflict("register", "user", err.Error()))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// ChangePassword update user password
// @Summary update user password
// @Description This endpoint to update user password
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "The User identifier."
// @Param users body auth.InputChangePassword true "The User update a new password."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/password/{id} [put]
func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var input auth.InputChangePassword

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.ChangePassword(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Password berhasil diperbarui",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// ResetPassword reset user password
// @Summary reset user password
// @Description This endpoint to reset user password
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "The User identifier."
// @Param users body auth.InputChangePassword true "The User reset a new password."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/password/reset/{id} [put]
func (u *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input auth.InputChangePassword

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.ResetPassword(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Password berhasil diperbarui",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// ResolveAll list all user.
// @Summary Get list all user.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data user sesuai dengan filter yang dikirimkan.
// @Tags users
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [id | kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param idRole query string false "id role"
// @Param idUnor query string false "id unor"
// @Param idBidang query string false "id bidang"
// @Success 200 {object} auth.UserDTO
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user [get]
func (h *UserHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
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

	idRole := r.URL.Query().Get("idRole")
	idUnor := r.URL.Query().Get("idUnor")
	idBidang := r.URL.Query().Get("idBidang")

	tenantID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "tenantId").(string))
	if err != nil {
		fmt.Print("error tenant id")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var branchIDStr string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		branchIDStr = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	req := model.StandardRequestUser{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		IdRole:     idRole,
		IdUnor:     idUnor,
		IdBidang:   idBidang,
		TenantID:   tenantID.String(),
		IdBranch:   branchIDStr,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.UserService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// UpdateUser update user data
// @Summary update user data
// @Description This endpoint to update user entity
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "The User identifier."
// @Param users body auth.UserUpdateFormat true "The User update data"
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [put]
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input auth.UserUpdateFormat

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	var idBranch string
	branchID := middleware.GetClaimsValue(r.Context(), "branchId")
	if branchID != nil {
		idBranch = middleware.GetClaimsValue(r.Context(), "branchId").(string)
	}

	if input.IDBranch == nil {
		input.IDBranch = &idBranch
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	exist, err := u.UserService.UpdateUser(id, input, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())

	if exist {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// UpdateUserFcmToken update data fcm token user
// @Summary update data fcm token user
// @Description This endpoint to update user entity
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "The User identifier."
// @Param users body auth.UserUpdateFcmTokenFormat true "The User update Fcm Token data"
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/fcm-token/{id} [put]
func (u *UserHandler) UpdateUserFcmToken(w http.ResponseWriter, r *http.Request) {

	var input auth.UserUpdateFcmTokenFormat

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.UpdateUserFcmToken(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// UpdateUser delete user data
// @Summary delete user data
// @Description This endpoint to delete user entity
// @Tags users
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "The User identifier."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [delete]
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveUserByID adalah untuk mendapatkan satu data User berdasarkan ID.
// @Summary Mendapatkan satu data User berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan User ID.
// @Tags users
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=auth.User}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [get]
func (h *UserHandler) ResolveUserById(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	user, err := h.UserService.ResolveUserById(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, user)
}
