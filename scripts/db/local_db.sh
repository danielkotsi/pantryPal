#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DB_PATH="${DB_PATH:-$ROOT_DIR/database/sqlite/pantrypal.db}"
MIGRATION_PATH="$ROOT_DIR/backend/migrations/001_init_schema.sql"
SEED_PATH="$ROOT_DIR/backend/seeds/001_seed_demo.sql"
MODE="${1:-reset}"

if [[ ! -f "$MIGRATION_PATH" ]]; then
  printf "Missing migration file: %s\n" "$MIGRATION_PATH" >&2
  exit 1
fi

if [[ ! -f "$SEED_PATH" ]]; then
  printf "Missing seed file: %s\n" "$SEED_PATH" >&2
  exit 1
fi

mkdir -p "$(dirname "$DB_PATH")"

case "$MODE" in
  setup)
    if [[ -f "$DB_PATH" ]]; then
      printf "Database already exists: %s\n" "$DB_PATH"
      printf "Use '%s reset' to recreate it.\n" "$0"
      exit 0
    fi
    ;;
  reset)
    ;;
  *)
    printf "Usage: %s [setup|reset]\n" "$0" >&2
    exit 1
    ;;
esac

if [[ "$MODE" == "reset" ]]; then
  (
    cd "$ROOT_DIR/backend"
    go run ./cmd/dbseed -db "$DB_PATH" -migration "$MIGRATION_PATH" -seed "$SEED_PATH" -reset
  )
else
  (
    cd "$ROOT_DIR/backend"
    go run ./cmd/dbseed -db "$DB_PATH" -migration "$MIGRATION_PATH" -seed "$SEED_PATH"
  )
fi

printf "Database ready at %s\n" "$DB_PATH"
