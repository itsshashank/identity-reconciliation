package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/identity-reconciliation/db"
)

func NewServer(store db.UserStorer) *fiber.App {
	app := fiber.New()
	userHandler := NewUserHandler(store)
	app.Post("/order", userHandler.HandlePostOrder)
	app.Post("/identify", userHandler.HandlePostIdentify)
	return app
}
