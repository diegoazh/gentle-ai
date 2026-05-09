# Provider Config Root Selection — Specification

## Purpose

Define the behavior required for Gentle-AI to safely target default or custom provider configuration roots for Claude Code, Codex, and Gemini CLI without breaking existing installs.

---

## 1. Root Model

### Requirement: Config Root Candidate

The system MUST represent provider config choices as structured candidates containing:

- provider agent ID
- source (`default`, `environment`, `manual`, `saved-selection`)
- display label
- selected base path
- resolved provider config directory
- validation status
- validation message

#### Scenario: Default candidate

- GIVEN the current user home is `/home/me`
- WHEN Claude candidates are resolved
- THEN a default candidate exists with resolved config directory `/home/me/.claude`

#### Scenario: Candidate includes source

- GIVEN a candidate was created from `CLAUDE_CONFIG_DIR`
- WHEN the candidate is displayed
- THEN its source is `environment` and its label identifies the environment variable

### Requirement: Root Values Match Provider Runtime Semantics

The system MUST resolve environment variables according to what the provider runtime actually reads, not according to Gentle-AI's historical default layout.

#### Scenario: Claude custom root from environment

- GIVEN `CLAUDE_CONFIG_DIR=/Users/me/.claude-work`
- WHEN the root is normalized for Claude
- THEN the resolved provider config directory is `/Users/me/.claude-work`
- AND Gentle-AI MUST NOT append an extra `.claude` segment

#### Scenario: Legacy nested Claude directory selected manually

- GIVEN a user manually selects `/Users/me/.claude-work/.claude`
- WHEN the root is normalized for Claude
- THEN the resolved provider config directory is `/Users/me/.claude-work/.claude`

#### Scenario: Gemini environment value is a base home, not direct config root

- GIVEN `GEMINI_CLI_HOME=/tmp/gemini-job-123`
- WHEN the root is normalized for Gemini CLI
- THEN the resolved provider config directory is `/tmp/gemini-job-123/.gemini`
- AND Gentle-AI MUST append exactly one `.gemini` segment to the cleaned base path

---

## 2. Claude Root Resolution

### Requirement: Default Claude Root

Claude MUST continue to use `~/.claude` when no custom root is selected.

#### Scenario: No environment and no manual selection

- GIVEN `CLAUDE_CONFIG_DIR` is unset
- AND the user does not choose a manual root
- WHEN Gentle-AI installs Claude artifacts
- THEN files are written under `~/.claude` exactly as before

### Requirement: Claude Environment Candidate

When `CLAUDE_CONFIG_DIR` is set, Gentle-AI SHOULD present it as a selectable candidate before writing Claude artifacts.

#### Scenario: Environment variable points to custom Claude root

- GIVEN `CLAUDE_CONFIG_DIR=/Users/me/.claude-work`
- WHEN candidates are resolved
- THEN a valid environment candidate is shown
- AND its resolved provider config directory is `/Users/me/.claude-work`

#### Scenario: Environment variable points to nested legacy directory

- GIVEN `CLAUDE_CONFIG_DIR=/Users/me/.claude-work/.claude`
- WHEN candidates are resolved
- THEN a valid environment candidate is shown
- AND its resolved provider config directory is `/Users/me/.claude-work/.claude`

#### Scenario: Environment variable points to invalid path

- GIVEN `CLAUDE_CONFIG_DIR=/tmp/missing-parent/not-claude`
- AND the path cannot be created safely because its parent does not exist
- WHEN candidates are resolved
- THEN the candidate is marked invalid
- AND Gentle-AI MUST NOT write Claude artifacts there unless the user fixes the path

---

## 3. Codex Root Resolution

### Requirement: Default Codex Root

Codex MUST continue to use `~/.codex` when no custom root is selected.

#### Scenario: Default Codex install

- GIVEN no custom Codex root is selected
- WHEN Gentle-AI installs Codex artifacts
- THEN files are written under `~/.codex` exactly as before

### Requirement: Manual Codex Root

Gentle-AI SHOULD allow a validated manual Codex root selection.

#### Scenario: Manual direct Codex config directory

- GIVEN the user selects `/Users/me/.codex-work`
- AND the path is accepted as a Codex config directory
- WHEN Gentle-AI installs Codex artifacts
- THEN `agents.md`, `skills/`, and `config.toml` are written under `/Users/me/.codex-work`

### Requirement: Codex Environment Candidate

When `CODEX_HOME` is set, Gentle-AI SHOULD present it as a selectable candidate before writing Codex artifacts.

#### Scenario: Environment variable points to custom Codex root

- GIVEN `CODEX_HOME=/Users/me/.codex-work`
- WHEN candidates are resolved
- THEN a valid environment candidate is shown
- AND its resolved provider config directory is `/Users/me/.codex-work`
- AND Gentle-AI MUST NOT append an extra `.codex` segment

#### Scenario: Empty CODEX_HOME

- GIVEN `CODEX_HOME` is empty
- WHEN candidates are resolved
- THEN the candidate is marked invalid
- AND Gentle-AI MUST NOT write Codex artifacts there

---

## 4. Gemini Root Resolution

### Requirement: Default Gemini Root

Gemini CLI MUST continue to use `~/.gemini` when no custom root is selected.

#### Scenario: Default Gemini install

- GIVEN no custom Gemini root is selected
- WHEN Gentle-AI installs Gemini artifacts
- THEN files are written under `~/.gemini` exactly as before

### Requirement: Gemini Environment Candidate

