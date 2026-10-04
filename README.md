

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

Environment variables (override the file): `NETBIRD_CLI_URL`, `NETBIRD_CLI_TOKEN`, `NETBIRD_CLI_CONFIG_FILE`, `NETBIRD_CLI_EDITION`.

### NetBird Cloud vs self-hosted

Some features only exist on one edition. The CLI detects it from the URL (`*.netbird.io` is NetBird Cloud, anything else is self-hosted) and adapts: commands the edition does not support are hidden from help and completion, and running one explains why instead of returning a raw 404.

| Only on NetBird Cloud | Only on self-hosted |
|---|---|
| event streaming, notification channels, EDR + peer bypass, IdP sync (Google, Azure, Okta SCIM, SCIM), ingress peers/ports, MSP tenants, billing, network traffic events, managed Agent Network gateway | `setup`, `instancestatus`, user invites (embedded IdP) |

Force the edition if the detection is wrong (e.g. a custom domain in front of NetBird Cloud):

```yaml
edition: cloud        # cloud | selfhosted | auto (default)
```

`netbird info` shows the URL, the edition (and how it was decided), the management version, and the commands hidden for this edition. On self-hosted servers, a 404 on an optional feature (reverse proxy, Agent Network, embedded IdP) comes with a hint that it may not be deployed or need a newer version.

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
netbird events <type>                # audit, traffic, agent
netbird bypass|unbypass peer <name>  # EDR compliance bypass
netbird sync idp <google|azure>      # trigger an IdP user/group sync
netbird validate domain <domain>     # validate a custom reverse proxy domain
netbird regenerate scimtoken <okta|scim>
netbird msp <verify-dns|invite|accept|decline|unlink|subscribe> <tenant>
netbird billing <usage|subscription|plans|invoices|invoice|portal|checkout|aws>
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
| `agentproviders` | `ap` | | `agentpolicies` | `apol` |
| `guardrails` | `gr` | | `budgetrules` | `br` |
| `agentsettings` | `as` | | `agentcatalog` | `acat` |
| `agentusage` | `au` | | `agentconsumption` | `acons` |
| `agentmodels` | `am` | | `agentgateway` | `agw` |
| `eventstreams` | `es` | | `notificationchannels` | `nc` |
| `notificationtypes` | `nt` | | `idpsyncs` | `ids` |
| `domains` | `dom` | | `proxytokens` | `pxt` |
| `ingresspeers` | `ip` | | `ingressports` | `iport` |
| `tenants` | `tn` | | `edrbypassed` | `bypassed` |

## Supported resources

`users`, `groups`, `peers`, `policies`, `networks`, `networkresources`, `setupkeys`, `posturechecks`, `routes`, `nameservers`, `dnssettings`, `dnszones`, `dnsrecords`, `accounts`, `tokens`, `services`, `proxyclusters`, `identityproviders`, `jobs`, `auditevents`, `trafficevents`, `proxylogs`, `countries`, `cities`, `instancestatus`, `instanceversion`

Integrations: `eventstreams`, `notificationchannels`, `notificationtypes`, `edr`, `edrbypassed`, `idpsyncs`, `idpsynclogs`

Reverse proxy & ingress: `domains`, `proxytokens`, `proxyclusters`, `ingresspeers`, `ingressports`

MSP & billing: `tenants`, `msp ...`, `billing ...`

Agent Network: `agentproviders`, `agentpolicies`, `guardrails`, `budgetrules`, `agentsettings`, `agentgateway`, `agentcatalog`, `agentmodels`, `agentconfig`, `agentusage`, `agentconsumption`, `events agent`

## Integrations

```bash
# Stream audit events to a SIEM (datadog, s3, firehose, generic_http)
netbird create eventstream --platform datadog --config api_key=$DD_API_KEY,api_url=https://http-intake.logs.datadoghq.eu

# Notifications (see event codes: netbird get notificationtypes)
netbird create notificationchannel --type email --emails ops@example.com --events user.join,peer.add
netbird create notificationchannel --type webhook --url https://hooks.example.com/nb --header Authorization="Bearer x" --events peer.add

# EDR: only compliant devices can connect (intune, sentinelone, falcon, huntress, fleetdm)
netbird create edr intune --client-id <app-id> --tenant-id <tenant> --secret $SECRET --groups laptops
netbird get edr
netbird bypass peer laptop-bob                   # let one non-compliant peer through

# IdP user/group sync (google, azure, okta, scim)
netbird create idpsync google --customer-id C01abc --service-account-key-file sa.json --group-prefixes eng-
netbird sync idp google
netbird get idpsynclogs google
netbird create idpsync okta --connection-name okta-prod   # prints the SCIM token once

# Reverse proxy: custom domains and self-hosted proxy tokens
netbird create domain --domain apps.example.com --cluster eu.proxy.netbird.io
netbird validate domain apps.example.com
netbird create proxytoken --name eu-proxy        # token shown once

# Ingress peers / port forwarding (cloud)
netbird create ingresspeer web-1
netbird create ingressport --peer web-1 --name https --range 443/tcp
```

Secrets the API returns masked (`****`: event stream credentials, webhook header values, EDR/IdP credentials) are left out of `netbird edit`, so they are never overwritten with the mask; pass the flag again (`--config`, `--header`, `--secret`, ...) to change them.

## Agent Network (AI agents)

[Agent Network](https://netbird.io) is NetBird's gateway for AI agents: LLM traffic from peers goes through a NetBird endpoint (reachable only over the tunnel) that injects the upstream API key, enforces who can reach which provider, applies guardrails and token/USD caps, and logs usage and cost.

```bash
# 1. Bootstrap (once): self-hosted proxy cluster, or a managed gateway on NetBird Cloud
netbird create agentsettings --proxy-address proxy.example.com
netbird create agentgateway                      # cloud: managed gateway (idempotent)

# 2. Add an LLM provider (catalog ids: netbird get agentcatalog)
netbird get agentcatalog openai_api              # models + default prices
netbird create agentprovider --name openai --type openai_api \
    --url https://api.openai.com --api-key "$OPENAI_API_KEY" --models gpt-4o-mini,gpt-4o
netbird get agentmodels --provider openai        # what the key can actually reach

# 3. Guardrails, access policy, budgets
netbird create guardrail --name audited --models gpt-4o-mini --prompt-capture --redact-pii
netbird create agentpolicy --name devs-llm --source-groups devs --providers openai \
    --guardrails audited --budget-user-cap 10 --window 24h
netbird create budgetrule --name devs-daily --groups devs --budget-user-cap 5 --window 24h
netbird describe agentpolicy devs-llm            # groups, providers, guardrails, limits resolved

# 4. Observe
netbird get agentconfig                          # endpoint + providers available to *you*
netbird get agentusage --since 30d --granularity week
netbird get agentconsumption                     # live counters vs. caps
netbird events agent --since 24h --decision deny # access logs (deny reasons)
netbird events agent --sessions --user alice@example.com
```

Everything is also declarative: `netbird generate agentnetwork > ai.yaml`, then `netbird apply -f ai.yaml`. Resources are matched by name; groups, users, providers and guardrails can be referenced by name, and `api_key: ${OPENAI_API_KEY}` is expanded from the environment so keys stay out of the file.

## Limitations

- Cloud-only and self-hosted-only features are hidden on the other edition (see [NetBird Cloud vs self-hosted](#netbird-cloud-vs-self-hosted)).
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
