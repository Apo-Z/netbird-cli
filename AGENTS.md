# netbird-cli

Go CLI to interact with the Netbird API.

## Quick start

```bash
make build          # → ./netbird-cli
make install        # → /usr/local/bin/netbird-cli
```

There are no tests. CI (`.github/workflows/release.yml`) only builds and publishes release binaries on tag push.

## Config

YAML config at `~/.config/netbird-cli/config.yaml`, overridable via `NETBIRD_CLI_CONFIG_FILE`.
API URL, token and edition can be overridden with `NETBIRD_CLI_URL`, `NETBIRD_CLI_TOKEN` and `NETBIRD_CLI_EDITION` env vars. See `config/config.go:10-14`.
The URL is stripped of trailing `/api` and `/` automatically.

## Structure

```
cmd/cli/main.go       — entrypoint; constructs the global `*client.Client`
cmd/cli/root.go       — cobra root command (`netbird`) + top-level verbs (get/create/edit/delete) + global flags
cmd/cli/flags.go      — shared flag variables used across all resource commands
cmd/cli/completion.go — autocompletion helpers that fetch names from the API
cmd/cli/format.go     — output formatting (table / key-value / JSON / YAML)
cmd/cli/editor.go     — interactive YAML editor for `edit` and `create --edit`
cmd/cli/apply.go      — declarative config (`apply -f <file.yaml>`)
cmd/cli/generate.go   — YAML template generation (`generate user`, `generate group`, ...)
cmd/cli/describe.go   — detailed resource view with cross-referencing
cmd/cli/color.go      — ANSI color helpers (colorBool, colorStatus, colorState, colorDecision)
cmd/cli/integrations.go — event streaming, notification channels, EDR (+ bypass/unbypass verbs), IdP sync (+ sync/regenerate verbs); shared helpers editTransformed, stripMasked, maskSecrets
cmd/cli/reverseproxy.go — reverse proxy domains (+ validate verb), proxy tokens, proxy cluster delete, ingress peers and port allocations
cmd/cli/msp.go        — MSP tenants (get/create/edit + `msp` subcommands) and `billing` subcommands
cmd/cli/agentnetwork.go — Agent Network (AI agent gateway): providers, agent policies, guardrails, budget rules, settings, gateway, catalog, model discovery, usage, consumption, `events agent`, describe; apply support lives in apply.go
cmd/cli/*.go          — one file per resource (users.go, groups.go, peers.go, policies.go, networks.go, setupkeys.go, posturechecks.go, routes.go, dns.go, dnszones.go, accounts.go, events.go, services.go, geolocations.go, jobs.go, idp.go, instance.go)
internal/client/      — Netbird API client (Token auth, header `Authorization: Token <token>`)
  client.go           — core `Client` struct with doRequest/doGet/doPost/doPut/doDelete helpers, GetRaw/PutRaw/PostRaw, in-memory name→ID cache, IDToName reverse lookup
  *.go                — one file per API resource
config/config.go      — YAML config loader
```

## CLI verbs

Kubectl-style: `netbird <verb> <resource> [name] [flags]`

```
netbird get <resource> [name]        # list or display
netbird describe <resource> <name>   # rich detail, all IDs resolved to names
netbird create <resource> [flags]    # create (--edit opens the editor)
netbird edit <resource> <name>       # edit (interactive editor, with or without flags)
netbird delete <resource> <name>     # delete
netbird apply -f <file.yaml>         # declarative config
netbird generate <resource|all>      # YAML template
```

Special verbs: `info`, `approve`, `block`, `unblock`, `invite`, `whoami`, `setup`, `events`, `bypass`, `unbypass`, `sync`, `validate`, `regenerate`, `msp`, `billing`.

When adding a new resource, register it under the appropriate verb (`getCmd.AddCommand`, `createCmd.AddCommand`, etc.) in the resource's `init()`.

## Name resolution (name as ID)

The client exposes `Resolve*ID(nameOrID string) (string, error)` methods for groups, users, peers, policies, and networks. These accept either a name or an ID and return the real API ID. Resolution uses an in-memory cache.

Key resolution methods:
- `ResolveGroupID(nameOrID)` — looks up by name first, falls back to ID
- `ResolveUserID(nameOrEmailOrID)` — looks up by email, then name, then ID
- `ResolvePeerID(nameOrID)` — looks up by name, falls back to ID
- `ResolvePolicyID(nameOrID)`, `ResolveNetworkID(nameOrID)`
- Agent Network: `ResolveAgentProviderID`, `ResolveAgentPolicyID`, `ResolveAgentGuardrailID`, `ResolveAgentBudgetRuleID` (one list call via `resolveByList`, which also fills the `IDToName` cache)

Also available: `GetGroupByName`, `GetUserByEmail`, `GetPeerByName`, `GetPolicyByName`, `GetNetworkByName`.

