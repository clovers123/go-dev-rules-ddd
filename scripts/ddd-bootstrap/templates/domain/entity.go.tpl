package entity

import "time"

type {{.DomainPascal}} struct {
	ID          string
	Name        string
	Description string
	OwnedBy     string
	IsDeleted   bool
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  *time.Time
	CreatedBy   *string
	UpdatedBy   *string
	DeletedBy   *string
}

func New{{.DomainPascal}}(name, description, ownedBy string) *{{.DomainPascal}} {
	return &{{.DomainPascal}}{
		Name:        name,
		Description: description,
		OwnedBy:     ownedBy,
		IsDeleted:   false,
	}
}
