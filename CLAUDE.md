# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

The main project guide lives in AGENTS.md (structure, CLI verbs, name resolution, editor workflow, API client internals, gotchas). It is imported here so there is one copy to maintain:

@AGENTS.md

## Additions and corrections to AGENTS.md

### Commands

```bash
make build                                   # stripped binary at ./netbird-cli
go vet ./...                                 # the only static check; no linter config, no tests
go run ./cmd/cli get groups                  # run without building (needs config/env, see below)
NETBIRD_CLI_URL=... NETBIRD_CLI_TOKEN=... ./netbird-cli get peers --dry-run
```

CI: `.github/workflows/release.yml` runs only on tag push, cross-compiles linux/darwin/windows × amd64/arm64 with `CGO_ENABLED=0`, and publishes a GitHub release with checksums. There is no test or lint job.

### Behavior worth knowing before editing

- `main.go` loads the config and builds the global client `c` **before** cobra runs. Every subcommand (including `completion`, `generate` and `--help`) exits early when the URL or token is missing, so set `NETBIRD_CLI_URL`/`NETBIRD_CLI_TOKEN` (any non-empty values work for commands that make no API calls).
- `exitErr(msg, err)` in `root.go` only prints to stderr. It does **not** exit, so callers must `return` right after it. Follow that pattern.
- `apply` handles only `groups`, `users`, `peers`, `policies`, `setupkeys`, `accounts` and the Agent Network lists `agentproviders`, `guardrails`, `agentpolicies`, `budgetrules` (the `spec` struct in `apply.go`). Groups are applied first, and agent providers/guardrails before agent policies, so later resources can reference earlier ones by name. It matches existing resources by name and updates them, otherwise it creates them.
- Table column widths are computed on `stripANSI` output (`cmd/cli/color.go`), so colored cells keep columns aligned.
- When a table needs different columns than the API struct gives (logs, catalog, models), `agentnetwork.go` maps to a small local row struct whose JSON tags are whitelisted; follow that instead of adding more fields to the global whitelist.
- Flag variables in `flags.go` are shared, and cobra writes a flag's default into the variable when the flag is registered, so the last registration's default wins for every command. If two commands need different defaults for the same variable, use separate variables (see `portalReturnFlag` vs `baseURLFlag` in `msp.go`), or only read the variable when `cmd.Flags().Changed(...)`.
- GET responses mask secrets (`****`). Edit flows that PUT the fetched object must drop them first (`editTransformed` + `stripMasked` in `integrations.go`), otherwise the mask overwrites the real secret. Some GET shapes also differ from the PUT shape (EDR `groups` are objects on GET, IDs on PUT; ingress ports return `port_range_mappings` but take `port_ranges`); `editTransformed` is where to reshape them.
- A new command whose endpoints are `x-cloud-only` in the spec (at operation or tag level) must be added to the `requireEdition(editionCloud, ...)` list in `edition.go`, otherwise self-hosted users see it and get a raw 404. Self-hosted-only features (embedded IdP, instance setup) go in the `editionSelfHosted` list.
- `createWithEditor` returns `nil, nil` when the user saves without changes; every caller must `return` on a nil result before POSTing.
- The personal-access-token API is in `internal/client/tokens.go`, but its commands (`get tokens <user-id>`, `create token`, `delete token <user-id> <token-id>`) are defined in `cmd/cli/accounts.go`. There is no `cmd/cli/tokens.go`.
- `dryRunCheck(payload)` prints the JSON body and `dryRunMsg(msg)` prints a one-line message. Both return true when `--dry-run` is set, so the caller can return before calling the API.
- `config.Config` has `json:` tags but is decoded with `yaml.v3`, which ignores them. It works only because yaml.v3 lowercases field names by default. Add `yaml:` tags if you add a field whose name is not just the lowercase field name.
