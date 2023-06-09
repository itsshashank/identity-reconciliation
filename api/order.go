package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/identity-reconciliation/db"
	"github.com/itsshashank/identity-reconciliation/types"
)

type UserHandler struct {
	userStore db.UserStorer
}

func NewUserHandler(s db.UserStorer) *UserHandler {
	return &UserHandler{
		userStore: s,
	}
}

func (h *UserHandler) HandlePostOrder(c *fiber.Ctx) error {

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

	if err := h.userStore.PutOrder(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
	})
}