When `GEMINI_CLI_HOME` is set, Gentle-AI SHOULD present it as a selectable candidate before writing Gemini artifacts.
This requirement is source-backed by Gemini CLI upstream docs: configuration reference and enterprise reference.

#### Scenario: Environment variable points to Gemini base home

- GIVEN `GEMINI_CLI_HOME=/Users/me/gemini-work`
- WHEN candidates are resolved
- THEN a valid environment candidate is shown
- AND its base path is `/Users/me/gemini-work`
- AND its resolved provider config directory is `/Users/me/gemini-work/.gemini`

#### Scenario: Direct-looking environment value still nests `.gemini`

- GIVEN `GEMINI_CLI_HOME=/Users/me/.gemini`
- WHEN candidates are resolved
- THEN a valid environment candidate is shown
- AND its resolved provider config directory is `/Users/me/.gemini/.gemini`

#### Scenario: Empty GEMINI_CLI_HOME

- GIVEN `GEMINI_CLI_HOME` is empty
- WHEN candidates are resolved
- THEN the candidate is marked invalid
- AND Gentle-AI MUST NOT write Gemini artifacts there

---

## 5. CLI Behavior

### Requirement: Non-Interactive Root Selection

The CLI MUST expose a non-interactive way to choose provider config roots for automation.

#### Scenario: CLI custom Claude root

- GIVEN the user runs install/sync with a custom Claude root option
- WHEN the command executes
- THEN all Claude component paths resolve through the selected root

### Requirement: CLI Safety Output

The CLI SHOULD print or report the resolved provider write target before applying changes.

#### Scenario: Dry run shows target

- GIVEN a dry run with a custom Claude root
- WHEN the plan is shown
- THEN the output includes the resolved Claude config directory

---

## 6. TUI Behavior

### Requirement: Root Selection Screen

The TUI SHOULD present root choices when an alternate root is detected or when the user requests manual selection.

#### Scenario: Claude environment candidate displayed

- GIVEN `CLAUDE_CONFIG_DIR` is set to a valid alternate root
- WHEN the user reaches agent/config selection
- THEN the TUI shows both default and environment-derived Claude choices

### Requirement: Manual Path Validation

The TUI MUST validate manual paths before allowing install/sync to proceed.

#### Scenario: Invalid manual path

- GIVEN the user enters an invalid manual Claude path
- WHEN validation runs
- THEN the TUI shows the validation reason
- AND install/sync cannot continue with that candidate

---

## 7. Component Path Consistency

### Requirement: Components Use Resolved Adapter Paths

Components MUST write to paths resolved from the selected provider root, not hardcoded defaults.

#### Scenario: SDD injection honors custom Claude root

- GIVEN Claude root resolves to `/Users/me/.claude-work/.claude`
- WHEN SDD injection runs
- THEN `CLAUDE.md`, `skills/`, commands, output styles, and sub-agents are written under that directory

#### Scenario: Engram injection honors custom Codex root

- GIVEN Codex root resolves to `/Users/me/.codex-work`
- WHEN Engram injection runs
- THEN `config.toml`, `engram-instructions.md`, and `engram-compact-prompt.md` are written under `/Users/me/.codex-work`

#### Scenario: Gemini components honor derived config directory

- GIVEN Gemini base root resolves from `GEMINI_CLI_HOME=/Users/me/gemini-work`
- WHEN system prompt, skills, or MCP settings are written
- THEN files are written under `/Users/me/gemini-work/.gemini`

### Requirement: Agent Builder Uses Resolver

Agent Builder MUST stop hardcoding Claude and Codex paths and MUST use the same selected-root resolver as install/sync.

#### Scenario: Custom agent installed to selected Claude root

- GIVEN Agent Builder installs a generated custom agent
- AND Claude root resolves to `/Users/me/.claude-work/.claude`
- WHEN installation completes
- THEN the custom skill and SDD reference are written under `/Users/me/.claude-work/.claude`

---

## 8. Persistence

### Requirement: Remember Selected Roots For Sync/Uninstall

When a root selection is applied, Gentle-AI SHOULD persist enough state for later sync and uninstall operations to target the same provider root.

#### Scenario: Sync after custom install

- GIVEN a Claude install used `/Users/me/.claude-work/.claude`
- WHEN the user runs sync later
- THEN Gentle-AI can reuse or prompt with the previously selected root

#### Scenario: Uninstall after custom install

- GIVEN managed files were written to a custom Claude root
- WHEN uninstall runs for Claude components
- THEN uninstall targets the same custom root, not only `~/.claude`

---

## 9. Documentation

### Requirement: User Guide

Documentation MUST explain how to use custom provider roots for personal/work accounts.

#### Scenario: Claude work/personal example

- GIVEN a user has `~/.claude-work/.claude` and `~/.claude-personal/.claude`
- WHEN they read the docs
- THEN they can identify which path to pass or select in the CLI/TUI

#### Scenario: Gemini home semantics explained

- GIVEN a user reads provider root docs
- WHEN they review Gemini examples
- THEN they understand `GEMINI_CLI_HOME` is a base directory and effective config root is `<base>/.gemini`

### Requirement: Safety Notes

Documentation MUST warn users that Gentle-AI does not switch provider authentication; it only targets where managed files are written.

---

## 10. Verification

### Requirement: Test Coverage

The implementation MUST include tests for default roots, environment candidates, manual validation, CLI/TUI path use, Agent Builder path use, sync, and uninstall.

#### Scenario: Full Go test suite

- GIVEN implementation is complete
- WHEN `go test ./...` runs
- THEN all tests pass
