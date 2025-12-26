package vault

import (
	"time"
)

type Vault struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VaultUser struct {
	Id        int
	Path      string
	UserId    int
	VaultId   int
	Role      VaultRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VaultWithUsers struct {
	Id    int           `json:"id"`
	Name  string        `json:"name"`
	Users []UserInVault `json:"users"`
}

type UserInVault struct {
	Id    int       `json:"id"`
	Email string    `json:"email"`
	Role  VaultRole `json:"role"`
	Path  string    `json:"path"`
}

type VaultUserCreateRequest struct {
	Role  VaultRole `json:"role"`
	Path  string    `json:"path"`
	Email string    `json:"email"`
}

type VaultUserUpdateRequest struct {
	Role        VaultRole `json:"role"`
	Path        string    `json:"path"`
	VaultUserId int       `json:"vaultUserId"`
}

type PathBodyValidation struct {
	Path string `json:"path"`
}

type VaultRole int

const (
	VaultRoleOwner VaultRole = iota + 1
	VaultRoleAdmin
	VaultRoleEditor
	VaultRoleViewer
)
