package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"{{.ModulePath}}/app/domain/{{.Domain}}/entity"
	portrepository "{{.ModulePath}}/app/acl/port/repository/{{.Domain}}"
)

type I{{.DomainPascal}}DomainService interface {
	Create{{.DomainPascal}}(ctx context.Context, item *entity.{{.DomainPascal}}) error
}

type {{.DomainPascal}}DomainService struct {
	repo portrepository.{{.DomainPascal}}Repository
}

func New{{.DomainPascal}}DomainService(repo portrepository.{{.DomainPascal}}Repository) *{{.DomainPascal}}DomainService {
	return &{{.DomainPascal}}DomainService{repo: repo}
}

func (svc *{{.DomainPascal}}DomainService) Create{{.DomainPascal}}(ctx context.Context, ent *entity.{{.DomainPascal}}) error {
	ent.ID = uuid.NewString()
	ent.CreateTime = time.Now()
	ent.UpdateTime = time.Now()
	return svc.repo.Create(ctx, ent)
}
