package {{.Domain}}

import "{{.ModulePath}}/app"

type Create{{.DomainPascal}}Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Create{{.DomainPascal}}Command struct {
	Name        string
	Description string
	User        app.Authorization
}

func (c Create{{.DomainPascal}}Command) GetOwner() string {
	return c.User.AccountID
}

func (c Create{{.DomainPascal}}Command) GetCreator() string {
	return c.User.Id
}
