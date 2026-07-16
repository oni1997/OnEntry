package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID              string    `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	MasterKeySalt   string    `json:"-" db:"master_key_salt"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Vault struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	EncryptedVault []byte    `json:"-" db:"encrypted_vault"`
	Nonce          []byte    `json:"-" db:"nonce"`
	Version        int       `json:"version" db:"version"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type PasswordEntry struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	VaultID          string    `json:"vault_id" db:"vault_id"`
	Title            string    `json:"title" db:"title"`
	Username         string    `json:"username" db:"username"`
	EncryptedPassword []byte   `json:"-" db:"encrypted_password"`
	Website          string    `json:"website" db:"website"`
	Notes            string    `json:"notes" db:"notes"`
	Folder           string    `json:"folder" db:"folder"`
	Favorite         bool      `json:"favorite" db:"favorite"`
	Tags             []string  `json:"tags" db:"tags"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	TokenHash       string    `json:"-" db:"token_hash"`
	RefreshTokenHash string   `json:"-" db:"refresh_token_hash"`
	ExpiresAt       time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	LastUsedAt      time.Time `json:"last_used_at" db:"last_used_at"`
}

type AuditLog struct {
	ID          string                 `json:"id" db:"id"`
	UserID      string                 `json:"user_id" db:"user_id"`
	Action      string                 `json:"action" db:"action"`
	ResourceType string                `json:"resource_type" db:"resource_type"`
	ResourceID  *string                `json:"resource_id" db:"resource_id"`
	IPAddress   string                 `json:"ip_address" db:"ip_address"`
	UserAgent   string                 `json:"user_agent" db:"user_agent"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	MasterKey   string `json:"master_key" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PasswordEntryRequest struct {
	Title       string   `json:"title" validate:"required"`
	Username    string   `json:"username"`
	Password    string   `json:"password" validate:"required"`
	Website     string   `json:"website"`
	Notes       string   `json:"notes"`
	Folder      string   `json:"folder"`
	Favorite    bool     `json:"favorite"`
	Tags        []string `json:"tags"`
}

type GeneratePasswordRequest struct {
	Length          int    `json:"length" validate:"min=4,max=128"`
	Uppercase       bool   `json:"uppercase"`
	Lowercase       bool   `json:"lowercase"`
	Numbers         bool   `json:"numbers"`
	Symbols         bool   `json:"symbols"`
	ExcludeSimilar  bool   `json:"exclude_similar"`
	Pronounceable   bool   `json:"pronounceable"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
