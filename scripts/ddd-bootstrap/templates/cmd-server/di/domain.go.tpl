package di

import (
	"go.uber.org/fx"

	// ddd-bootstrap:import domain
)

var DomainModule = fx.Options(
	fx.Provide(
		// ddd-bootstrap:provide domain
	),
)
