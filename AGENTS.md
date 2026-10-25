# Repository Guidelines

## Project Structure & Module Organization

Killgrave is a Go HTTP mock server. `cmd/killgrave/main.go` is the executable entry point; `internal/app/` wires the application together, and `internal/app/cmd/` defines the Cobra CLI. Configuration and file watching live in `internal/`; HTTP routing, imposters, proxying, and response handling live in `internal/server/http/`.

Unit tests sit beside source files as `*_test.go`. Fixtures live in `internal/test/testdata/` and `internal/server/http/test/testdata/`. Process-level acceptance tests and `.txtar` fixtures live in `acceptance/`. Consult `README.md` for configuration and imposter examples.

## Build, Test, and Development Commands

Use Go 1.27, matching `go.mod` and CI.

- `make build`: produce `bin/killgrave` with branch and commit version metadata.
- `./bin/killgrave --imposters ./internal/server/http/test/testdata/imposters`: run locally with existing fixtures on localhost:3000.
- `go test ./...`: run the default unit test suite.
- `go test -v -vet=off -race ./...`: reproduce the CI unit test command with race detection.
- `make acceptance`: build the executable, then run acceptance tests with the `acceptance` build tag.
- `gofmt -w path/to/changed.go`: format changed Go files; use `go vet ./...` for static checks.

## Coding Style & Naming Conventions

Follow existing Go conventions: use `gofmt` formatting and tab indentation, short lowercase package names, exported `CamelCase` identifiers, and unexported `camelCase` identifiers. Keep related behavior in the existing packages. Group standard-library, external, and repository imports consistently with nearby code. Name imposter fixtures with `.imp.json`, `.imp.yml`, or `.imp.yaml` extensions.

## Testing Guidelines

Use Go's `testing` package with Testify assertions. Name tests `Test<Behavior>` and use descriptive table-driven subtests for input variations. Add regression coverage for behavior changes, including error cases. Acceptance scenarios belong under `acceptance/tests/<scenario>/`, with `config.txtar` and HTTP request/response archives in `http/`. Run unit and acceptance suites before submitting. No numeric coverage threshold is configured.

## Commit & Pull Request Guidelines

Recent commits use concise action-oriented subjects, such as `Fix imposters search` and `Replace spf13/afero with io/fs`; merged subjects often include a PR number. Follow that style.

Follow `CONTRIBUTING.md`: ask to take an issue, confirm its milestone, and base work on the milestone branch. PR descriptions should explain the change, link the issue, and list validation performed. For bugs, include reproduction steps, expected versus actual behavior, and relevant configuration or imposters.
