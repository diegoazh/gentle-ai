# Tasks: Provider Config Root Selection

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

## Phase 2: Selection State and CLI Wiring

- [ ] T-09 Extend install/sync selection state to carry provider config-root selections.
- [ ] T-10 Add CLI flag parsing for non-interactive provider root selection.
- [ ] T-11 Add CLI option to opt into environment-derived provider roots.
- [ ] T-12 Update install runtime to pass selected roots into adapter/component path resolution.
- [ ] T-13 Update sync runtime and backup target calculation to use selected roots.
- [ ] T-14 Update post-install/post-sync verification paths to use selected roots.
- [ ] T-15 Add CLI tests for parsing, dry-run output, default parity, and custom Claude/Codex/Gemini root behavior.

## Phase 3: Adapter and Component Path Integration

- [ ] T-16 Update Claude adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-17 Update Codex adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-17b Update Gemini adapter path resolution to honor selected config root while preserving default behavior.
- [ ] T-18 Update SDD injection tests to verify Claude/Gemini writes under a custom root.
- [ ] T-19 Update skills injection tests to verify Claude/Codex/Gemini custom roots.
- [ ] T-20 Update Engram/MCP/permissions tests for custom root path consistency.
- [ ] T-21 Update discovery/config scan behavior to report selected or environment-derived roots where appropriate.

## Phase 4: TUI Root Selection

- [ ] T-22 Add TUI state for provider root candidates and selections.
- [ ] T-23 Add root selection screen or step when alternate candidates are present.
- [ ] T-24 Add manual path input and validation feedback.
- [ ] T-25 Display resolved write target and safety note before installation/sync.
- [ ] T-26 Add TUI tests for default, environment candidate, manual invalid path, and selected custom path.

## Phase 5: Agent Builder Cleanup

- [ ] T-27 Replace `agentBuilderSkillsDir()` hardcoded Claude/Codex paths with resolver/adapter-backed paths.
- [ ] T-28 Replace `agentBuilderSystemPromptPath()` hardcoded Claude/Codex paths with resolver/adapter-backed paths.
- [ ] T-29 Ensure generated custom agents install into the selected provider root.
- [ ] T-30 Ensure Agent Builder SDD reference injection uses the selected provider root.
- [ ] T-31 Add Agent Builder tests for Claude, Codex, and Gemini custom roots.

## Phase 6: Uninstall and Continuity

- [ ] T-32 Persist selected provider roots in install state.
- [ ] T-33 Update sync to reuse or prompt with saved provider root selections.
- [ ] T-34 Update uninstall to target saved custom roots.
- [ ] T-35 Add uninstall tests proving custom roots are cleaned and default roots are untouched.
- [ ] T-36 Add fallback behavior for missing saved custom roots with clear user feedback.

## Phase 7: Documentation

- [ ] T-37 Document default vs custom provider roots.
- [ ] T-38 Add Claude personal/work account example using direct roots like `~/.claude-work` and `~/.claude-personal`.
- [ ] T-39 Document `CLAUDE_CONFIG_DIR` detection and direct-root normalization rules.
- [ ] T-40 Document `CODEX_HOME` detection and direct-root normalization rules.
- [ ] T-40b Document `GEMINI_CLI_HOME` detection and base-home normalization rules (`<base>/.gemini`).
- [ ] T-41 Add safety note that Gentle-AI does not switch provider authentication.

## Phase 8: Verification

- [ ] T-42 Run focused resolver tests.
- [ ] T-43 Run CLI tests.
- [ ] T-44 Run TUI tests.
- [ ] T-45 Run component path tests.
- [ ] T-46 Run Agent Builder tests.
- [ ] T-47 Run uninstall tests.
- [ ] T-48 Run `go test ./...`.
- [ ] T-49 Run `go vet ./...`.

## Suggested Chain Boundaries

### PR 1: Resolver Foundation

Delivers candidate/selection types, Claude/Codex/Gemini normalization, validation, unit tests, and source-backed documentation. No runtime install/sync behavior changes.

### PR 2: CLI Path Selection

Adds CLI flags/options, install/sync path usage, dry-run reporting, and verification/backup path updates.

### PR 3: Component and Adapter Integration

Ensures SDD, skills, MCP, Engram, permissions, and discovery/config scan consistently use selected roots for Claude, Codex, and Gemini.

### PR 4: TUI Selection UX

Adds root-selection UI, manual path validation, and resolved-target display.

### PR 5: Agent Builder and Uninstall Continuity

Removes hardcoded Agent Builder Claude/Codex paths and persists/reuses selections for sync/uninstall.

### PR 6: Documentation and Final Verification

Adds docs, examples, safety notes, and final full-suite verification.
