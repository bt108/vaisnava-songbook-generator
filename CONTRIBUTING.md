# Contributing

## Prerequisites

See [README](README.md#prerequisites): Go toolchain, LuaLaTeX with the
`Charis SIL` font, `Noto Serif Bengali` font, and `yamlfmt`.

## Content edits

Everything is YAML-driven — edit `data/`, not Go code, to change content.
Follow the YAML style rules in `AGENTS.md`, then normalize:

```sh
yamlfmt .
```

## Checks

```sh
make test   # go test ./...
make lint   # gofmt check + go vet
```

Both must pass before opening a PR.

## Branch / PR flow

1. Create a short-lived branch from `main` (`feat/...`, `fix/...`, `docs/...`).
2. Keep commits small with conventional messages (`feat:`, `fix:`, `docs:`, `chore:`, `ci:`).
3. Open a PR against `main` describing what changed and why.
