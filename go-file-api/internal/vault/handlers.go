package vault

import (
	"go-file-api/internal/locals"
	"go-file-api/internal/users"

	"github.com/gofiber/fiber/v2"
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
		path := locals.RequestedVaultPath(c)

		body := new(VaultUserUpdateRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		vaultUser := VaultUser{
			Id:   body.VaultUserId,
			Path: path,
			Role: body.Role,
		}

		_, err := vaultRepo.UpdateVaultUser(ctx, &vaultUser)
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

		var idsToDelete []int
		for _, target := range targetEntries {
			if target.Role == VaultRoleOwner {
				continue
			}
			for _, admin := range adminEntries {
				if pathAllowed(admin.Path, target.Path) {
					idsToDelete = append(idsToDelete, target.Id)
					break
				}
			}
		}

		if len(idsToDelete) == 0 {
			return c.JSON([]VaultUser{})
		}

		deleted, err := vaultRepo.DeleteVaultUsersByIds(ctx, idsToDelete)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(deleted)
	}
}
