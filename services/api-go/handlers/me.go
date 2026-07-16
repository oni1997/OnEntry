package handlers

import (
	"net/http"

	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/middleware"
	"github.com/oni1997/onentry/services/api-go/models"
)

type MeHandler struct {
	db *database.DB
}

func NewMeHandler(db *database.DB) *MeHandler {
	return &MeHandler{db: db}
}

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetUserFromContext(r.Context())
	if session == nil {
		models.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), "")
	_ = user
	_ = err

	models.JSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"user_id": session.UserID},
	})
}
