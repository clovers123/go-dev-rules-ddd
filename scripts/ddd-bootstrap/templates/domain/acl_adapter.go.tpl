package {{.Domain}}

import (
	"time"

	"{{.ModulePath}}/app/acl/adapter/repository/postgres/model"
	"{{.ModulePath}}/app/domain/{{.Domain}}/entity"
)

type I{{.DomainPascal}}AclAdapter interface {
	// To{{.DomainPascal}}Aggregation maps DB model to domain entity
	To{{.DomainPascal}}Aggregation(record *model.{{.TablePascal}}) *entity.{{.DomainPascal}}
	// To{{.DomainPascal}}Aggregations maps a slice of DB models to domain entities
	To{{.DomainPascal}}Aggregations(records []*model.{{.TablePascal}}) []*entity.{{.DomainPascal}}
	// From{{.DomainPascal}}AggregationToInfo maps domain entity to DB model for creation
	From{{.DomainPascal}}AggregationToInfo(item *entity.{{.DomainPascal}}) *model.{{.TablePascal}}
}

type {{.DomainPascal}}AclAdapter struct{}

func New{{.DomainPascal}}AclAdapter() I{{.DomainPascal}}AclAdapter {
	return &{{.DomainPascal}}AclAdapter{}
}

func (o *{{.DomainPascal}}AclAdapter) From{{.DomainPascal}}AggregationToInfo(item *entity.{{.DomainPascal}}) *model.{{.TablePascal}} {
	if item == nil {
		return nil
	}
	return &model.{{.TablePascal}}{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		OwnedBy:     item.OwnedBy,
		IsDeleted:   item.IsDeleted,
		CreatedBy:   item.CreatedBy,
		UpdatedBy:   item.UpdatedBy,
		DeletedBy:   item.DeletedBy,
	}
}

func (o *{{.DomainPascal}}AclAdapter) To{{.DomainPascal}}Aggregation(record *model.{{.TablePascal}}) *entity.{{.DomainPascal}} {
	if record == nil {
		return nil
	}
	return &entity.{{.DomainPascal}}{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		OwnedBy:     record.OwnedBy,
		IsDeleted:   record.IsDeleted,
		CreateTime:  {{.DomainPascal}}TimeValue(record.CreateTime),
		UpdateTime:  {{.DomainPascal}}TimeValue(record.UpdateTime),
		DeleteTime:  record.DeleteTime,
		CreatedBy:   record.CreatedBy,
		UpdatedBy:   record.UpdatedBy,
		DeletedBy:   record.DeletedBy,
	}
}

func (o *{{.DomainPascal}}AclAdapter) To{{.DomainPascal}}Aggregations(records []*model.{{.TablePascal}}) []*entity.{{.DomainPascal}} {
	items := make([]*entity.{{.DomainPascal}}, 0, len(records))
	for _, record := range records {
		items = append(items, o.To{{.DomainPascal}}Aggregation(record))
	}
	return items
}

func {{.DomainPascal}}TimeValue(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}
