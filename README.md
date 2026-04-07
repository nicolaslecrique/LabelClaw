# LabelClaw

LabelClaw is a full-stack dataset labeling application with a Go backend and a React/TypeScript frontend. The first iteration focuses on the configuration flow: defining task schemas, generating a React labeling panel with an LLM, previewing it, and saving the active configuration.

## Stack

- Go 1.26.x
- Node.js 22.x
- pnpm 10.x
- React 19 + Vite + TypeScript
- Playwright + Vitest

## Layout

- `docs/spec.md`: product behavior and feature scope
- `backend/`: REST API, persistence, LLM integration, static asset serving
- `frontend/`: SPA, UI runtime for generated panels, unit/e2e tests
- `Makefile`: common entrypoint for install, lint, test, type-check, build, and dev tasks

## Quick Start

1. Install Go 1.24.x, Node.js 22.x, and pnpm 10.x.
2. Run `make install`.
3. Start the frontend with `pnpm --dir frontend dev`.
4. Start the backend with `make dev`.

The Vite dev server runs on `http://127.0.0.1:5173` and talks to the Go API on `http://127.0.0.1:8080`.

## API

- `GET /api/health`
- `GET /api/configuration/current`
- `POST /api/configuration/generate`
- `PUT /api/configuration/current`

## Notes

- The backend persists a single active configuration at `backend/data/active-config.json` by default.
