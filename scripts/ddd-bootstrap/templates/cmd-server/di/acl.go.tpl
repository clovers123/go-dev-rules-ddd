package di

import (
	"go.uber.org/fx"

	// ddd-bootstrap:import acl
)

var AclModule = fx.Options(
	fx.Provide(
		// ddd-bootstrap:provide acl
	),
)
