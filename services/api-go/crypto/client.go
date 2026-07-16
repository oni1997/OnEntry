package crypto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/oni1997/onentry/services/api-go/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) DeriveMasterKey(ctx context.Context, password string, salt string) (string, error) {
	req := map[string]string{"password": password, "salt": salt}
	resp, err := c.post(ctx, "/derive-master-key", req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		MasterKey string `json:"master_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.MasterKey, nil
}

func (c *Client) EncryptVault(ctx context.Context, plaintext string, key string) (*models.Vault, error) {
	req := map[string]string{"plaintext": plaintext, "key": key}
	resp, err := c.post(ctx, "/encrypt", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &models.Vault{
		EncryptedVault: []byte(result.Ciphertext),
		Nonce:          []byte(result.Nonce),
	}, nil
}

func (c *Client) DecryptVault(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error) {
	req := map[string]string{
		"ciphertext": string(ciphertext),
		"nonce":      string(nonce),
		"key":        key,
	}
	resp, err := c.post(ctx, "/decrypt", req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Plaintext string `json:"plaintext"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Plaintext, nil
}

func (c *Client) EncryptPassword(ctx context.Context, password string, key string) ([]byte, []byte, error) {
	req := map[string]string{"plaintext": password, "key": key}
	resp, err := c.post(ctx, "/encrypt", req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}
	return []byte(result.Ciphertext), []byte(result.Nonce), nil
}

func (c *Client) DecryptPassword(ctx context.Context, ciphertext []byte, nonce []byte, key string) (string, error) {
	return c.DecryptVault(ctx, ciphertext, nonce, key)
}

func (c *Client) GeneratePassword(ctx context.Context, req models.GeneratePasswordRequest) (string, error) {
	if req.Length == 0 {
		req.Length = 16
	}
	if !req.Uppercase && !req.Lowercase && !req.Numbers && !req.Symbols {
		req.Uppercase = true
		req.Lowercase = true
		req.Numbers = true
		req.Symbols = true
	}

	resp, err := c.post(ctx, "/generate-password", req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Password, nil
}

func (c *Client) HashPassword(ctx context.Context, password string) (string, string, error) {
	req := map[string]string{"password": password}
	resp, err := c.post(ctx, "/hash-password", req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Hash string `json:"hash"`
		Salt string `json:"salt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	return result.Hash, result.Salt, nil
}

func (c *Client) VerifyPassword(ctx context.Context, password, hash string) (bool, error) {
	req := map[string]string{"password": password, "hash": hash}
	resp, err := c.post(ctx, "/verify-password", req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Valid bool `json:"valid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.Valid, nil
}

func (c *Client) Health(ctx context.Context) error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("crypto service error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
