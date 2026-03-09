package invites

import "time"

type VaultInvite struct {
	Id         int
	VaultId    int
	InvitedBy  int
	Email      string
	Role       int
	Path       string
	Token      string
	ExpiresAt  time.Time
	AcceptedAt *time.Time
	CreatedAt  time.Time
}

type InviteInfoResponse struct {
	Email     string `json:"email"`
	VaultName string `json:"vaultName"`
	Token     string `json:"token"`
}
