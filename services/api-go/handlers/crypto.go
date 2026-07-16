package handlers

import (
	"context"

	"github.com/oni1997/onentry/services/api-go/models"
)

type CryptoClient interface {
	HashPassword(ctx context.Context, password string) (string, string, error)
	VerifyPassword(ctx context.Context, password, hash string) (bool, error)
	DeriveMasterKey(ctx context.Context, password string, salt string) (string, error)
	GeneratePassword(ctx context.Context, req models.GeneratePasswordRequest) (string, error)
	EncryptPassword(ctx context.Context, password string, key string) ([]byte, []byte, error)
	DecryptPassword(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error)
	EncryptVault(ctx context.Context, plaintext string, key string) (*models.Vault, error)
	DecryptVault(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error)
}
