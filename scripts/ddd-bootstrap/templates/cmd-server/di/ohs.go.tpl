package di

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"

	// ddd-bootstrap:import ohs
)

var OhsModule = fx.Options(
	fx.Provide(
		func() *fiber.App { return fiber.New() },
		func(app *fiber.App) fiber.Router { return app.Group("/api") },
		// ddd-bootstrap:provide ohs
	),
)
