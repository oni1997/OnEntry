package database

import (
	"context"
	"time"

	"github.com/oni1997/onentry/services/api-go/models"
)

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, master_key_salt, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.MasterKeySalt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, master_key_salt)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.MasterKeySalt).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (db *DB) GetVault(ctx context.Context, userID string) (*models.Vault, error) {
	query := `
		SELECT id, user_id, encrypted_vault, nonce, version, created_at, updated_at
		FROM vaults WHERE user_id = $1
	`
	vault := &models.Vault{}
	err := db.QueryRowContext(ctx, query, userID).Scan(
		&vault.ID, &vault.UserID, &vault.EncryptedVault, &vault.Nonce,
		&vault.Version, &vault.CreatedAt, &vault.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return vault, nil
}

func (db *DB) CreateVault(ctx context.Context, vault *models.Vault) error {
	query := `
		INSERT INTO vaults (user_id, encrypted_vault, nonce, version)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query, vault.UserID, vault.EncryptedVault, vault.Nonce, vault.Version).
		Scan(&vault.ID, &vault.CreatedAt, &vault.UpdatedAt)
}

func (db *DB) UpdateVault(ctx context.Context, vault *models.Vault) error {
	query := `
		UPDATE vaults SET encrypted_vault = $1, nonce = $2, version = version + 1, updated_at = $3
		WHERE user_id = $4 RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query, vault.EncryptedVault, vault.Nonce, time.Now(), vault.UserID).
		Scan(&vault.ID, &vault.CreatedAt, &vault.UpdatedAt)
}

func (db *DB) CreatePasswordEntry(ctx context.Context, entry *models.PasswordEntry) error {
	query := `
		INSERT INTO passwords (user_id, vault_id, title, username, encrypted_password, website, notes, folder, favorite, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		entry.UserID, entry.VaultID, entry.Title, entry.Username,
		entry.EncryptedPassword, entry.Website, entry.Notes,
		entry.Folder, entry.Favorite, entry.Tags,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

func (db *DB) UpdatePasswordEntry(ctx context.Context, entry *models.PasswordEntry) error {
	query := `
		UPDATE passwords SET title = $1, username = $2, encrypted_password = $3, website = $4,
		notes = $5, folder = $6, favorite = $7, tags = $8, updated_at = $9
		WHERE id = $10 AND user_id = $11 RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		entry.Title, entry.Username, entry.EncryptedPassword, entry.Website,
		entry.Notes, entry.Folder, entry.Favorite, entry.Tags,
		time.Now(), entry.ID, entry.UserID,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

func (db *DB) DeletePasswordEntry(ctx context.Context, userID, entryID string) error {
	query := `DELETE FROM passwords WHERE id = $1 AND user_id = $2`
	_, err := db.ExecContext(ctx, query, entryID, userID)
	return err
}

func (db *DB) GetPasswordEntries(ctx context.Context, userID string) ([]models.PasswordEntry, error) {
	query := `
		SELECT id, user_id, vault_id, title, username, encrypted_password, website,
		notes, folder, favorite, tags, created_at, updated_at
		FROM passwords WHERE user_id = $1 ORDER BY updated_at DESC
	`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.PasswordEntry
	for rows.Next() {
		var entry models.PasswordEntry
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.VaultID, &entry.Title, &entry.Username,
			&entry.EncryptedPassword, &entry.Website, &entry.Notes,
			&entry.Folder, &entry.Favorite, &entry.Tags, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (db *DB) SearchPasswordEntries(ctx context.Context, userID, query string) ([]models.PasswordEntry, error) {
	searchQuery := `
		SELECT id, user_id, vault_id, title, username, encrypted_password, website,
		notes, folder, favorite, tags, created_at, updated_at
		FROM passwords WHERE user_id = $1 AND (
			title ILIKE $2 OR username ILIKE $2 OR website ILIKE $2 OR $2 = ANY(tags)
		) ORDER BY updated_at DESC
	`
	rows, err := db.QueryContext(ctx, searchQuery, userID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.PasswordEntry
	for rows.Next() {
		var entry models.PasswordEntry
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.VaultID, &entry.Title, &entry.Username,
			&entry.EncryptedPassword, &entry.Website, &entry.Notes,
			&entry.Folder, &entry.Favorite, &entry.Tags, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (db *DB) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (user_id, action, resource_type, resource_id, ip_address, user_agent, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at
	`
	return db.QueryRowContext(ctx, query,
		log.UserID, log.Action, log.ResourceType, log.ResourceID,
		log.IPAddress, log.UserAgent, log.Metadata,
	).Scan(&log.ID, &log.CreatedAt)
}

func (db *DB) GetAuditLogs(ctx context.Context, userID string, limit int) ([]models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, ip_address, user_agent, metadata, created_at
		FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2
	`
	rows, err := db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&log.IPAddress, &log.UserAgent, &log.Metadata, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (db *DB) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, last_used_at
	`
	return db.QueryRowContext(ctx, query, session.UserID, session.TokenHash, session.RefreshTokenHash, session.ExpiresAt).
		Scan(&session.ID, &session.CreatedAt, &session.LastUsedAt)
}

func (db *DB) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token_hash, refresh_token_hash, expires_at, created_at, last_used_at
		FROM sessions WHERE token_hash = $1
	`
	session := &models.Session{}
	err := db.QueryRowContext(ctx, query, tokenHash).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.RefreshTokenHash,
		&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (db *DB) DeleteSession(ctx context.Context, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := db.ExecContext(ctx, query, sessionID)
	return err
}