`c.IDToName(resourceType, id) string` — reverse lookup (ID → name), batch-fetches once per resource type.

## Editor workflow

`edit` commands:
1. Resolve resource name → ID
2. GET the resource (via `GetRaw` or typed method)
3. Marshal to YAML with ID annotations (`# name`)
4. Open `$EDITOR` (default `vim`) with the annotated YAML
5. Strip comments, resolve names back to IDs
6. PUT the full object
7. If file unchanged → skip PUT, print "Edit cancelled"

`create --edit` commands:
1. Generate a template YAML with defaults
2. Open `$EDITOR` with the template
3. Strip comments, resolve names to IDs
4. POST the result
5. If file unchanged → skip POST, print "Creation cancelled"

Resources with parent dependencies (networkresource, dnsrecord, token) include the parent ID in the YAML template (`network_id`, `zone_id`, `user_id`) rather than requiring a separate flag.

## Flags

- `-o / --output json|yaml` — output format (default: human-readable table/key-value)
- `--dry-run` — show what would be done without executing (print request body, skip API calls)
- `--edit` on `create` — open editor with template instead of using flags

Flag variables are shared across command files via `flags.go`. When adding new flags, add the variable there.

## Autocompletion

`completion.go` provides `validArgsFunc(fetch func() ([]string, error))` that wraps API lookups for dynamic shell completion.
Each `get`/`edit`/`delete` command with a name parameter has `ValidArgsFunction: validArgsFunc(xxxNames)`.

## API client internals

- `bodyToStructure[T]` — unmarshals response body into `*T`
- `bodyToSlice[E]` — unmarshals response body into `[]E` (element type, not slice type)
- `doRequest(method, path, body)` — core request builder, sets Accept/Authorization/Content-Type headers, checks HTTP status
- `doGet(path, params)` — GET with optional query params map
- `doPost(path, payload)`, `doPut(path, payload)` — JSON POST/PUT
- `doDelete(path)` — DELETE
- `GetRaw(path)` — GET returning raw bytes
- `PutRaw(path, body)` — PUT returning raw bytes
- `PostRaw(path, body)` — POST returning raw bytes
- `IDToName(resource, id)` — reverse name lookup with lazy batch-fetch

The URL is passed in at construction (`NewNetbirdClient(url, token)`). All API paths are relative (e.g., `/api/users`). Trailing `/` and `/api` are stripped from the URL.

## Gotchas

- all user-facing error strings are in English
- `types.go` under `internal/client/` is empty (superseded by per-resource files)
- `utils/` directory exists but contains no files
- flag variables are shared across command files via `flags.go` — when adding new flags, add the var there
- Go 1.26.2, dependencies: cobra + yaml.v3 only (no external HTTP deps beyond stdlib)
- Netbird API uses PUT (full replace), not PATCH: always send the complete object when editing
- `GET /api/users/{id}` does NOT exist in the Netbird API — user lookup is done via listing all users
- `GET /api/accounts/{id}` does NOT exist either — account data is fetched via listing
- `create --edit` uses `createCmd.PersistentFlags()` for the `--edit` flag, so it's inherited by all create subcommands
- edition support lives in `cmd/cli/edition.go`: `requireEdition(editionCloud|editionSelfHosted, cmds...)` tags commands; the other edition hides them (`applyEdition`, run from `main.go`) and refuses them (`checkEdition`, root `PersistentPreRun`). Edition = config `edition` / `NETBIRD_CLI_EDITION`, else detected from the URL (`*.netbird.io` = cloud). Tag every new command whose endpoints are `x-cloud-only` in the OpenAPI spec (operation or tag level)
- print API errors with `printErr(err)` or `exitErr(msg, err)` (stderr, adds a hint on self-hosted 404s), not `fmt.Printf`
- the table field whitelist in `format.go:getTableFields` must be updated when adding new resource types with new JSON fields, otherwise only `id` will appear in table output
- describe commands now respect `-o json|yaml`; in default mode, rich human-readable output is used
- the Agent Network OpenAPI spec lives upstream at `shared/management/http/api/openapi.yml` in netbirdio/netbird (tag `Agent Network`); ID-list fields (`source_groups`, `destination_provider_ids`, `guardrail_ids`, `target_groups`, `target_users`) are registered in `editor.go:knownIDFields` and `format.go:idListFields` so they are annotated/resolved by name
- `internal/client/integrations.go` holds event streaming, notifications, reverse proxy domains/tokens, ingress, EDR, IdP sync, MSP and billing; EDR and IdP sync use one union struct each (`EDRIntegration`, `IdPSync`) over the per-vendor endpoints
- secrets come back masked (`****`) on GET: edit commands must drop them via `editTransformed` + `stripMasked` rather than PUT them back
