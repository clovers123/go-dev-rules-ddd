package {{.Domain}}

import (
	"context"

	domainservice "{{.ModulePath}}/app/domain/{{.Domain}}/service"
	ohspl "{{.ModulePath}}/app/ohs/pl/{{.Domain}}"
	pladapter "{{.ModulePath}}/app/ohs/pl/{{.Domain}}/adapter"
	portrepository "{{.ModulePath}}/app/acl/port/repository/{{.Domain}}"
)

type {{.DomainPascal}}AppService struct {
	domainService domainservice.I{{.DomainPascal}}DomainService
	repo          portrepository.{{.DomainPascal}}Repository
	adapter       pladapter.I{{.DomainPascal}}Adapter
}

func New{{.DomainPascal}}AppService(
	domainService domainservice.I{{.DomainPascal}}DomainService,
	repo portrepository.{{.DomainPascal}}Repository,
	adapter pladapter.I{{.DomainPascal}}Adapter,
) *{{.DomainPascal}}AppService {
	return &{{.DomainPascal}}AppService{
		domainService: domainService,
		repo:          repo,
		adapter:       adapter,
	}
}

func (s *{{.DomainPascal}}AppService) Create(ctx context.Context, cmd ohspl.Create{{.DomainPascal}}Command) (bo ohspl.{{.DomainPascal}}BO, err error) {
	item, err := s.adapter.FromCreate{{.DomainPascal}}Command(ctx, cmd)
	if err != nil {
		return
	}

	err = s.domainService.Create{{.DomainPascal}}(ctx, item)
	if err != nil {
		return
	}

	return s.adapter.From{{.DomainPascal}}Entity(ctx, item)
}

func (s *{{.DomainPascal}}AppService) List(ctx context.Context, query ohspl.List{{.DomainPascalPlural}}Query) (bo ohspl.{{.DomainPascalPlural}}BO, err error) {
	items, err := s.repo.List(ctx, query.User.AccountID, query.Keyword, query.Limit, query.Offset)
	if err != nil {
		return
	}

	return s.adapter.From{{.DomainPascal}}Entities(ctx, items)
}
