package {{.Domain}}

import (
	"context"

	"{{.ModulePath}}/app/domain/{{.Domain}}/entity"
)

type {{.DomainPascal}}Repository interface {
	Create(ctx context.Context, item *entity.{{.DomainPascal}}) error
	List(ctx context.Context, ownedBy string, keyword string, limit int, offset int) ([]*entity.{{.DomainPascal}}, error)
}
