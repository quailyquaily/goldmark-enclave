# Repository Guidelines

**NOTICE: THIS IS A GOLANG REPOSITORY, USE TAB FOR INDENTATION.**

## Project Structure & Module Organization
- Root package `enclave` wires the Goldmark extension: `enclave.go`, `transformer.go`, `render.go`.
- Core types/config: `core/` (`core/core.go`). Helper utilities: `helper/`.
- Feature packages: `callout/`, `fence/`, `href/`, `mark/`, `kbd/`, and `object/` (YouTube, Bilibili, X/Twitter, TradingView, Dify, Quail, Spotify, audio, image helpers).
- Tests live alongside packages as `*_test.go` (e.g., `callout/callout_test.go`).

## Build, Test, and Development Commands
- `go build ./...` — compile all packages to ensure the tree builds.
- `go test ./...` — run unit tests; add `-v` for verbose, `-run <Regex>` to filter, `-short` to skip any long/integration tests.
- `go test ./... -cover` — run tests with coverage.
- `go vet ./...` — static checks for common issues.
- `go fmt ./...` — format code; CI/code review expects `gofmt`-clean diffs.

Go version: module targets Go 1.22+ (toolchain 1.24). Use a recent Go toolchain.

## Coding Style & Naming Conventions
- Follow standard Go style: tabs for indentation; `gofmt` and `goimports` clean.
- Exported identifiers use PascalCase with doc comments; unexported use lowerCamelCase.
- Package names are short, lower-case nouns (match existing layout).
- Errors: return `error`, avoid `panic` in library code. Prefer small, focused funcs.

## Testing Guidelines
- Use table-driven tests in `*_test.go` in the same package.
- Prefer deterministic, offline tests. Avoid network I/O (e.g., oEmbed calls) or guard with `testing.Short()`.
- Validate rendered HTML strings precisely where feasible. Example: `go test ./object -run TestFormalizeImageSize`.

## Commit & Pull Request Guidelines
- Commits are short, imperative, and often prefixed (observed: `feat:`, `docs:`, `Refactor`). Prefer Conventional Commits where possible: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`.
- PRs should include:
  - Clear description and rationale, linked issues.
  - Tests for new behavior and updates to existing tests.
  - Notes/screenshots when HTML output changes.
  - No unrelated refactors.

## Security & Configuration Tips
- `core.Config` controls behavior (e.g., `IframeDisabled`); use fallbacks in constrained environments.
- Sanitize/escape when introducing new HTML paths; keep parity with existing renderers.
- Do not introduce network calls in hot paths without timeouts and clear test strategy.
