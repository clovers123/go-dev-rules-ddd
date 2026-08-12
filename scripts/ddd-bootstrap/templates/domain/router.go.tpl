package routers

import (
	"github.com/gofiber/fiber/v3"

	controller "{{.ModulePath}}/app/ohs/remote/controller/{{.Domain}}"
)

func Mount{{.DomainPascal}}RouteGroup(router fiber.Router, controller *controller.{{.DomainPascal}}Controller) {
	group := router.Group("/v1/{{.DomainPlural}}")
	group.Post("/", controller.Create)
	group.Get("/", controller.List)
}
