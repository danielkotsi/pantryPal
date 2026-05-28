# Backend DB Setup

Local SQLite setup for backend Phase 1 tasks.

Food source note:
- App pantry and recipe ingredients now reference USDA foods directly via `fdc_id`.
- The old local `ingredients` table is no longer part of the schema.

## One-command full reset

Run from repo root:

```bash
./scripts/db/reset-local-db
```

This will:
- Recreate `database/sqlite/pantrypal.db`
- Apply `backend/migrations/001_init_schema.sql`
- Apply `backend/migrations/010_usda_schema.sql`
- Seed demo data from `backend/seeds/001_seed_demo.sql`
- Import USDA dataset from `dataset/FoodData_Central_foundation_food_json_2026-04-30.json`

## Optional commands

Create DB only if missing:

```bash
./scripts/db/local_db.sh setup
```

Reset app tables only (preserve USDA tables and imported data):

```bash
./scripts/db/local_db.sh reset-app
```

Force recreate DB (app + USDA):

```bash
./scripts/db/local_db.sh reset-all
```

Alias kept for compatibility:

```bash
./scripts/db/local_db.sh reset
```

Override DB file path:

```bash
DB_PATH=/tmp/pantrypal-demo.db ./scripts/db/local_db.sh reset
```

Disable USDA import on setup/reset-all:

```bash
IMPORT_USDA=0 ./scripts/db/local_db.sh reset-all
```

Use a different USDA dataset file path:

```bash
USDA_DATASET_PATH=/tmp/fdc.json ./scripts/db/local_db.sh reset-all
```
