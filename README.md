# netbird-cli

Go CLI to interact with the [Netbird](https://netbird.io) API.

Inspired by `kubectl`: verb-resource interface, API-powered autocompletion, name→ID resolution, interactive editor, declarative apply.

## Installation

```bash
make build          # → ./netbird-cli
make install        # → /usr/local/bin/netbird-cli
```

Or via `go install`:

```bash
go install ./cmd/cli
```

## Configuration

YAML config file expected at `~/.config/netbird-cli/config.yaml`:

```yaml
auth:
  url: https://api.netbird.io       # or your instance URL (without /api)
  token: <your-api-token>
```

Environment variables (override the file): `NETBIRD_CLI_URL`, `NETBIRD_CLI_TOKEN`, `NETBIRD_CLI_CONFIG_FILE`.

Create the API token in the Netbird dashboard → Users → your user → Access Tokens.

## Usage

```
netbird get <resource> [name]        # list or display
netbird describe <resource> <name>   # full detail, IDs resolved to names
netbird create <resource> [flags]    # create (--edit for interactive editor)
netbird edit <resource> <name>       # edit (interactive editor)
netbird delete <resource> <name>     # delete
netbird apply -f <file.yaml>         # declarative config
netbird generate <resource|all>      # generate a template
```

### Examples

```bash
# List / display
netbird get users                   # table
netbird get users alice             # key-value detail
netbird get users alice -o json     # JSON output
netbird get groups                  # table (alias: grp)

# Inspect
netbird describe user alice         # linked groups, peers, policies
netbird describe group devs         # peers, resources, users
netbird describe peer stage-host-1  # groups, OS, IP, user

# Create
netbird create user --email alice@example.com --name Alice --role admin --auto-groups devs
netbird create user --edit          # opens editor with a template
netbird create group --name devs --peers alice-laptop
netbird create group --edit         # YAML template in editor

# Edit
netbird edit user alice             # opens $EDITOR, annotated YAML
netbird edit user alice --role admin # pre-fills "role: admin"

# Delete
netbird delete group devs

# Events
netbird events audit                  # last 50 events
netbird events audit --page-size 100  # last 100 events
netbird events audit --page-size 0    # all events
netbird events traffic                # cloud-only
netbird get proxylogs                 # cloud-only

# Special actions
netbird approve user bob
netbird block user bob
netbird unblock user bob
netbird invite create --email ... --name ... --role ... --auto-groups ...
netbird whoami               # current user
netbird setup --email ... --password ...  # initialize instance
```

### Global flags

| Flag | Description |
|------|-------------|
| `-o json`, `-o yaml` | Output format |
| `--dry-run` | Show what would be done without executing |
| `--edit` | (on `create`) Open editor with a template |

### Autocompletion

```bash
source <(netbird-cli completion bash)   # bash
source <(netbird-cli completion zsh)    # zsh
```

`netbird get <TAB>` → lists known resources.
`netbird get users <TAB>` → lists existing users (via API).

### Default output

Human-readable format: tables for lists, key-value for details. `-o json` or `-o yaml` for machine output.

### Interactive editor

`netbird edit` and `netbird create --edit` open `$EDITOR` (or `vim`) with an annotated YAML. IDs are resolved to names:

```yaml
auto_groups:
    - d7q5t33pu5as73fmoqd0  # admin
    - ch8i4ug6lnn4g9hqv7m0  # devs
```

Names accepted instead of IDs (`devs` instead of `ch8i4...`), resolved automatically. No saved changes = no request sent.

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

## Resources

`users`, `groups`, `peers`, `policies`, `networks`, `networkresources`, `setupkeys`, `posturechecks`, `routes`, `nameservers`, `dnssettings`, `dnszones`, `dnsrecords`, `accounts`, `tokens`, `services`, `proxyclusters`, `identityproviders`, `jobs`, `auditevents`, `trafficevents`, `proxylogs`, `countries`, `cities`, `instancestatus`, `instanceversion`

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
internal/client/  — Netbird API client (std lib only)
  client.go       —   Client, doRequest/doGet/doPost/doPut/doDelete, GetRaw/PutRaw/PostRaw, name→ID cache, IDToName
  *.go            —   one file per API resource
config/config.go  —   YAML loading + env vars
```

## License

BSD 3-Clause — see [LICENSE](LICENSE).
