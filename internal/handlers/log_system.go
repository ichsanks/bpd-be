package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/auth"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/middleware"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type LogSystemHandler struct {
	LogSystemService auth.LogSystemService
}

func ProvideLogSystemHandler(LogSystemService auth.LogSystemService) LogSystemHandler {
	return LogSystemHandler{
		LogSystemService: LogSystemService,
	}
}

func (h *LogSystemHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/log-system", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.CreateLogSystem)
		})
	})
}

// CreateLogSystem adalah untuk menambah data Log System.
// @Summary menambahkan data Log System.
// @Description Endpoint ini adalah untuk menambahkan data Log System.
// @Tags log-system
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param logSystem body auth.RequestLogSystemFormat true "Log System yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.LogSystem}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/log-system [post]
func (h *LogSystemHandler) CreateLogSystem(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestLogSystemFormat
	fmt.Println("reqFormat", reqFormat)
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	newMenu, err := h.LogSystemService.CreateLogSystem(reqFormat, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenu)
}
