package vault

import "github.com/gofiber/fiber/v2"

func GetVault(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
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
		userId := c.Locals("userId").(int)
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
		userId := c.Locals("userId").(int)
		ctx := c.UserContext()

		vaults, err := vaultRepo.Create(ctx, "my_vault", userId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.JSON(vaults)
	}
}

func AssignUserToVault(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		userId := c.Locals("userId").(int)
		vaultId := c.Locals("vaultId").(int)

		body := new(VaultUserRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		_, err := vaultRepo.AddUserToVault(ctx, vaultId, userId, body.Path, body.Role)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
	}
}

func UpdateVaultUser(vaultRepo *Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		vaultId := c.Locals("vaultId").(int)
		path := c.Locals("requestedVaultPath").(string)

		body := new(VaultUserRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		vaultUser := VaultUser{
			Id:   vaultId,
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
