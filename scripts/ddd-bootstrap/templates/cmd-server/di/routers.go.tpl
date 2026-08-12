package di

import (
	{{if .WithHealth}}
	"{{.ModulePath}}/app/ohs/remote/routers"
	{{end}}
	// ddd-bootstrap:import routers
	"go.uber.org/fx"
)

var RouterModule = fx.Options(
	fx.Invoke(
		{{if .WithHealth}}
		routers.Mount{{.DomainPascal}}RouteGroup,
		{{end}}
		// ddd-bootstrap:invoke routers
	),
)
