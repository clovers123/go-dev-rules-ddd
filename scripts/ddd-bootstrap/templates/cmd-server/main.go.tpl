package main

import (
	"go.uber.org/fx"

	"{{.ModulePath}}/cmd/server/di"
)

func main() {
	fx.New(
		di.InfrastructureModule,
		di.AclModule,
		di.DomainModule,
		di.OhsModule,
		di.RouterModule,
		di.InvokeModule,
	).Run()
}
