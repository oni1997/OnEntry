package handlers

import (
	"context"
	"net/http"

	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/middleware"
	"github.com/oni1997/onentry/services/api-go/models"
	"github.com/oni1997/onentry/services/api-go/utils"
)

type GenerateHandler struct {
	db     *database.DB
	crypto CryptoClient
}

type CryptoClient interface {
	GeneratePassword(ctx context.Context, req models.GeneratePasswordRequest) (string, error)
	EncryptPassword(ctx context.Context, password string, key string) ([]byte, []byte, error)
}

func NewGenerateHandler(db *database.DB, crypto CryptoClient) *GenerateHandler {
	return &GenerateHandler{db: db, crypto: crypto}
}

func (h *GenerateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		models.JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.GeneratePasswordRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		models.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	password, err := h.crypto.GeneratePassword(r.Context(), req)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to generate password")
		return
	}

	models.JSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"password": password},
	})
}
