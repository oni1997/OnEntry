package handlers

import (
	"context"
	"net/http"

	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/middleware"
	"github.com/oni1997/onentry/services/api-go/models"
	"github.com/oni1997/onentry/services/api-go/utils"
)

type VaultHandler struct {
	db     *database.DB
	crypto CryptoClient
}

type CryptoClient interface {
	DecryptVault(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error)
	EncryptVault(ctx context.Context, plaintext string, key string) (*models.Vault, error)
	EncryptPassword(ctx context.Context, password string, key string) ([]byte, []byte, error)
	DecryptPassword(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error)
}

func NewVaultHandler(db *database.DB, crypto CryptoClient) *VaultHandler {
	return &VaultHandler{db: db, crypto: crypto}
}

func (h *VaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetUserFromContext(r.Context())
	if session == nil {
		models.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getVault(w, r, session)
	case http.MethodPost:
		h.createOrUpdateVault(w, r, session)
	default:
		models.JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *VaultHandler) getVault(w http.ResponseWriter, r *http.Request, session *models.Session) {
	masterKey := r.URL.Query().Get("master_key")
	if masterKey == "" {
		models.JSONError(w, http.StatusBadRequest, "master_key is required")
		return
	}

	vault, err := h.db.GetVault(r.Context(), session.UserID)
	if err != nil {
		models.JSONError(w, http.StatusNotFound, "Vault not found")
		return
	}

	plaintext, err := h.crypto.DecryptVault(r.Context(), vault.EncryptedVault, vault.Nonce, masterKey)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to decrypt vault")
		return
	}

	models.JSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"vault": plaintext},
	})
}

func (h *VaultHandler) createOrUpdateVault(w http.ResponseWriter, r *http.Request, session *models.Session) {
	var req struct {
		Vault    string `json:"vault" validate:"required"`
		MasterKey string `json:"master_key" validate:"required"`
	}
	if err := utils.DecodeJSON(r, &req); err != nil {
		models.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	vault, err := h.crypto.EncryptVault(r.Context(), req.Vault, req.MasterKey)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to encrypt vault")
		return
	}
	vault.UserID = session.UserID

	existing, _ := h.db.GetVault(r.Context(), session.UserID)
	if existing != nil {
		vault.ID = existing.ID
		if err := h.db.UpdateVault(r.Context(), vault); err != nil {
			models.JSONError(w, http.StatusInternalServerError, "Failed to update vault")
			return
		}
	} else {
		if err := h.db.CreateVault(r.Context(), vault); err != nil {
			models.JSONError(w, http.StatusInternalServerError, "Failed to create vault")
			return
		}
	}

	models.JSON(w, http.StatusOK, models.APIResponse{Success: true})
}
