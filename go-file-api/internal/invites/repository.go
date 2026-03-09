package invites

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Repository) Create(ctx context.Context, vaultId, invitedBy int, email, path string, role int) (*VaultInvite, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	var inv VaultInvite
	err = r.DB.QueryRow(ctx, `
		INSERT INTO vault_invites (vault_id, invited_by, email, role, path, token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, vault_id, invited_by, email, role, path, token, expires_at, accepted_at, created_at
	`, vaultId, invitedBy, email, role, path, token, expiresAt).Scan(
		&inv.Id, &inv.VaultId, &inv.InvitedBy, &inv.Email, &inv.Role, &inv.Path,
		&inv.Token, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt,
	)

	return &inv, err
}

func (r *Repository) FindByToken(ctx context.Context, token string) (*VaultInvite, error) {
	var inv VaultInvite
	err := r.DB.QueryRow(ctx, `
		SELECT id, vault_id, invited_by, email, role, path, token, expires_at, accepted_at, created_at
		FROM vault_invites
		WHERE token = $1
	`, token).Scan(
		&inv.Id, &inv.VaultId, &inv.InvitedBy, &inv.Email, &inv.Role, &inv.Path,
		&inv.Token, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt,
	)
	return &inv, err
}

func (r *Repository) FindPendingByEmail(ctx context.Context, email string) ([]VaultInvite, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, vault_id, invited_by, email, role, path, token, expires_at, accepted_at, created_at
		FROM vault_invites
		WHERE email = $1 AND accepted_at IS NULL AND expires_at > NOW()
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VaultInvite
	for rows.Next() {
		var inv VaultInvite
		if err := rows.Scan(
			&inv.Id, &inv.VaultId, &inv.InvitedBy, &inv.Email, &inv.Role, &inv.Path,
			&inv.Token, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, inv)
	}
	return result, rows.Err()
}

func (r *Repository) FindPendingByVault(ctx context.Context, vaultId int) ([]VaultInvite, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, vault_id, invited_by, email, role, path, token, expires_at, accepted_at, created_at
		FROM vault_invites
		WHERE vault_id = $1 AND accepted_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
	`, vaultId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VaultInvite
	for rows.Next() {
		var inv VaultInvite
		if err := rows.Scan(
			&inv.Id, &inv.VaultId, &inv.InvitedBy, &inv.Email, &inv.Role, &inv.Path,
			&inv.Token, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, inv)
	}
	return result, rows.Err()
}

func (r *Repository) Accept(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE vault_invites SET accepted_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM vault_invites WHERE id = $1`, id)
	return err
}
