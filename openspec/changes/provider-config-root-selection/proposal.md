# Proposal: Provider Config Root Selection

## Answer First

Provider roots are **saved provider profiles**, not a single active environment-derived selection.

- Environment variables are **hints/candidates only**.
- Source of truth is a persisted **provider profile registry** managed by CLI + TUI.
- A provider can have **multiple profiles** (for example, Claude `work` and `personal`).

## Intent

Allow Gentle-AI to install, sync, verify, and uninstall managed artifacts against persisted provider profile roots, including non-default locations. This supports users who maintain separate personal and work agent accounts, for example:

```text
~/.claude-personal/.claude/
~/.claude-personal/.claude.json
~/.claude-work/.claude/
~/.claude-work/.claude.json
```

Today Gentle-AI assumes the current OS home directory and hardcodes Claude, Codex, and Gemini roots as `~/.claude`, `~/.codex`, and `~/.gemini`. That makes the default path safe, but it prevents users from managing multiple durable targets across sessions.

## Scope

### In Scope

- Add a provider profile registry abstraction that persists named roots per provider.
- Support Claude Code first-class custom root selection.
- Support Codex custom root selection with the same abstraction.
- Support `CODEX_HOME` as Codex environment-derived candidate root.
- Support Gemini CLI custom root selection with source-backed semantics for `GEMINI_CLI_HOME`.
- Preserve current behavior when no extra profile is configured.
- Update CLI install/sync/verification path resolution to use the selected root.
- Add CLI profile management (`profiles add|list|remove|update`).
- Add TUI profile management flow (`Manage provider profiles`).
- Update TUI detection/selection flows to surface registered roots and opportunistic candidates before writing managed files.
- Remove hardcoded Claude/Codex paths in Agent Builder and route them through adapter/config-root resolution.
- Document how users manage default/work/personal profiles.
- Add tests for default behavior, environment candidates, manual path candidates, validation, and fallback.

### Out of Scope

- Managing provider authentication itself.
- Switching Claude or Codex accounts inside their native CLIs.
- Migrating existing provider config directories.
- Auto-discovering every possible user-created profile directory under `~`.
- Changing OpenCode's own config directory or authentication model.

## Canonical Management Surfaces

| Surface | Purpose | Canonical? |
|--------|---------|------------|
| CLI (`gentle-ai profiles ...`) | Add/list/remove/update provider profiles | Yes |
| TUI (`Manage provider profiles`) | Interactive profile CRUD and validation | Yes |
| `upgrade` / `sync` env detection | Opportunistic suggestions to register | No |

## User Experience

When Gentle-AI needs a target for install/sync/upgrade, it SHOULD ask users which **registered provider profile** to use. If no relevant profile exists, it SHOULD help them register one.

| Choice | Meaning |
|--------|---------|
| Registered profile | Use persisted path, e.g. `claude-code/work` |
| Register detected env path | Promote env candidate into a saved profile |
| Add another path | Create a new profile interactively |

For Claude, `CLAUDE_CONFIG_DIR` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat the environment value as the directory Claude Code reads directly. It MUST NOT append an extra `.claude` segment to that path.

This matters for work/personal layouts such as:

```text
CLAUDE_CONFIG_DIR=~/.claude-work
~/.claude-work/settings.json       # read by Claude Code
~/.claude-work/.claude/settings.json # wrong target for runtime assets
```

For Codex, `CODEX_HOME` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat `CODEX_HOME` as the directory Codex reads directly. It MUST NOT append an extra `.codex` segment.

For Gemini CLI, `GEMINI_CLI_HOME` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat this value as Gemini's base home and MUST resolve the provider config root as `<GEMINI_CLI_HOME>/.gemini`.

If a registered profile path disappears, Gentle-AI MUST:

1. warn that updates/sync will skip that target,
2. mark/remove it as stale in the registry, and
3. show how to re-add or update it later.

If an environment path is detected during upgrade/sync and is not registered, Gentle-AI MUST offer to register it and optionally add other profiles.

Gentle-AI MUST display the resolved write target before applying changes.

## Approach

1. **Introduce a provider profile registry**: Persist named profile targets by provider.
2. **Introduce a config-root resolver**: Resolve runtime targets from selected registered profiles and validated candidates.
2. **Keep adapter APIs stable where possible**: Prefer a contextual path resolver or wrapper over broad switch statements.
3. **Normalize root semantics**: Track profile metadata, selected base path, and resolved provider config directory to avoid ambiguity.
4. **Update CLI/TUI consumers**: Install/sync/uninstall mostly flow through adapters already, but Agent Builder has hardcoded Claude/Codex paths and must be corrected.
5. **Document the workflow**: Add a concise guide explaining default vs custom account profiles and common examples.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/agents/*` | Modified | Claude/Codex/Gemini path resolution must support selected roots |
| `internal/system/config_scan.go` | Modified | Scan should include default and selected/available custom roots |
| `internal/agents/discovery.go` | Modified | Discovery should report resolved config dirs, not only default dirs |
| `internal/cli/run.go` | Modified | Install/apply/verification should use selected registered profiles |
| `internal/cli/sync.go` | Modified | Sync should use selected profiles and backup correct targets |
| `internal/cli/profiles*.go` | Added/Modified | CLI profile CRUD management |
| `internal/components/*` | Modified | Components should keep using adapters, but receive selected config context |
| `internal/components/uninstall/service.go` | Modified | Uninstall must target the same selected roots used at install/sync time |
| `internal/tui/model.go` | Modified | TUI and Agent Builder must stop hardcoding Claude/Codex paths |
| `internal/tui/*profiles*` | Added/Modified | Manage provider profiles flow |
| `docs/` or README | Modified | Add user-facing custom config root documentation |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Writing to the wrong provider account | Medium | Always show resolved write targets and validate expected layout |
| Ambiguous root semantics | High | Persist profile registry entries + resolved provider config directory separately |
| Breaking default installs | Low | Default behavior remains unchanged and heavily tested |
| Misreading Codex root semantics | Low | Use verified `CODEX_HOME` semantics: the env value is the direct config root |
| Misreading Gemini home semantics | Medium | Use verified Gemini docs: `GEMINI_CLI_HOME` is a base path and config lives in `<base>/.gemini` |
| TUI/CLI divergence | Medium | Route both through the same registry + resolver model |
| Stale profile paths after directory moves | Medium | Detect missing path, mark stale, explain how to repair (`profiles update` / TUI manage flow) |

## Rollback Plan

1. Remove provider profile registry, resolver, and selection state.
2. Revert Claude/Codex adapters to direct `filepath.Join(homeDir, ".claude")` and `filepath.Join(homeDir, ".codex")` behavior.
3. Revert CLI/TUI path-selection screens and flags.
4. Remove docs for custom provider roots.
5. Existing default config directories remain untouched.

## Success Criteria

- [ ] Existing default install/sync/uninstall behavior remains unchanged.
- [ ] Claude can target a profile home such as `~/.claude-work` containing `.claude/`.
- [ ] Claude can target a direct config directory such as `/path/to/.claude` when validated.
- [ ] Codex can target a validated non-default config root.
- [ ] Gemini CLI can target a validated non-default base path and resolve config under `<base>/.gemini`.
- [ ] TUI exposes Manage provider profiles and install/sync selection from registered profiles.
- [ ] CLI exposes profile management commands and non-interactive root/profile selection.
- [ ] Agent Builder writes generated agents and SDD references to the selected provider root.
- [ ] Documentation explains personal/work account examples and env-vars-as-hints mental model.
- [ ] `go test ./...` passes.
