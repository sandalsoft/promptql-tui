# PromptQL TUI

A terminal user interface (TUI) client for the [Hasura PromptQL](https://hasura.io/docs/promptql/) API, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Uses the [PromptQL Go SDK](https://github.com/sandalsoft/promptql-sdk-golang) for API communication.

## Features

- **Project browser** — List and select your PromptQL projects
- **Thread management** — Create new conversation threads or resume existing ones
- **Interactive chat** — Send natural language queries to PromptQL and see responses inline
- **Direct query mode** — Use an API key + DDN URL for stateless queries
- **Persistent config** — Credentials saved to `~/.config/promptql-tui/config.json`
- **Environment overrides** — Set `PROMPTQL_PAT`, `PROMPTQL_API_KEY`, `PROMPTQL_DDN_URL`

## Install

```bash
go install github.com/sandalsoft/promptql-tui/cmd/promptql-tui@latest
```

Or build from source:

```bash
git clone https://github.com/sandalsoft/promptql-tui.git
cd promptql-tui
go build -o promptql-tui ./cmd/promptql-tui/
```

## Usage

```bash
# Run with config file (interactive setup on first run)
./promptql-tui

# Or set credentials via environment
export PROMPTQL_PAT="your-personal-access-token"
./promptql-tui

# For direct query mode (no thread management)
export PROMPTQL_API_KEY="your-api-key"
export PROMPTQL_DDN_URL="https://your-project.ddn.hasura.app/graphql"
./promptql-tui
```

## Navigation

| View | Key | Action |
|------|-----|--------|
| All | `ctrl+c` | Quit |
| All | `esc` | Go back |
| Setup | `tab`/`shift+tab` | Navigate fields |
| Setup | `enter` | Save and continue |
| Projects | `j`/`k` or arrows | Navigate list |
| Projects | `enter` | Select project |
| Projects | `r` | Refresh |
| Projects | `s` | Go to setup |
| Threads | `j`/`k` or arrows | Navigate list |
| Threads | `enter` | Select/resume thread |
| Threads | `n` | New thread |
| Threads | `r` | Refresh |
| Chat | `ctrl+s` | Send message |

## Architecture

```
cmd/promptql-tui/    # Entry point
internal/
  config/            # Persistent configuration (~/.config/promptql-tui/)
  sdk/               # Vendored PromptQL Go SDK
  tui/               # Bubble Tea TUI (views, styles, messages)
```

## Configuration

On first run, the TUI presents a setup screen to configure:

| Field | Required | Description |
|-------|----------|-------------|
| PAT | Yes | Hasura Personal Access Token |
| API Key | No | Runtime API key for direct queries |
| DDN URL | No | DDN GraphQL endpoint URL |
| Timezone | No | Defaults to UTC |

Config is stored at `~/.config/promptql-tui/config.json`.

## SDK

The embedded SDK (in `internal/sdk/`) is vendored from [sandalsoft/promptql-sdk-golang](https://github.com/sandalsoft/promptql-sdk-golang) and provides:

- **Projects** — List, lookup, enable/disable PromptQL
- **Threads** — Create, list, send messages, get events
- **Query** — Execute natural language queries via REST
- **Prompts** — CRUD operations on sample prompts
- **API Keys** — Generate and manage runtime API keys
- **Users** — List and lookup PromptQL users
