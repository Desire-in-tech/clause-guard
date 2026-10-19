#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ -f "$ROOT_DIR/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT_DIR/.env"
  set +a
fi

cleanup() {
  local status=$?
  trap - EXIT INT TERM
  kill "${backend_pid:-}" "${frontend_pid:-}" 2>/dev/null || true
  wait "${backend_pid:-}" 2>/dev/null || true
  wait "${frontend_pid:-}" 2>/dev/null || true
  exit "$status"
}

trap cleanup EXIT INT TERM

(
  cd "$ROOT_DIR/backend"
  exec go run ./cmd/api
) &
backend_pid=$!

(
  cd "$ROOT_DIR/frontend"
  exec npm run dev -- --host 0.0.0.0
) &
frontend_pid=$!

wait -n "$backend_pid" "$frontend_pid"
