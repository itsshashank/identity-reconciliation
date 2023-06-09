package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/identity-reconciliation/types"
)

func (h *UserHandler) HandlePostIdentify(c *fiber.Ctx) error {

	var req types.Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Email == "" && req.PhoneNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email or phoneNumber parameter is required",
		})
	}

	response, err := h.userStore.GetContacts(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"contact": response,
	})
}
