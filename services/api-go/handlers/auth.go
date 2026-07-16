package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/models"
	"github.com/oni1997/onentry/services/api-go/utils"
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		models.JSONError(w, http.StatusServiceUnavailable, "Database unhealthy")
		return
	}
	models.JSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"status": "healthy"},
	})
}

type RegisterHandler struct {
	db     *database.DB
	crypto CryptoClient
}

type CryptoClient interface {
	HashPassword(ctx context.Context, password string) (string, string, error)
	DeriveMasterKey(ctx context.Context, password string, salt string) (string, error)
}

func NewRegisterHandler(db *database.DB, crypto CryptoClient) *RegisterHandler {
	return &RegisterHandler{db: db, crypto: crypto}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		models.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.ValidateEmail(req.Email) {
		models.JSONError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	existing, _ := h.db.GetUserByEmail(r.Context(), req.Email)
	if existing != nil {
		models.JSONError(w, http.StatusConflict, "Email already registered")
		return
	}

	passwordHash, salt, err := h.crypto.HashPassword(r.Context(), req.Password)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	_, err = h.crypto.DeriveMasterKey(r.Context(), req.Password, salt)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to derive master key")
		return
	}

	user := &models.User{
		Email:         req.Email,
		PasswordHash:  passwordHash,
		MasterKeySalt: salt,
	}

	if err := h.db.CreateUser(r.Context(), user); err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	go func() {
		auditLog := &models.AuditLog{
			UserID:      user.ID,
			Action:      "register",
			ResourceType: "user",
			ResourceID:  &user.ID,
			IPAddress:   r.RemoteAddr,
			UserAgent:   r.UserAgent(),
		}
		_ = h.db.CreateAuditLog(context.Background(), auditLog)
	}()

	models.JSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    map[string]string{"user_id": user.ID},
	})
}
