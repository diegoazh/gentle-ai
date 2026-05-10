# Tasks: Provider Config Root Selection

## Answer First

Implement as **registry-first**:

1. Persist provider profiles (source of truth).
2. Treat env vars as candidates/hints.
3. Wire canonical management via CLI + TUI.
4. Keep upgrade/sync opportunistic (register suggestions, not ownership).

## Review Workload Forecast

| Field | Forecast |
|-------|----------|
| Estimated changed lines | 700-1100 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Decision needed before apply | No, delivery strategy is `auto-chain` |
| Recommended chain strategy | Ask before apply: `stacked-to-main` or `feature-branch-chain` |

This feature affects shared path resolution, CLI, TUI, Agent Builder, uninstall, docs, and tests. It should be delivered as chained PRs unless a maintainer grants `size:exception`.

## Phase 1: Resolver Foundation

- [x] T-01 Add provider config-root candidate and selection types.
- [x] T-02 Implement default candidate resolution for Claude (`~/.claude`), Codex (`~/.codex`), and Gemini (`~/.gemini`).
- [x] T-03 Implement Claude environment candidate resolution from `CLAUDE_CONFIG_DIR`.
- [x] T-04 Implement Claude path normalization using direct `CLAUDE_CONFIG_DIR` root semantics.
- [x] T-05 Implement Codex manual root normalization.
- [x] T-06 Implement Codex environment candidate resolution from `CODEX_HOME` using direct root semantics.
- [x] T-06b Implement Gemini environment candidate resolution from `GEMINI_CLI_HOME` using base-home semantics (`<base>/.gemini`).
- [x] T-07 Add validation helpers for directory existence, parent existence, and file-vs-directory conflicts.
- [x] T-08 Add resolver unit tests for default parity, valid env paths, invalid env paths, and manual candidates, including Gemini `.gemini` nesting behavior.

## Phase 2: Profile Registry Foundation and CLI Wiring

- [ ] T-09 Add persisted provider profile registry (provider + profile name + base/config dir + source + stale state).
- [ ] T-10 Extend install/sync selection state to carry provider profile selections.
- [ ] T-11 Add CLI profile management commands (`profiles add|list|remove|update`).
- [ ] T-12 Add CLI flag parsing for non-interactive provider profile/root selection.
- [ ] T-13 Add CLI option to opt into environment-derived provider root candidates.
- [ ] T-14 Update install runtime to pass selected roots into adapter/component path resolution.
- [ ] T-15 Update sync runtime and backup target calculation to use selected roots.
- [ ] T-16 Update post-install/post-sync verification paths to use selected roots.
- [ ] T-17 Add CLI tests for profile CRUD, parsing, dry-run output, default parity, and custom Claude/Codex/Gemini root behavior.

## Phase 3: Adapter and Component Path Integration

- [ ] T-18 Update Claude adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-19 Update Codex adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-19b Update Gemini adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-20 Update SDD injection tests to verify Claude/Gemini writes under a custom root.
- [ ] T-21 Update skills injection tests to verify Claude/Codex/Gemini custom roots.
- [ ] T-22 Update Engram/MCP/permissions tests for custom root path consistency.
- [ ] T-23 Update discovery/config scan behavior to report selected or environment-derived roots where appropriate.

## Phase 4: TUI Root Selection

- [ ] T-24 Add TUI state for provider profiles, candidates, and selections.
- [ ] T-25 Add `Manage provider profiles` flow (add/list/edit/remove).
- [ ] T-26 Add root/profile selection step when alternate candidates are present.
- [ ] T-27 Add manual path input and validation feedback.
- [ ] T-28 Display resolved write target and safety note before installation/sync.
- [ ] T-29 Add TUI tests for default, environment candidate, manual invalid path, and selected custom path.

## Phase 5: Agent Builder Cleanup

- [ ] T-30 Replace `agentBuilderSkillsDir()` hardcoded Claude/Codex paths with resolver/adapter-backed paths.
- [ ] T-31 Replace `agentBuilderSystemPromptPath()` hardcoded Claude/Codex paths with resolver/adapter-backed paths.
- [ ] T-32 Ensure generated custom agents install into the selected provider root.
- [ ] T-33 Ensure Agent Builder SDD reference injection uses the selected provider root.
- [ ] T-34 Add Agent Builder tests for Claude, Codex, and Gemini custom roots.

## Phase 6: Uninstall and Continuity

- [ ] T-35 Persist selected provider roots in install state.
- [ ] T-36 Update sync to reuse or prompt with saved provider root selections.
- [ ] T-37 Update uninstall to target saved custom roots.
- [ ] T-38 Add uninstall tests proving custom roots are cleaned and default roots are untouched.
- [ ] T-39 Add fallback behavior for missing saved custom roots with clear user feedback.
- [ ] T-40 Mark/remove stale profiles when registered paths disappear and print repair guidance.

## Phase 7: Documentation

- [ ] T-41 Document mental model: profiles are persisted targets; env vars are hints.
- [ ] T-42 Add Claude personal/work account example using direct roots like `~/.claude-work` and `~/.claude-personal`.
- [ ] T-43 Document `CLAUDE_CONFIG_DIR` detection and direct-root normalization rules.
- [ ] T-44 Document `CODEX_HOME` detection and direct-root normalization rules.
- [ ] T-44b Document `GEMINI_CLI_HOME` detection and base-home normalization rules (`<base>/.gemini`).
- [ ] T-45 Document CLI/TUI canonical management surfaces and opportunistic sync/upgrade behavior.
- [ ] T-46 Add safety note that Gentle-AI does not switch provider authentication.
- [ ] T-47 Document OpenCode difference (internal provider profiles/subscriptions).

## Phase 8: Verification

- [ ] T-48 Run focused resolver tests.
- [ ] T-49 Run CLI tests.
- [ ] T-50 Run TUI tests.
- [ ] T-51 Run component path tests.
- [ ] T-52 Run Agent Builder tests.
- [ ] T-53 Run uninstall tests.
- [ ] T-54 Run `go test ./...`.
- [ ] T-55 Run `go vet ./...`.

## Suggested Chain Boundaries

### PR 1: Resolver + Registry Design Foundation

Delivers candidate/selection types, Claude/Codex/Gemini normalization, validation, unit tests, and registry-first specs/design. No registry persistence or runtime install/sync behavior changes.

### PR 2: Registry Persistence + CLI Profile Management

Adds persisted profile registry storage, `profiles` CRUD commands, CLI flags/options, install/sync path usage, dry-run reporting, and verification/backup path updates.

### PR 3: Component and Adapter Integration

Ensures SDD, skills, MCP, Engram, permissions, and discovery/config scan consistently use selected roots for Claude, Codex, and Gemini.

### PR 4: TUI Manage Profiles + Selection UX

Adds Manage provider profiles flow, root-selection UI, manual path validation, and resolved-target display.

### PR 5: Agent Builder and Uninstall Continuity

Removes hardcoded Agent Builder Claude/Codex paths and persists/reuses selections for sync/uninstall, including stale-path handling.

### PR 6: Documentation and Final Verification

Adds docs, examples, safety notes, and final full-suite verification.
