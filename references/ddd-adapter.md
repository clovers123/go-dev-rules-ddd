# Adapter and Object Assembly Rules

## OHS Adapter

Place OHS adapters in `app/ohs/pl/{bc}/adapter`.

Use them for:

- Fiber-bound DTO to Command/Query conversion;
- extracting authorization and tenant metadata from the protocol context;
- Command/Envelope to Domain Entity conversion;
- Domain Entity to BO/response conversion.

Do not access PostgreSQL, Redis, Kafka, generated models, or repository implementations from an OHS adapter.

Keep DTO and Query types separate even when they share fields. DTOs carry protocol tags; Queries and Commands are protocol-neutral use-case inputs. Convert explicitly in the adapter.

## ACL PL Adapter

Place ACL mapping adapters in `app/acl/pl/{bc}`.

Use them only for:

- database/external model to Domain Entity conversion;
- Domain Entity to database/external model conversion;
- batch conversion, audit-field mapping, null handling, and JSON serialization.

ACL PL must be pure mapping code with no database, cache, message, network, or business-rule execution.

## Assembly Boundary

Do not spread field-by-field assembly across controllers, appservices, domain services, or repositories. If conversion assigns more than three meaningful fields, move it to the matching adapter unless it is an entity constructor enforcing a domain invariant.

Use these naming patterns:

- OHS input: `FromCreateXxxDTO`, `FromXxxCommand`, `FromEnvelope`.
- OHS output: `FromXxxEntity`, `FromXxxEntities`, `ToXxxBO`.
- ACL input/output: `ToEntity`, `ToEntities`, `ToModel`, or semantic equivalents used consistently by the repository.
- Interfaces: `I{Context}Adapter` and `I{Context}AclAdapter` when the repository uses the existing `I` convention.

Application services live in `app/ohs/local/appservice/{bc}` and may call both OHS adapters and Domain services. Domain services must receive domain objects or semantic values, never Fiber contexts, DTOs, Commands, BOs, or persistence models.
