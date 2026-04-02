# LabelClaw

LabelClaw is a full-stack dataset labeling application with a Go backend and a React/TypeScript frontend. The first iteration focuses on the configuration flow: defining task schemas, generating a React labeling panel with an LLM, previewing it, and saving the active configuration.

## Stack

- Go 1.24.x
- Node.js 22.x
- pnpm 10.x
- React 19 + Vite + TypeScript
- Playwright + Vitest

## Layout

- `backend/`: REST API, persistence, LLM integration, static asset serving
- `frontend/`: SPA, UI runtime for generated panels, unit/e2e tests
- `Makefile`: common entrypoint for install, lint, test, type-check, build, and dev tasks

## Prerequisites

LabelClaw expects a local Go installation. This workspace currently does not have `go` available in `PATH`, so backend commands and cross-stack `make` targets will fail until Go 1.24.x is installed locally.

## Quick Start

1. Install Go 1.24.x, Node.js 22.x, and pnpm 10.x.
2. Copy env files:
   - `cp backend/.env.example backend/.env`
   - `cp frontend/.env.example frontend/.env`
3. Run `make install`.
4. Start the frontend with `pnpm --dir frontend dev`.
5. Start the backend with `make dev`.

The Vite dev server runs on `http://127.0.0.1:5173` and talks to the Go API on `http://127.0.0.1:8080`.

## API

- `GET /api/health`
- `GET /api/configuration/current`
- `POST /api/configuration/generate`
- `PUT /api/configuration/current`

## Notes

- The backend persists a single active configuration at `backend/data/active-config.json` by default.
- The labelling tab is intentionally a placeholder in v1.
- Generated React code is treated as trusted application code in this iteration and is rendered inside the SPA.

