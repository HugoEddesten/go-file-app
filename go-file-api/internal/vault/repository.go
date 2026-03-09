package vault

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) Create(ctx context.Context, name string, userId int) (*Vault, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var vault Vault
	err = tx.QueryRow(ctx,
		`INSERT INTO vaults (name)
		 VALUES ($1)
		 RETURNING id, name, created_at, updated_at`,
		name,
	).Scan(
		&vault.Id,
		&vault.Name,
		&vault.CreatedAt,
		&vault.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO vault_users (vault_id, user_id, path, role)
		 VALUES ($1, $2, $3, $4)`,
		vault.Id,
		userId,
		"/",
		VaultRoleOwner,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *Repository) AddUserToVault(
	ctx context.Context,
	vaultId int,
	userId int,
	path string,
	role VaultRole,
) (*VaultUser, error) {

	var vaultUser VaultUser

	err := r.DB.QueryRow(ctx,
		`INSERT INTO vault_users (vault_id, user_id, path, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, vault_id, user_id, path, role, created_at, updated_at`,
		vaultId,
		userId,
		path,
		role,
	).Scan(
		&vaultUser.Id,
		&vaultUser.VaultId,
		&vaultUser.UserId,
		&vaultUser.Path,
		&vaultUser.Role,
		&vaultUser.CreatedAt,
		&vaultUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &vaultUser, nil
}

func (r *Repository) UpdateVaultUser(ctx context.Context, vu *VaultUser) (*VaultUser, error) {
	var updated VaultUser

	err := r.DB.QueryRow(ctx, `
		UPDATE vault_users
		SET 
			role = COALESCE($1, role),
			path = COALESCE($2, path),
			updated_at = NOW()
		WHERE id = $3
		RETURNING id, vault_id, user_id, path, role, created_at, updated_at
	`, vu.Role, vu.Path, vu.Id).Scan(
		&updated.Id,
		&updated.VaultId,
		&updated.UserId,
		&updated.Path,
		&updated.Role,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *Repository) GetVaultUsers(ctx context.Context, vaultId int, userId int) ([]VaultUser, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, vault_id, user_id, path, role
		FROM vault_users
		WHERE vault_id = $1 AND user_id = $2
	`, vaultId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VaultUser
	for rows.Next() {
		var vu VaultUser
		if err := rows.Scan(
			&vu.Id,
			&vu.VaultId,
			&vu.UserId,
			&vu.Path,
			&vu.Role,
		); err != nil {
			return nil, err
		}
		result = append(result, vu)
	}

	return result, rows.Err()
}

func (r *Repository) GetVault(ctx context.Context, vaultId int) (*VaultWithUsers, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT
			v.id,
			v.name,
			u.id,
			u.email,
			vu.role,
			vu.path,
			vu.id
		FROM vaults v
		JOIN vault_users vu ON vu.vault_id = v.id
		JOIN users u ON u.id = vu.user_id
		WHERE v.id = $1
		ORDER BY u.id
	`, vaultId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vault *VaultWithUsers

	for rows.Next() {
		var (
			vaultId int
			v       VaultWithUsers
			u       UserInVault
		)

		if err := rows.Scan(
			&vaultId,
			&v.Name,
			&u.Id,
			&u.Email,
			&u.Role,
			&u.Path,
			&u.VaultUserId,
		); err != nil {
			return nil, err
		}

		if vault == nil {
			v.Id = vaultId
			v.Users = []UserInVault{}
			vault = &v
		}

		vault.Users = append(vault.Users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if vault == nil {
		return nil, pgx.ErrNoRows
	}

	return vault, nil
}

func (r *Repository) DeleteVaultUsersByIds(ctx context.Context, ids []int) ([]VaultUser, error) {
	rows, err := r.DB.Query(ctx, `
		DELETE FROM vault_users
		WHERE id = ANY($1)
		RETURNING id, vault_id, user_id, path, role, created_at, updated_at
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VaultUser
	for rows.Next() {
		var vu VaultUser
		if err := rows.Scan(
			&vu.Id,
			&vu.VaultId,
			&vu.UserId,
			&vu.Path,
			&vu.Role,
			&vu.CreatedAt,
			&vu.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, vu)
	}

	return result, rows.Err()
}

func (r *Repository) GetVaultsForUser(
	ctx context.Context,
	userId int,
) ([]VaultWithUsers, error) {

	rows, err := r.DB.Query(ctx, `
		SELECT
			v.id,
			v.name,
			u.id,
			u.email,
			vu.role,
			vu.path,
			vu.id
		FROM vaults v
		JOIN vault_users vu ON vu.vault_id = v.id
		JOIN users u ON u.id = vu.user_id
		WHERE v.id IN (
			SELECT vault_id
			FROM vault_users
			WHERE user_id = $1
		)
		ORDER BY v.id
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vaultMap := make(map[int]*VaultWithUsers)

	for rows.Next() {
		var (
			vaultId int
			v       VaultWithUsers
			u       UserInVault
		)

		if err := rows.Scan(
			&vaultId,
			&v.Name,
			&u.Id,
			&u.Email,
			&u.Role,
			&u.Path,
			&u.VaultUserId,
		); err != nil {
			return nil, err
		}

		existing, ok := vaultMap[vaultId]
		if !ok {
			v.Id = vaultId
			v.Users = []UserInVault{}
			vaultMap[vaultId] = &v
			existing = &v
		}

		existing.Users = append(existing.Users, u)
	}

	result := make([]VaultWithUsers, 0, len(vaultMap))
	for _, v := range vaultMap {
		result = append(result, *v)
	}

	return result, rows.Err()
}

func (r *Repository) GetVaultName(ctx context.Context, vaultId int) (string, error) {
	var name string
	err := r.DB.QueryRow(ctx, `SELECT name FROM vaults WHERE id = $1`, vaultId).Scan(&name)
	return name, err
}

func (r *Repository) GetVaultUser(ctx context.Context, vaultUserId int) (*VaultUser, error) {
	var vaultUser VaultUser

	err := r.DB.QueryRow(ctx, `
		SELECT
			vu.id,
			vu.path,
			vu.user_id,
			vu.vault_id,
			vu.created_at,
			vu.updated_at,
			vu.role
		FROM vault_users vu
		WHERE vu.id = $1`, vaultUserId).Scan(
		&vaultUser.Id,
		&vaultUser.Path,
		&vaultUser.UserId,
		&vaultUser.VaultId,
		&vaultUser.CreatedAt,
		&vaultUser.UpdatedAt,
		&vaultUser.Role,
	)

	if err != nil {
		return nil, err
	}

	return &vaultUser, nil
}
