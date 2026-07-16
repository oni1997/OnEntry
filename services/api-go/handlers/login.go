package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oni1997/onentry/services/api-go/config"
	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/models"
	"github.com/oni1997/onentry/services/api-go/utils"
)

type LoginHandler struct {
	db     *database.DB
	crypto CryptoClient
	cfg    *config.Config
}

func NewLoginHandler(db *database.DB, crypto CryptoClient, cfg *config.Config) *LoginHandler {
	return &LoginHandler{db: db, crypto: crypto, cfg: cfg}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		models.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		models.JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	valid, err := h.crypto.VerifyPassword(r.Context(), req.Password, user.PasswordHash)
	if err != nil || !valid {
		models.JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	masterKey, err := h.crypto.DeriveMasterKey(r.Context(), req.Password, user.MasterKeySalt)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to derive master key")
		return
	}

	accessToken, err := h.generateAccessToken(user, masterKey)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	session := &models.Session{
		UserID:           user.ID,
		TokenHash:        utils.HashSHA256(accessToken),
		RefreshTokenHash: utils.HashSHA256(refreshToken),
		ExpiresAt:        time.Now().Add(time.Duration(h.cfg.AccessTokenExpiry) * time.Minute),
	}

	if err := h.db.CreateSession(r.Context(), session); err != nil {
		models.JSONError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	go func() {
		_ = h.db.CreateAuditLog(context.Background(), &models.AuditLog{
			UserID:      user.ID,
			Action:      "login",
			ResourceType: "session",
			ResourceID:  &session.ID,
			IPAddress:   r.RemoteAddr,
			UserAgent:   r.UserAgent(),
		})
	}()

	models.JSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user": map[string]interface{}{
				"id":            user.ID,
				"email":         user.Email,
				"master_key":    masterKey,
			},
			"expires_in": h.cfg.AccessTokenExpiry * 60,
		},
	})
}

func (h *LoginHandler) generateAccessToken(user *models.User, masterKey string) (string, error) {
	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(h.cfg.AccessTokenExpiry) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
