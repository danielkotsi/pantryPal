#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

load_env_file() {
  local env_file="$1"

  if [[ -f "$env_file" ]]; then
    set -a
    # shellcheck disable=SC1090
    source "$env_file"
    set +a
  fi
}

load_env_file "$ROOT_DIR/.env"
load_env_file "$ROOT_DIR/.env.local"
load_env_file "$ROOT_DIR/backend/.env"
load_env_file "$ROOT_DIR/backend/.env.local"

"$ROOT_DIR/scripts/db/local_db.sh" reset-all

cd "$ROOT_DIR/backend"
exec go run ./cmd/api
