#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DB_PATH="${DB_PATH:-$ROOT_DIR/database/sqlite/pantrypal.db}"
APP_MIGRATION_PATH="$ROOT_DIR/backend/migrations/001_init_schema.sql"
USDA_SCHEMA_MIGRATION_PATH="$ROOT_DIR/backend/migrations/010_usda_schema.sql"
SEED_PATH="$ROOT_DIR/backend/seeds/001_seed_demo.sql"
USDA_IMPORT_PATH="$ROOT_DIR/scripts/db/import_usda.py"
USDA_DATASET_PATH="${USDA_DATASET_PATH:-$ROOT_DIR/dataset/FoodData_Central_foundation_food_json_2026-04-30.json}"
IMPORT_USDA="${IMPORT_USDA:-1}"
MODE="${1:-reset-app}"

if [[ ! -f "$APP_MIGRATION_PATH" ]]; then
  printf "Missing migration file: %s\n" "$APP_MIGRATION_PATH" >&2
  exit 1
fi

if [[ ! -f "$USDA_SCHEMA_MIGRATION_PATH" ]]; then
  printf "Missing migration file: %s\n" "$USDA_SCHEMA_MIGRATION_PATH" >&2
  exit 1
fi

if [[ ! -f "$SEED_PATH" ]]; then
  printf "Missing seed file: %s\n" "$SEED_PATH" >&2
  exit 1
fi

mkdir -p "$(dirname "$DB_PATH")"

apply_migrations() {
  (
    cd "$ROOT_DIR/backend"
    go run ./cmd/dbseed -db "$DB_PATH" -migration "$APP_MIGRATION_PATH"
  )

  sqlite3 "$DB_PATH" < "$USDA_SCHEMA_MIGRATION_PATH"
}

seed_demo() {
  sqlite3 "$DB_PATH" < "$SEED_PATH"
}

usda_food_count() {
  if [[ ! -f "$DB_PATH" ]]; then
    printf "0"
    return
  fi

  sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM usda_foods;" 2>/dev/null || printf "0"
}

import_usda_if_configured() {
  if [[ "$IMPORT_USDA" != "1" ]]; then
    return
  fi

  if [[ ! -f "$USDA_IMPORT_PATH" ]]; then
    printf "USDA import script missing, skipping: %s\n" "$USDA_IMPORT_PATH"
    return
  fi

  if [[ ! -f "$USDA_DATASET_PATH" ]]; then
    printf "USDA dataset missing, skipping: %s\n" "$USDA_DATASET_PATH"
    return
  fi

  python3 "$USDA_IMPORT_PATH" --db-path "$DB_PATH" --dataset-path "$USDA_DATASET_PATH"
}

reset_app_tables() {
  sqlite3 "$DB_PATH" <<'SQL'
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS consumption_log;
DROP TABLE IF EXISTS plan_meals;
DROP TABLE IF EXISTS meal_plans;
DROP TABLE IF EXISTS pantry_items;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS user_body_metrics;
DROP TABLE IF EXISTS users;

PRAGMA foreign_keys = ON;
SQL
}

case "$MODE" in
  setup)
    if [[ -f "$DB_PATH" ]]; then
      printf "Database already exists: %s\n" "$DB_PATH"
      printf "Use '%s reset-all' to recreate it.\n" "$0"
      exit 0
    fi
    apply_migrations
    import_usda_if_configured
    seed_demo
    ;;
  reset-app)
    if [[ -f "$DB_PATH" ]]; then
      reset_app_tables
    fi
    apply_migrations
    if [[ "$(usda_food_count)" == "0" ]]; then
      import_usda_if_configured
    fi
    seed_demo
    ;;
  reset-all|reset)
    rm -f "$DB_PATH"
    apply_migrations
    import_usda_if_configured
    seed_demo
    ;;
  *)
    printf "Usage: %s [setup|reset-app|reset-all|reset]\n" "$0" >&2
    exit 1
    ;;
esac

printf "Database ready at %s\n" "$DB_PATH"
