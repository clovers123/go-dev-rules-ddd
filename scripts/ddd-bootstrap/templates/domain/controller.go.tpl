package {{.Domain}}

import (
	"github.com/gofiber/fiber/v3"

	appservice "{{.ModulePath}}/app/ohs/local/appservice/{{.Domain}}"
	ohspl "{{.ModulePath}}/app/ohs/pl/{{.Domain}}"
	pladapter "{{.ModulePath}}/app/ohs/pl/{{.Domain}}/adapter"
	resp "git.yugeeker.com/SHARED/go-lazy/app/entities"
)

type {{.DomainPascal}}Controller struct {
	adapter    pladapter.I{{.DomainPascal}}Adapter
	appService *appservice.{{.DomainPascal}}AppService
}

func New{{.DomainPascal}}Controller(adapter pladapter.I{{.DomainPascal}}Adapter, appService *appservice.{{.DomainPascal}}AppService) *{{.DomainPascal}}Controller {
	return &{{.DomainPascal}}Controller{
		adapter:    adapter,
		appService: appService,
	}
}

func (c *{{.DomainPascal}}Controller) Create(ctx fiber.Ctx) error {
	var req ohspl.Create{{.DomainPascal}}Request
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	cmd, err := c.adapter.FromCreate{{.DomainPascal}}DTO(ctx, req)
	if err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	bo, err := c.appService.Create(ctx.Context(), cmd)
	if err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	return ctx.JSON(resp.ResponseOK(bo))
}

func (c *{{.DomainPascal}}Controller) List(ctx fiber.Ctx) error {
	var req ohspl.List{{.DomainPascalPlural}}Request
	if err := ctx.Bind().Query(&req); err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	q, err := c.adapter.FromList{{.DomainPascalPlural}}DTO(ctx, req)
	if err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	data, err := c.appService.List(ctx.Context(), q)
	if err != nil {
		return ctx.JSON(resp.ResponseFromError(err))
	}
	return ctx.JSON(resp.ResponseOK(data))
}
