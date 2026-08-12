# Fx Dependency Injection Rules

Organize DI by architectural layer, matching both reference projects:

```text
cmd/server/di/
├── infrastructure.go  # config, DB, Redis, log, telemetry, runtime clients
├── acl.go             # repositories, ACL PL, clients, publishers, subscribers
├── domain.go          # domain services, factories, domain handlers
├── ohs.go             # appservices, OHS adapters, controllers, Fiber app
├── routers.go         # fx.Invoke route mounts only
├── realtime.go        # optional process-specific invocations
└── invoke.go          # lifecycle hooks and process startup
```

For other processes, keep the same layer-based naming under `cmd/{process}/di` and register only the modules that process needs.

## Registration Rules

- Register a new repository implementation and ACL mapper in `acl.go`.
- Register a domain service or factory in `domain.go`.
- Register an appservice, OHS adapter, controller, and protocol middleware in `ohs.go`.
- Register each route mount with `fx.Invoke` in `routers.go`.
- Register configuration, databases, caches, telemetry, message clients, and shared runtime providers in `infrastructure.go`.
- Put server/worker lifecycle startup and shutdown hooks in `invoke.go` or a clearly named process-specific file.

Do not create one DI file per bounded context. Do not write `fx.Provide` directly in `main.go`. Do not instantiate repositories or adapters inside services.

Use Fx groups for multiple implementations:

```go
fx.Annotate(NewSender, fx.ResultTags(`group:"channel_senders"`))
fx.Annotate(NewFactory, fx.ParamTags(`group:"channel_senders"`))
```

After adding any `NewXxx` provider, verify that its layer file includes it and that every route mount is invoked.
