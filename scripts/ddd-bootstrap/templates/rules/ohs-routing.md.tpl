# OHS Routing and Controller Rules

## Directory Layout

```text
app/ohs/
├── local/appservice/{bc}/
├── pl/{bc}/
│   └── adapter/
└── remote/
    ├── controller/{bc}/
    ├── routers/{bc}.go
    └── middleware/
```

## Router Rules

- Put one context-oriented router file in `app/ohs/remote/routers`.
- Accept the preconfigured `fiber.Router` when registering authenticated `/api` routes; do not duplicate `/api` in route strings.
- Use `/v1/...` paths inside the router mount unless the repository defines another versioning convention.
- Put public, SDK, webhook, or WebSocket groups on explicitly injected routers/apps with their own middleware.
- Register mount functions through `fx.Invoke` in `cmd/server/di/routers.go`.
- Keep only path, method, middleware, and route-name declarations in router files.

## Controller Rules

A controller performs only this sequence:

1. Bind URI, query, header, or body fields into a DTO.
2. Call an OHS adapter to create a Command or Query and extract user/tenant metadata.
3. Call `app/ohs/local/appservice/{bc}`.
4. Wrap the result with the repository's standard response type.

Controllers must not call repositories or Domain services directly, open transactions, assemble entities, or contain product decisions.

## Appservice Rules

Application services belong to OHS local. They may own transactions, workflows, idempotency, and coordination across Domain services and ACL ports. They must not receive Fiber contexts or return persistence models.
