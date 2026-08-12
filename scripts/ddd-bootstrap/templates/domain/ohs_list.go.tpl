package {{.Domain}}

import "{{.ModulePath}}/app"

type List{{.DomainPascalPlural}}Request struct {
	Keyword string `json:"keyword" query:"keyword"`
	Limit   int    `json:"limit"  query:"limit"`
	Offset  int    `json:"offset" query:"offset"`
}

type List{{.DomainPascalPlural}}Query struct {
	Keyword string
	Limit   int
	Offset  int
	User    app.Authorization
}

func (q List{{.DomainPascalPlural}}Query) GetOwner() string {
	return q.User.AccountID
}
