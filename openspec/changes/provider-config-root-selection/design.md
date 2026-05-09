# Design: Provider Config Root Selection

## Outcome

Introduce a single config-root resolution layer that lets Gentle-AI target default or custom provider roots for Claude, Codex, and Gemini CLI while preserving the adapter-driven architecture.

## Current Behavior

| Area | Current Behavior |
|------|------------------|
| Claude adapter | Hardcodes paths under `homeDir/.claude` |
| Codex adapter | Hardcodes paths under `homeDir/.codex` |
| Gemini adapter | Hardcodes paths under `homeDir/.gemini` |
| Detection | Uses `adapter.GlobalConfigDir(homeDir)` or fixed scan entries |
| CLI install/sync | Passes one `homeDir` into adapters/components |
| TUI Agent Builder | Hardcodes Claude/Codex skills and prompt paths directly |

The adapter pattern is good. The missing concept is that `homeDir` is not always the right root for provider-managed files.

## Architecture Decision Records

| Decision | Choice | Alternatives | Rationale |
|----------|--------|--------------|-----------|
| Root abstraction | Add provider config-root resolver | Add ad-hoc Claude-only env checks | A shared resolver avoids duplicating path logic across install, sync, TUI, Agent Builder, and uninstall. |
| Root semantics | Treat provider env vars as direct provider config roots | Append default subdirectories to env roots | The provider runtime owns env var semantics. For Claude, `CLAUDE_CONFIG_DIR=/path/.claude-work` means Claude reads `/path/.claude-work`, not `/path/.claude-work/.claude`. |
| Default behavior | Preserve current `~/.claude` and `~/.codex` | Force prompt every install | Backwards compatibility matters; users without custom roots should see no behavior change. |
| Claude env detection | Use `CLAUDE_CONFIG_DIR` as candidate | Ignore env vars | The session env is a strong signal that the user is running in a non-default Claude profile. |
| Codex env detection | Use `CODEX_HOME` as a direct config-root candidate | Manual-only Codex support | Codex docs and installed binary confirm `CODEX_HOME` controls config/auth home and defaults to `~/.codex`. |
| Gemini env detection | Use `GEMINI_CLI_HOME` as base-home candidate and resolve `<base>/.gemini` | Treat `GEMINI_CLI_HOME` as direct config root | Official Gemini docs define `~/.gemini` as config location and document `GEMINI_CLI_HOME` as a unique home/base directory where Gemini creates/uses `.gemini`. |
| TUI integration | Shared resolver/state consumed by screens | TUI-only duplicate logic | CLI and TUI must behave consistently. |
| Agent Builder | Replace hardcoded paths with resolver/adapter paths | Leave as follow-up | Leaving hardcoded paths would silently write custom agents to the wrong account. |

## Proposed Types

```go
type ConfigRootSource string

const (
    ConfigRootDefault        ConfigRootSource = "default"
    ConfigRootEnvironment    ConfigRootSource = "environment"
    ConfigRootManual         ConfigRootSource = "manual"
    ConfigRootSavedSelection ConfigRootSource = "saved-selection"
)

type ConfigRootCandidate struct {
    AgentID           model.AgentID
    Source            ConfigRootSource
    Label             string
    BasePath          string
    ConfigDir         string
    Valid             bool
    ValidationMessage string
}

type ConfigRootSelection struct {
    AgentID   model.AgentID
    Source    ConfigRootSource
    BasePath  string
    ConfigDir string
}
```

The exact package is implementation-defined, but `internal/agents` or a small `internal/configroots` package are good candidates. Avoid putting this in `system` if it creates an import cycle with adapters.

## Resolver Flow

```text
Resolve candidates(agent, osHome, env)
  │
  ├─ default candidate
  │    Claude: osHome/.claude
  │    Codex:  osHome/.codex
  │    Gemini: osHome/.gemini
  │
  ├─ environment candidates
  │    Claude: CLAUDE_CONFIG_DIR when set
  │    Codex: CODEX_HOME when set
  │    Gemini: GEMINI_CLI_HOME when set (base path; resolved config dir is <base>/.gemini)
  │
  └─ optional manual candidate
       validate provider-specific layout

User/CLI/TUI selects candidate
  │
  ▼
Selection passed into install/sync/uninstall runtime
  │
  ▼
Adapters/components resolve paths through selection
```

## Claude Normalization

Claude environment roots must match Claude Code runtime behavior. `CLAUDE_CONFIG_DIR` points to the directory Claude reads directly.

| Input | Meaning | Resolved ConfigDir |
|-------|---------|--------------------|
| `/Users/me/.claude-work` | Direct Claude config root from `CLAUDE_CONFIG_DIR` | `/Users/me/.claude-work` |
| `/Users/me/.claude-work/.claude` | Legacy/manual nested config directory | `/Users/me/.claude-work/.claude` |

Validation SHOULD accept an existing directory or a creatable directory whose parent exists. If the path exists as a file, reject it.

