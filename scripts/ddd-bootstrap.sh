#!/usr/bin/env bash
# Wrapper for the bundled ddd-bootstrap CLI.
#
# The CLI loads templates relative to its own working directory, so this
# wrapper always cd's into the module root before invoking it.
#
# Usage:
#   ddd-bootstrap.sh <subcommand> [args...]
#   ddd-bootstrap.sh init --project-name <name> --module <path> --output <dir> --db postgres
#   ddd-bootstrap.sh validate --project-root <dir>
#   ddd-bootstrap.sh add-domain --project-root <dir> --domain <d> --aggregate <P> --schema <s> --table <t>
#   ddd-bootstrap.sh add-entity --project-root <dir> --domain <d> --entity <e> --schema <s> --table <t>
#   ddd-bootstrap.sh remove-health --project-root <dir>
#
# The script auto-builds the CLI on first use and reuses the cached binary.

set -euo pipefail

SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_DIR="$SKILL_DIR/ddd-bootstrap"
BIN="$SKILL_DIR/.ddd-bootstrap-bin"

if [[ ! -d "$MODULE_DIR" ]]; then
  echo "error: ddd-bootstrap module not found at $MODULE_DIR" >&2
  exit 1
fi

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" || $# -eq 0 ]]; then
  cat <<USAGE
ddd-bootstrap.sh - scaffold a four-layer Go DDD project (bundled with go-dev-rules-ddd)

Subcommands (forwarded to the Go CLI):
  init           Create a new project skeleton with the health module.
  validate       Verify an existing project skeleton.
  add-domain     Add a full OHS -> Domain -> ACL -> SQL vertical domain.
  add-entity     Add only the persistent entity/repository/SQL/GORM Gen path.
  remove-health  Delete the bootstrap health module after verification.

Examples:
  ddd-bootstrap.sh init --project-name demo --module github.com/acme/demo \\
                       --output ./demo --db postgres
  ddd-bootstrap.sh validate --project-root ./demo
  ddd-bootstrap.sh add-domain --project-root ./demo --domain order \\
                       --aggregate Order --schema public --table orders

Each subcommand supports its own flags. Run the binary directly for full help:
  $BIN <subcommand> --help
USAGE
  exit 0
fi

if [[ ! -x "$BIN" ]] || [[ "${FORCE_REBUILD:-0}" == "1" ]]; then
  echo ">> building ddd-bootstrap (first use)..." >&2
  ( cd "$MODULE_DIR" && go build -o "$BIN" ./cmd/ddd-bootstrap )
fi

cd "$MODULE_DIR"
exec "$BIN" "$@"
