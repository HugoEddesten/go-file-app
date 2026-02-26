package locals

import "github.com/gofiber/fiber/v2"

func UserId(c *fiber.Ctx) int                { return c.Locals("userId").(int) }
func Email(c *fiber.Ctx) string              { return c.Locals("email").(string) }
func VaultId(c *fiber.Ctx) int               { return c.Locals("vaultId").(int) }
func VaultRole(c *fiber.Ctx) int             { return c.Locals("vaultRole").(int) }
func RequestedVaultPath(c *fiber.Ctx) string { return c.Locals("requestedVaultPath").(string) }

func SetUserId(c *fiber.Ctx, v int)                { c.Locals("userId", v) }
func SetEmail(c *fiber.Ctx, v string)              { c.Locals("email", v) }
func SetVaultId(c *fiber.Ctx, v int)               { c.Locals("vaultId", v) }
func SetVaultRole(c *fiber.Ctx, v int)             { c.Locals("vaultRole", v) }
func SetRequestedVaultPath(c *fiber.Ctx, v string) { c.Locals("requestedVaultPath", v) }
