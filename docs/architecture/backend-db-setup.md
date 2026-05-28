# Backend DB Setup

Local SQLite setup for backend Phase 1 tasks.

## One-command reset

Run from repo root:

```bash
./scripts/db/reset-local-db
```

This will:
- Recreate `database/sqlite/pantrypal.db`
- Apply `backend/migrations/001_init_schema.sql`
- Seed demo data from `backend/seeds/001_seed_demo.sql`

## Optional commands

Create DB only if missing:

```bash
./scripts/db/local_db.sh setup
```

Force recreate DB:

```bash
./scripts/db/local_db.sh reset
```

Override DB file path:

```bash
DB_PATH=/tmp/pantrypal-demo.db ./scripts/db/local_db.sh reset
```
