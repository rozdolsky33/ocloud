# Contributing to ocloud

Thank you for your interest in contributing! This document explains how to work with the codebase and how to propose changes. It also documents the project architecture and data flow so contributions stay consistent with our separation of concerns.

If you’re new to the project, please skim README.md first to understand what ocloud does.


## Getting help, asking questions
- For questions and usage help, prefer Discussions: https://github.com/rozdolsky33/ocloud/discussions
- For bugs and feature requests, use GitHub Issues with the provided templates.


## Filing issues
We use GitHub Issue Forms tailored for this project:
- Bug report: include exact CLI commands, environment (OS/arch, region, tenancy/compartment context), and debug logs. Run with `--log-level debug` when possible.
- Feature request: describe the CLI UX, affected domain area (compute/network/storage/database/identity/configuration/search), and architecture impact across domain, mapping, services, oci adapters, and cmd layers. Provide acceptance criteria.

Before opening a new issue, search for existing issues and read this guide.


## Development prerequisites
- Go 1.22+ recommended
- Make (optional, for convenience targets)
- Access to an OCI tenancy if you plan to run live commands

Set up your environment:
- Ensure your OCI config and keys are set (typically in ~/.oci/config). You can override region with environment variable `OCI_REGION`.
- Useful environment variables consumed by the CLI: see internal/config/flags and the flags added in cmd/ and internal/app.


## Build, test, and run
- Build: `go build ./...`
- Lint/format: `gofmt -s -w .`
- Run locally: `go run . ...`
- Show version: `ocloud version`
- Tests: `go test ./...`

Some tests interact with mapping and services layers and can run without OCI credentials; anything hitting the live SDK should be behind interfaces and mocked in unit tests.


## Project architecture and separation of concerns
The codebase follows a clean layering that separates CLI, application services, domain models, mapping, and provider adapters. Please keep changes within the appropriate layer and depend only inward (upper layers depend on lower-level abstractions, not implementations).

Top-level directories you will most often touch:
- cmd/...: Cobra commands and flags. This is the CLI shell and presentation/wiring to services. No direct SDK calls here.
- internal/app: Application bootstrap and context (ApplicationContext) including OCI provider setup, identity client, tenancy/compartment resolution, logging, and shared IO writers.
- internal/services/...: Application use-cases. Services depend on domain interfaces (repositories) and orchestrate operations. They do not talk to the OCI SDK directly.
- internal/domain/...: Domain models and repository interfaces. No external SDK types should leak here.
- internal/mapping: Mappers between OCI SDK structs and internal domain models.
- internal/oci/...: Provider adapters and low-level clients for OCI SDK. Implements domain interfaces by calling the SDK and using mapping to produce domain models.
- internal/config, internal/logger, internal/printer, internal/tui: Supporting utilities, logging, printing/formatting, and optional TUI components.


### Data flow overview
1. User invokes `ocloud ...` in a terminal.
2. cmd/root.go initializes flags and either:
   - runs commands that do not require context, or
   - builds an ApplicationContext via internal/app (InitApp/InitializeAppContext) which sets up OCI provider, Identity client, region, tenancy/compartment context, and logger.
3. A Cobra command handler in cmd/... constructs or obtains a Service from internal/services/... and calls a use-case method.
4. The Service depends on domain repository interfaces (from internal/domain/...) and uses those to perform work. No direct SDK calls occur in the service.
5. A repository implementation in internal/oci/... (an adapter) calls the OCI Go SDK. It converts SDK responses to domain models via internal/mapping.
6. The Service returns domain models or view models to the cmd layer, which uses internal/printer or cmd/shared/display to render results (table/json, etc.).


### Wiring through adapters
- Services accept dependencies via constructors (interfaces). Example: compute/instance Service uses a compute.InstanceRepository.
- OCI adapters implement the repository interfaces by delegating to SDK clients. Example: internal/oci/storage/objectstorage/adapter.go wraps objectstorage.ObjectStorageClient and returns domain types using mapping functions.
- Mapping is centralized in internal/mapping (e.g., bucket_mapper.go, instance_mapper.go, etc.). Avoid ad-hoc conversions in services or cmd.
- ApplicationContext (internal/app) holds shared configuration, logger, and resolved tenancy/compartment context and is often used to build concrete adapters wired into services in the command layer.

This separation ensures:
- Domain and service logic is testable without hitting the OCI SDK.
- Mapping logic is reusable and validated via unit tests.
- CLI layer remains thin, focused on argument parsing, validation, and presentation.


## Making changes that span layers
When implementing features or fixes:
- Domain (internal/domain/...): Introduce/adjust models and repository interfaces only when necessary.
- Mapping (internal/mapping): Add/extend mappers to translate SDK structs to domain models and back if needed.
- OCI adapters (internal/oci/...): Implement or extend repository implementations. Do not return SDK types from public methods—convert to domain types.
- Services (internal/services/...): Orchestrate use cases using repository interfaces. Keep business/application logic here and make the behavior easy to unit test.
- cmd (...): Add Cobra commands/flags and wire services with appropriate adapters. Use internal/printer or cmd/shared/display for output.

Whenever you touch multiple layers, describe these changes in your PR (see PR checklist below).


## Coding style and guidelines
- Go modules are used; keep imports organized and avoid circular dependencies.
- Keep logs meaningful. Use internal/logger and structured logging where applicable. Respect the global log level.
- Never leak secrets to logs. For issue reports, ask users to redact before sharing.
- Follow existing patterns in the respective package for naming and error wrapping (`fmt.Errorf("context: %w", err)`).


## Tests
- Unit tests: Cover mapping and services with table-driven tests where possible.
- Adapters: Prefer testing via interfaces with mocked SDK clients; keep live integration tests opt-in if introduced.
- Keep tests fast and deterministic. Avoid real network calls in unit tests.


## Commit messages and PRs
- Write clear commit messages; reference issues when applicable (e.g., "Fixes #123").
- Open a PR with a concise description, screenshots (if applicable), and a checklist:
  - [ ] Changes align with the architecture (domain, mapping, services, cmd)
  - [ ] New/updated unit tests included
  - [ ] Docs updated (README/COMMAND help if needed)
  - [ ] Backward compatibility considered


## Local troubleshooting tips
- Use `--log-level debug` to get detailed logs.
- Check the resolved tenancy/compartment in output (root cmd prints OCI configuration when no subcommands are given).
- Use environment variables like `OCI_REGION`, and the flags under internal/config/flags; ApplicationContext resolves tenancy/compartment through flags, env, and mapping files (`~/.oci/tenancy-map.yaml`).


## Code ownership and review
- Keep PRs small and focused.
- Expect CI to run `go test ./...` at minimum.
- Reviewers will check architectural alignment and separation of concerns first.

Thank you for contributing!