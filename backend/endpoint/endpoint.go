package endpoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/strawstacks/strawhouse/strawhouse-backend/endpoint/get"
	"github.com/strawstacks/strawhouse/strawhouse-backend/endpoint/system"
	"github.com/strawstacks/strawhouse/strawhouse-backend/type/response"
)

func HandleRoot(c *fiber.Ctx) error {
	return c.JSON(response.Success("STRAWHOUSE BACKEND"))
}

func HandleNotFound(c *fiber.Ctx) error {
	return fiber.ErrNotFound
}

func Bind(app *fiber.App, systemHandler *system.Handler, getHandler *get.Handler) {
	// * System route
	systemRoute := app.Group("_")
	systemRoute.Post("upload", systemHandler.Upload)
	systemRoute.Use(HandleNotFound)

	// * Get route
	app.Get("/", HandleRoot)
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.Redirect("https://static1.pixcee.dev/external/strawstacks/favicon.ico", fiber.StatusMovedPermanently)
	})
	app.Get("*", getHandler.Get)
	app.Use(HandleNotFound)
}