Do not auto-append `.claude` to `CLAUDE_CONFIG_DIR`. That was the original failure mode: Gentle-AI-managed assets landed under `.claude-work/.claude`, while Claude Code read `.claude-work` directly.

## Codex Normalization

Codex is simpler because Gentle-AI currently treats `~/.codex` as the config directory itself, and Codex's `CODEX_HOME` follows the same direct-root semantics.

| Input | Meaning | Resolved ConfigDir |
|-------|---------|--------------------|
| `/Users/me/.codex-work` | Direct Codex config directory | `/Users/me/.codex-work` |

Validation SHOULD accept an existing directory or allow creating it when the parent exists, depending on existing install conventions. If the path exists as a file, reject it.

Do not auto-append `.codex` to `CODEX_HOME`. Codex reads `CODEX_HOME/config.toml`, `CODEX_HOME/auth.json`, and other runtime files directly.

## Gemini Normalization

Gemini semantics differ from Claude/Codex direct-root env values. By default, Gemini CLI stores config in `~/.gemini`, and when `GEMINI_CLI_HOME` is set, Gemini treats that value as a base home and creates/uses a `.gemini` directory under it.

Sources:

- Gemini configuration reference: <https://github.com/google-gemini/gemini-cli/blob/main/docs/reference/configuration.md>
- Gemini enterprise reference: <https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/enterprise.md>

| Input | Meaning | Resolved ConfigDir |
|-------|---------|--------------------|
| `/Users/me/gemini-work` | Base home from `GEMINI_CLI_HOME` | `/Users/me/gemini-work/.gemini` |
| `/Users/me/.gemini` | Base home that looks like config dir | `/Users/me/.gemini/.gemini` |

Validation SHOULD run against the resolved config directory (`<base>/.gemini`) and follow the same file/parent safety rules as other providers.

## Adapter Integration Options

### Preferred

Add a runtime path context used by adapters without changing every component signature drastically.

Possible shape:

```go
type PathResolver interface {
    ConfigDir(agent model.AgentID, homeDir string) string
}
```

Then install/sync runtimes can wrap adapters or pass a resolver to components.

### Acceptable

Extend `model.Selection` with provider root selections and pass that selection into runtime/component layers.

This is more invasive but explicit.

### Avoid

Hardcoding `CLAUDE_CONFIG_DIR` inside individual component functions. That repeats logic and makes sync/uninstall inconsistent.

## CLI Design

CLI should support automation. Candidate names are examples and can be adjusted during implementation:

```text
gentle-ai install --agent claude-code --agent-config-root claude-code=/Users/me/.claude-work
gentle-ai sync --agent-config-root claude-code=/Users/me/.claude-work
gentle-ai install --use-provider-env-roots
```

Dry runs SHOULD include:

```text
Claude Code config root: /Users/me/.claude-work (source: environment CLAUDE_CONFIG_DIR)
Gemini CLI config root: /Users/me/gemini-work/.gemini (source: environment GEMINI_CLI_HOME)
```

## TUI Design

When alternate candidates exist:

```text
Claude Code configuration

Where should Gentle-AI install Claude Code artifacts?

› Default: ~/.claude
  Environment: ~/.claude-work
  Other path...

Gentle-AI writes prompts, skills, MCP config, commands, output styles, and sub-agents here.
It does not switch your Claude login.
```

Manual path entry must validate before moving forward.

## State and Uninstall

The selected root should be persisted with install state so future sync/uninstall can target the same root.

State should record at least:

- agent ID
- source
- selected base path
- resolved config dir

Uninstall must prefer saved selections for managed files. This is important because uninstalling from `~/.claude` after installing to `~/.claude-work/.claude` would look successful but leave managed files behind.

## Documentation Plan

Add a user-facing section to the install/sync docs explaining:

- default behavior
- Claude work/personal layouts
- `CLAUDE_CONFIG_DIR` detection
- manual path selection
- Codex custom root selection
- Gemini `GEMINI_CLI_HOME` base-home semantics (`<base>/.gemini`)
- safety note: this writes Gentle-AI files only; it does not authenticate or switch provider accounts

## Testing Strategy

| Layer | Tests |
|-------|-------|
| Resolver unit tests | default, Claude env profile home, Claude direct config, invalid path, Codex manual root, Gemini base-home env nesting |
| Adapter tests | Claude/Codex path methods honor selected root where integration point exists |
| CLI tests | flags parse into selections; dry run reports resolved root |
| TUI tests | alternate candidate appears; invalid manual path blocks continue |
| Component tests | SDD, skills, engram, MCP, permissions write under selected root (including Gemini resolved config dir) |
| Agent Builder tests | generated skills and SDD references use selected root |
| Uninstall tests | saved custom roots are removed/cleaned instead of defaults |

## Rollout

This should be implemented in slices because it touches core pathing and UX:

1. Resolver + Claude/Codex/Gemini path tests with default parity.
2. CLI install/sync path selection.
3. TUI root selection.
4. Agent Builder hardcoded path cleanup.
5. Uninstall/saved selection continuity.
6. Documentation.
