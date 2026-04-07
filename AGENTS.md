# LabelClaw Agent Guide

These instructions apply to the whole repository.

## Source Of Truth

- Use `README.md` for stack, layout, and local development commands.
- Use `docs/spec.md` for product behavior and feature scope.

## Coding Style

- Prefer the most standard, well-supported solution for the current ecosystem and year.
- Keep implementations simple: follow KISS and YAGNI.
- Preserve separation of concerns, strong cohesion, and low coupling.
- Prefer pure functions, explicit data flow, immutability, and precise types when practical.
- Use strict, standard formatters, linters, and type checkers already established in the repo.
- Add or update automated tests for behavior changes; include end-to-end coverage for user-facing specifications when practical.

## Architecture Boundaries

- Backend owns business logic, LLM calls, configuration persistence, and serving frontend assets.
- Frontend owns UI logic and should interact with the backend only through the REST API.
- Keep the labelling tab empty unless the task explicitly expands its scope.
- Generated React code is treated as trusted application code in this iteration.
