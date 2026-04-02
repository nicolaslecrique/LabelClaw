#!/usr/bin/env sh
set -eu

cleanup() {
  if [ -n "${BACKEND_PID:-}" ]; then
    kill "${BACKEND_PID}" 2>/dev/null || true
  fi

  if [ -n "${FRONTEND_PID:-}" ]; then
    kill "${FRONTEND_PID}" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

(cd backend && go run ./cmd/server) &
BACKEND_PID=$!

(cd frontend && pnpm dev --host 127.0.0.1) &
FRONTEND_PID=$!

wait "${BACKEND_PID}" "${FRONTEND_PID}"
