# Proposal: Provider Config Root Selection

## Intent

Allow Gentle-AI to install, sync, verify, and uninstall managed artifacts against non-default agent configuration locations. This supports users who maintain separate personal and work agent accounts, for example:

```text
~/.claude-personal/.claude/
~/.claude-personal/.claude.json
~/.claude-work/.claude/
~/.claude-work/.claude.json
```

Today Gentle-AI assumes the current OS home directory and hardcodes Claude, Codex, and Gemini roots as `~/.claude`, `~/.codex`, and `~/.gemini`. That makes the default path safe, but it prevents users from targeting an active alternate account/profile without moving files around or creating shell-level workarounds.

## Scope

### In Scope

- Add a provider config-root abstraction that can resolve default, environment-derived, and manually supplied roots.
- Support Claude Code first-class custom root selection.
- Support Codex custom root selection with the same abstraction.
- Support `CODEX_HOME` as Codex's environment-derived config root.
- Support Gemini CLI custom root selection with source-backed semantics for `GEMINI_CLI_HOME`.
- Preserve current behavior when no alternate root is configured.
- Update CLI install/sync/verification path resolution to use the selected root.
- Update TUI detection/selection flows to surface alternate roots before writing managed files.
- Remove hardcoded Claude/Codex paths in Agent Builder and route them through adapter/config-root resolution.
- Document how users choose default vs work/personal roots.
- Add tests for default behavior, environment candidates, manual path candidates, validation, and fallback.

### Out of Scope

- Managing provider authentication itself.
- Switching Claude or Codex accounts inside their native CLIs.
- Migrating existing provider config directories.
- Auto-discovering every possible user-created profile directory under `~`.
- Changing OpenCode's own config directory or authentication model.

## User Experience

When Gentle-AI detects an alternate provider config root, it SHOULD ask users which target to use before installation/sync writes files:

| Choice | Meaning |
|--------|---------|
| Default | Use current behavior, e.g. `~/.claude` or `~/.codex` |
| Detected custom root | Use the root found from environment/session context |
| Other path | Let the user provide a path, then validate it before use |

For Claude, `CLAUDE_CONFIG_DIR` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat the environment value as the directory Claude Code reads directly. It MUST NOT append an extra `.claude` segment to that path.

This matters for work/personal layouts such as:

```text
CLAUDE_CONFIG_DIR=~/.claude-work
~/.claude-work/settings.json       # read by Claude Code
~/.claude-work/.claude/settings.json # wrong target for runtime assets
```

For Codex, `CODEX_HOME` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat `CODEX_HOME` as the directory Codex reads directly. It MUST NOT append an extra `.codex` segment.

For Gemini CLI, `GEMINI_CLI_HOME` SHOULD be considered a candidate when present in the session. Gentle-AI MUST treat this value as Gemini's base home and MUST resolve the provider config root as `<GEMINI_CLI_HOME>/.gemini`.

Gentle-AI MUST display the resolved write target before applying changes.

## Approach

1. **Introduce a config-root resolver**: Centralize root selection so adapters and components do not invent paths independently.
2. **Keep adapter APIs stable where possible**: Prefer a contextual path resolver or wrapper over broad switch statements.
3. **Normalize root semantics**: Track both the user-selected root and the resolved provider config directory to avoid ambiguity.
4. **Update CLI/TUI consumers**: Install/sync/uninstall mostly flow through adapters already, but Agent Builder has hardcoded Claude/Codex paths and must be corrected.
5. **Document the workflow**: Add a concise guide explaining default vs custom account profiles and common examples.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/agents/*` | Modified | Claude/Codex/Gemini path resolution must support selected roots |
| `internal/system/config_scan.go` | Modified | Scan should include default and selected/available custom roots |
| `internal/agents/discovery.go` | Modified | Discovery should report resolved config dirs, not only default dirs |
| `internal/cli/run.go` | Modified | Install/apply/verification should use selected roots |
| `internal/cli/sync.go` | Modified | Sync should use selected roots and backup correct targets |
| `internal/components/*` | Modified | Components should keep using adapters, but receive selected config context |
| `internal/components/uninstall/service.go` | Modified | Uninstall must target the same selected roots used at install/sync time |
| `internal/tui/model.go` | Modified | TUI and Agent Builder must stop hardcoding Claude/Codex paths |
| `docs/` or README | Modified | Add user-facing custom config root documentation |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Writing to the wrong provider account | Medium | Always show resolved write targets and validate expected layout |
| Ambiguous root semantics | High | Store selected root + resolved provider config directory separately |
| Breaking default installs | Low | Default behavior remains unchanged and heavily tested |
| Misreading Codex root semantics | Low | Use verified `CODEX_HOME` semantics: the env value is the direct config root |
| Misreading Gemini home semantics | Medium | Use verified Gemini docs: `GEMINI_CLI_HOME` is a base path and config lives in `<base>/.gemini` |
| TUI/CLI divergence | Medium | Route both through the same resolver and state model |

## Rollback Plan

1. Remove config-root resolver and selection state.
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
- [ ] TUI displays default/custom/manual root choices before writes.
- [ ] CLI offers a non-interactive way to select a root.
- [ ] Agent Builder writes generated agents and SDD references to the selected provider root.
- [ ] Documentation explains personal/work account examples.
- [ ] `go test ./...` passes.
