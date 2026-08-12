package adapter

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"git.yugeeker.com/SHARED/go-lazy/tools"
	"{{.ModulePath}}/app"
	"{{.ModulePath}}/app/domain/{{.Domain}}/entity"
	ohspl "{{.ModulePath}}/app/ohs/pl/{{.Domain}}"
)

type I{{.DomainPascal}}Adapter interface {
	FromCreate{{.DomainPascal}}DTO(ctx fiber.Ctx, dto ohspl.Create{{.DomainPascal}}Request) (cmd ohspl.Create{{.DomainPascal}}Command, err error)
	FromList{{.DomainPascalPlural}}DTO(ctx fiber.Ctx, req ohspl.List{{.DomainPascalPlural}}Request) (q ohspl.List{{.DomainPascalPlural}}Query, err error)

	FromCreate{{.DomainPascal}}Command(ctx context.Context, cmd ohspl.Create{{.DomainPascal}}Command) (ent *entity.{{.DomainPascal}}, err error)

	From{{.DomainPascal}}Entity(ctx context.Context, ent *entity.{{.DomainPascal}}) (bo ohspl.{{.DomainPascal}}BO, err error)
	From{{.DomainPascal}}Entities(ctx context.Context, entities []*entity.{{.DomainPascal}}) (bo ohspl.{{.DomainPascalPlural}}BO, err error)
}

type {{.DomainPascal}}Adapter struct{}

func New{{.DomainPascal}}Adapter() I{{.DomainPascal}}Adapter {
	return &{{.DomainPascal}}Adapter{}
}

func (o {{.DomainPascal}}Adapter) FromCreate{{.DomainPascal}}DTO(ctx fiber.Ctx, dto ohspl.Create{{.DomainPascal}}Request) (cmd ohspl.Create{{.DomainPascal}}Command, err error) {
	u, err := app.ExtractorUserInfo(ctx)
	if err != nil {
		return
	}
	cmd, err = tools.ConvertObject[ohspl.Create{{.DomainPascal}}Request, ohspl.Create{{.DomainPascal}}Command](&dto)
	if err != nil {
		return
	}
	cmd.User = u
	return
}

func (o {{.DomainPascal}}Adapter) FromList{{.DomainPascalPlural}}DTO(ctx fiber.Ctx, req ohspl.List{{.DomainPascalPlural}}Request) (q ohspl.List{{.DomainPascalPlural}}Query, err error) {
	u, err := app.ExtractorUserInfo(ctx)
	if err != nil {
		return
	}
	q, err = tools.ConvertObject[ohspl.List{{.DomainPascalPlural}}Request, ohspl.List{{.DomainPascalPlural}}Query](&req)
	if err != nil {
		return
	}
	q.User = u
	return
}

func (o {{.DomainPascal}}Adapter) FromCreate{{.DomainPascal}}Command(ctx context.Context, cmd ohspl.Create{{.DomainPascal}}Command) (ent *entity.{{.DomainPascal}}, err error) {
	ent, err = tools.ConvertObjectPtr[ohspl.Create{{.DomainPascal}}Command, entity.{{.DomainPascal}}](&cmd)
	if err != nil {
		return
	}
	ent.OwnedBy = cmd.User.AccountID
	return
}

func (o {{.DomainPascal}}Adapter) From{{.DomainPascal}}Entity(ctx context.Context, ent *entity.{{.DomainPascal}}) (bo ohspl.{{.DomainPascal}}BO, err error) {
	return tools.ConvertObject[entity.{{.DomainPascal}}, ohspl.{{.DomainPascal}}BO](ent)
}

func (o {{.DomainPascal}}Adapter) From{{.DomainPascal}}Entities(ctx context.Context, entities []*entity.{{.DomainPascal}}) (bo ohspl.{{.DomainPascalPlural}}BO, err error) {
	for _, ent := range entities {
		boItem, err := o.From{{.DomainPascal}}Entity(ctx, ent)
		if err != nil {
			return bo, err
		}
		bo.Items = append(bo.Items, boItem)
	}
	bo.Total = int64(len(entities))
	return bo, nil
}
