package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

func DecodeJSON(r interface{ ReadJSON(interface{}) error }, v interface{}) error {
	if err := r.ReadJSON(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashSHA256(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func GetClientIP(r interface{ RemoteAddr string; Header Getter }) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	return r.RemoteAddr
}

type Getter interface {
	Get(string) string
}

var passwordPattern = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*]{8,}$`)

func ValidatePasswordStrength(password string) bool {
	return len(password) >= 8
}

func RequestID() string {
	return middleware.RequestIDHeader
}
