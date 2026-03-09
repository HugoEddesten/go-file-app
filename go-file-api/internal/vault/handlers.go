package vault

import (
	"errors"
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

func AssignUserToVault(vaultRepo *Repository, usersRepo *users.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		vaultId := locals.VaultId(c)

		body := new(VaultUserCreateRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		user, err := usersRepo.FindByEmail(body.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "No user with provided email found")
		}

		_, err = vaultRepo.AddUserToVault(ctx, vaultId, user.Id, body.Path, body.Role)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
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
