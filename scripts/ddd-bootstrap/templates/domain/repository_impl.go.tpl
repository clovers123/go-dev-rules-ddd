package {{.Domain}}

import (
	"context"

	"{{.ModulePath}}/app/acl/adapter/repository/postgres/repository"
	aclpl "{{.ModulePath}}/app/acl/pl/{{.Domain}}"
	portrepository "{{.ModulePath}}/app/acl/port/repository/{{.Domain}}"
	"{{.ModulePath}}/app/domain/{{.Domain}}/entity"
)

var _ portrepository.{{.DomainPascal}}Repository = (*{{.DomainPascal}}Repository)(nil)

type {{.DomainPascal}}Repository struct {
	adapter aclpl.I{{.DomainPascal}}AclAdapter
}

func New{{.DomainPascal}}Repository(adapter aclpl.I{{.DomainPascal}}AclAdapter) portrepository.{{.DomainPascal}}Repository {
	return &{{.DomainPascal}}Repository{adapter: adapter}
}

func (r *{{.DomainPascal}}Repository) Create(ctx context.Context, item *entity.{{.DomainPascal}}) error {
	dao := repository.Get{{.TablePascal}}Dao(ctx, nil)
	record := r.adapter.From{{.DomainPascal}}AggregationToInfo(item)
	return dao.Create(record)
}

func (r *{{.DomainPascal}}Repository) List(ctx context.Context, ownedBy string, keyword string, limit int, offset int) ([]*entity.{{.DomainPascal}}, error) {
	dao := repository.Get{{.TablePascal}}Dao(ctx, nil)
	// The GORM Gen SQL keeps owned_by and is_deleted=false constraints in the query path.
	records, err := dao.FindPage(ownedBy, keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	return r.adapter.To{{.DomainPascal}}Aggregations(records), nil
}
