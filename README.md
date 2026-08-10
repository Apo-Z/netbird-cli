

# netbird-cli

> [!WARNING]
> **Unofficial** — this CLI is not affiliated with, endorsed by, or maintained by Netbird. Use at your own risk.
> The Netbird API may introduce breaking changes at any time without notice.

A `kubectl`-style CLI to manage a [Netbird](https://netbird.io) instance via the REST API.

## Features

- Full CRUD (create, read, update, delete) for all Netbird resources
- Describe with cross-referencing (IDs resolved to human-readable names)
- Interactive YAML editor (`netbird edit`, `netbird create --edit`)
- Declarative config (`netbird apply -f <file.yaml>`)
- Shell autocompletion (bash, zsh)
- JSON / YAML / table / key-value output formats (`-o`)
- Name → ID resolution (use `alice` instead of `ch8i4ug6...`)
- Zero external HTTP dependencies (stdlib `net/http` only)
- Works with both cloud and self-hosted Netbird instances

## Prerequisites

- [Go](https://go.dev) 1.26+
- A Netbird API token (dashboard → Users → your user → Access Tokens)

## Installation

### From source (recommended)

```bash
git clone https://github.com/Apo-Z/netbird-cli.git
cd netbird-cli
make build          # → ./netbird-cli
make install        # → /usr/local/bin/netbird-cli (requires sudo)
```

### Manual build

```bash
go build -ldflags="-s -w" -o netbird-cli ./cmd/cli
```

> `go install ./cmd/cli` is **not recommended** — the resulting binary will be named `cli` instead of `netbird-cli`.

## Makefile

| Target | Description |
|--------|-------------|
| `build` | Compiles a statically-linked binary at `./netbird-cli` |
| `install` | Builds then copies to `/usr/local/bin/netbird-cli` (requires `sudo`) |
| `clean` | Removes `./netbird-cli` |

The binary is stripped (`-ldflags="-s -w"`) and relies only on the Go standard library — no CGO, no external C dependencies.

## Configuration

YAML config file expected at `~/.config/netbird-cli/config.yaml`:

```yaml
auth:
  url: https://api.netbird.io       # or your instance URL (without /api)
  token: <your-api-token>
```

Environment variables (override the file): `NETBIRD_CLI_URL`, `NETBIRD_CLI_TOKEN`, `NETBIRD_CLI_CONFIG_FILE`.

The URL is automatically cleaned — trailing `/` and `/api` are stripped.

## Usage

```
netbird get <resource> [name]        # list or display
netbird describe <resource> <name>   # full detail, IDs resolved to names
netbird create <resource> [flags]    # create (--edit for interactive editor)
netbird edit <resource> <name>       # edit (interactive editor)
netbird delete <resource> <name>     # delete
netbird apply -f <file.yaml>         # declarative config
netbird generate <resource|all>      # generate a template

# Special verbs
netbird approve user <name>          # approve a pending user
netbird block user <name>            # block a user
netbird unblock user <name>          # unblock a user
netbird invite create ...            # invite a user
netbird whoami                       # current authenticated user
netbird setup --email ...            # initialize a fresh instance
netbird events <type>                # audit, traffic, proxylogs
```

### Examples

```bash
# List / display
netbird get users                   # table
netbird get users alice             # key-value detail
netbird get users alice -o json     # JSON output
netbird get groups                  # table (aliases: grp)

# Inspect
netbird describe user alice         # linked groups, peers, policies
netbird describe group devs         # peers, resources, users
netbird describe peer stage-host-1  # groups, OS, IP, user

# Create
netbird create user --email alice@example.com --name Alice --role admin
netbird create user --edit          # opens editor with a template
netbird create group --name devs --peers alice-laptop
netbird create group --edit         # YAML template in editor

# Edit
netbird edit user alice             # opens $EDITOR, annotated YAML

# Delete
netbird delete group devs

# Events
netbird events audit                  # last 50 events
netbird events audit --page-size 100  # last 100 events
netbird events audit --page-size 0    # all events
```

### Global flags

| Flag | Description |
|------|-------------|
| `-o json`, `-o yaml` | Output format (default: human-readable) |
| `--dry-run` | Show what would be done without executing |
| `--edit` | (on `create`) Open editor with a template |

### Autocompletion

```bash
source <(netbird-cli completion bash)   # bash
source <(netbird-cli completion zsh)    # zsh
```

`netbird get <TAB>` → lists known resources.
`netbird get users <TAB>` → lists existing users fetched from the API.

### Interactive editor

`netbird edit` and `netbird create --edit` open `$EDITOR` (defaults to `vim`) with an annotated YAML. IDs are resolved to names:

```yaml
auto_groups:
    - d7q5t33pu5as73fmoqd0  # admin
    - ch8i4ug6lnn4g9hqv7m0  # devs
```

Names are accepted instead of IDs (`devs` instead of `ch8i4...`) and resolved automatically. If the file is unchanged, no request is sent.

### Declarative apply

```bash
netbird generate all > setup.yaml
vim setup.yaml
netbird apply -f setup.yaml --dry-run
netbird apply -f setup.yaml
```

### Aliases

| Command | Alias | | Command | Alias |
|----------|-------|-|----------|-------|
| `users` | `us` | | `nameservers` | `ns` |
| `groups` | `grp` | | `setupkeys` | `sk` |
| `peers` | `pr` | | `posturechecks` | `pc` |
| `networks` | `nw` | | `services` | `svc` |
| `dnszones` | `dz` | | `dnsrecords` | `dr` |
| `accounts` | `ac` | | `identityproviders` | `idp` |
| `proxyclusters` | `pxc` | | `networkresources` | `nwr` |

## Supported resources

`users`, `groups`, `peers`, `policies`, `networks`, `networkresources`, `setupkeys`, `posturechecks`, `routes`, `nameservers`, `dnssettings`, `dnszones`, `dnsrecords`, `accounts`, `tokens`, `services`, `proxyclusters`, `identityproviders`, `jobs`, `auditevents`, `trafficevents`, `proxylogs`, `countries`, `cities`, `instancestatus`, `instanceversion`

## Limitations

- `netbird events traffic` and `netbird events proxylogs` are **cloud-only** — they will return 404 on self-hosted instances.
- Edits use HTTP `PUT` (full object replacement), not `PATCH`. Always send the complete object body.
- `GET /api/users/{id}` does not exist in the Netbird API — user lookup is done by listing all users.
- No CI / tests yet.

## Project structure

```
cmd/cli/          — Cobra CLI (kubectl-style verbs + resources)
  main.go         —   entrypoint
  root.go         —   rootCmd, getCmd, createCmd, editCmd, deleteCmd, applyCmd, generateCmd, describeCmd
  flags.go        —   shared flag variables
  completion.go   —   dynamic autocompletion via API
  editor.go       —   interactive YAML editing (GET → $EDITOR → PUT / template → $EDITOR → POST)
  format.go       —   table / key-value / JSON / YAML output
  apply.go        —   declarative config -f
  generate.go     —   YAML templates
  describe.go     —   rich detail with cross-referencing
  *.go            —   one file per resource
internal/client/  — Netbird API client (stdlib only)
  client.go       —   Client, doRequest/doGet/doPost/doPut/doDelete, GetRaw/PutRaw/PostRaw, name→ID cache, IDToName
  *.go            —   one file per API resource
config/config.go  —   YAML loading + env vars
```

## License

BSD 3-Clause — see [LICENSE](LICENSE).
