package vault

import (
	"errors"
	"fmt"
	"os"

	"go-file-api/internal/email"
	"go-file-api/internal/invites"
	"go-file-api/internal/locals"
	"go-file-api/internal/users"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetVault(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		ctx := c.UserContext()

		vault, err := vaultRepo.GetVault(ctx, vaultId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(vault)
	}
}

func GetUserVaults(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := locals.UserId(c)
		ctx := c.UserContext()

		vaults, err := vaultRepo.GetVaultsForUser(ctx, userId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(vaults)
	}
}

func CreateVault(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := locals.UserId(c)
		ctx := c.UserContext()

		vaults, err := vaultRepo.Create(ctx, "my_vault", userId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(vaults)
	}
}

func AssignUserToVault(vaultRepo *Repository, usersRepo *users.Repository, inviteRepo *invites.Repository, emailSvc email.EmailService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		vaultId := locals.VaultId(c)
		invitedBy := locals.UserId(c)

		body := new(VaultUserCreateRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}
		if body.Path == "" {
			body.Path = "/"
		}

		existingUser, err := usersRepo.FindByEmail(body.Email)
		if err == nil {
			// User exists — grant access directly.
			if _, err := vaultRepo.AddUserToVault(ctx, vaultId, existingUser.Id, body.Path, body.Role); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError)
			}
			vaultName, _ := vaultRepo.GetVaultName(ctx, vaultId)
			go emailSvc.SendVaultAccessGranted(ctx, body.Email, vaultName)
			return c.SendStatus(fiber.StatusCreated)
		}

		// User doesn't exist — create invite and send email.
		if errors.Is(err, pgx.ErrNoRows) {
			vaultName, err := vaultRepo.GetVaultName(ctx, vaultId)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Couldnt get vault name")
			}

			inv, err := inviteRepo.Create(ctx, vaultId, invitedBy, body.Email, body.Path, int(body.Role))
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "vault_invite creation failed")
			}

			appURL := os.Getenv("APP_URL")
			if appURL == "" {
				appURL = "http://localhost:5173"
			}
			inviteURL := fmt.Sprintf("%s/register/%s", appURL, inv.Token)
			go emailSvc.SendVaultInvite(ctx, body.Email, vaultName, inviteURL)
			return c.SendStatus(fiber.StatusCreated)
		}

		return fiber.NewError(fiber.StatusBadRequest, "Couldnt create the permission")
	}
}

func UpdateVaultUser(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		vaultId := locals.VaultId(c)
		adminUserId := locals.UserId(c)
		path := locals.RequestedVaultPath(c)

		body := new(VaultUserUpdateRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}
		if body.Role == VaultRoleOwner {
			return fiber.NewError(fiber.StatusBadRequest, "Cant set role to owner")
		}

		adminEntries, err := vaultRepo.GetVaultUsers(ctx, vaultId, adminUserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		target, err := vaultRepo.GetVaultUser(ctx, body.VaultUserId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fiber.NewError(fiber.StatusNotFound)
			}
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		editable := editableByAdmin(adminEntries, []VaultUser{*target})
		if len(editable) == 0 {
			return fiber.NewError(fiber.StatusForbidden, "Insufficient access level")
		}

		vaultUser := VaultUser{
			Id:   body.VaultUserId,
			Path: path,
			Role: body.Role,
		}

		_, err = vaultRepo.UpdateVaultUser(ctx, &vaultUser)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func RemoveUserFromVault(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		adminUserId := locals.UserId(c)
		vaultId := locals.VaultId(c)

		body := new(RemoveUserFromVaultRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		adminEntries, err := vaultRepo.GetVaultUsers(ctx, vaultId, adminUserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		targetEntries, err := vaultRepo.GetVaultUsers(ctx, vaultId, body.UserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		editable := editableByAdmin(adminEntries, targetEntries)
		if len(editable) == 0 {
			return c.JSON([]VaultUser{})
		}

		idsToDelete := make([]int, len(editable))
		for i, u := range editable {
			idsToDelete[i] = u.Id
		}

		deleted, err := vaultRepo.DeleteVaultUsersByIds(ctx, idsToDelete)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(deleted)
	}
}

func RemoveVaultUser(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		adminUserId := locals.UserId(c)
		vaultId := locals.VaultId(c)

		body := new(RemoveVaultUserRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		adminEntries, err := vaultRepo.GetVaultUsers(ctx, vaultId, adminUserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		target, err := vaultRepo.GetVaultUser(ctx, body.VaultUserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		editable := editableByAdmin(adminEntries, []VaultUser{*target})
		if len(editable) == 0 {
			return c.JSON([]VaultUser{})
		}

		deletedVaultUser, err := vaultRepo.DeleteVaultUsersByIds(ctx, []int{body.VaultUserId})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		return c.JSON(deletedVaultUser)
	}
}

// GetPendingInvites handles GET /vault/invites/:vaultId (admin)
func GetPendingInvites(inviteRepo *invites.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		vaultId := locals.VaultId(c)

		pending, err := inviteRepo.FindPendingByVault(ctx, vaultId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		return c.JSON(pending)
	}
}

// GetInviteInfo handles GET /invites/:token (public)
func GetInviteInfo(inviteRepo *invites.Repository, vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		token := c.Params("token")

		inv, err := inviteRepo.FindByToken(ctx, token)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fiber.NewError(fiber.StatusNotFound, "Invite not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		if inv.AcceptedAt != nil {
			return fiber.NewError(fiber.StatusGone, "Invite already accepted")
		}

		vaultName, err := vaultRepo.GetVaultName(ctx, inv.VaultId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(invites.InviteInfoResponse{
			Email:     inv.Email,
			VaultName: vaultName,
			Token:     inv.Token,
		})
	}
}
