package agents

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestDefaultConfigRootCandidateParity(t *testing.T) {
	home := t.TempDir()

	claude := DefaultConfigRootCandidate(model.AgentClaudeCode, home)
	if claude.Source != ConfigRootDefault {
		t.Fatalf("claude source = %q, want %q", claude.Source, ConfigRootDefault)
	}
	if claude.BasePath != filepath.Join(home, ".claude") {
		t.Fatalf("claude base path = %q", claude.BasePath)
	}
	if claude.ConfigDir != filepath.Join(home, ".claude") {
		t.Fatalf("claude config dir = %q", claude.ConfigDir)
	}
	if !claude.Valid {
		t.Fatalf("claude default candidate should be valid: %s", claude.ValidationMessage)
	}

	codex := DefaultConfigRootCandidate(model.AgentCodex, home)
	if codex.Source != ConfigRootDefault {
		t.Fatalf("codex source = %q, want %q", codex.Source, ConfigRootDefault)
	}
	if codex.BasePath != filepath.Join(home, ".codex") {
		t.Fatalf("codex base path = %q", codex.BasePath)
	}
	if codex.ConfigDir != filepath.Join(home, ".codex") {
		t.Fatalf("codex config dir = %q", codex.ConfigDir)
	}
	if !codex.Valid {
		t.Fatalf("codex default candidate should be valid: %s", codex.ValidationMessage)
	}

	gemini := DefaultConfigRootCandidate(model.AgentGeminiCLI, home)
	if gemini.Source != ConfigRootDefault {
		t.Fatalf("gemini source = %q, want %q", gemini.Source, ConfigRootDefault)
	}
	if gemini.BasePath != filepath.Join(home, ".gemini") {
		t.Fatalf("gemini base path = %q", gemini.BasePath)
	}
	if gemini.ConfigDir != filepath.Join(home, ".gemini") {
		t.Fatalf("gemini config dir = %q", gemini.ConfigDir)
	}
	if !gemini.Valid {
		t.Fatalf("gemini default candidate should be valid: %s", gemini.ValidationMessage)
	}
}

func TestResolveClaudeEnvironmentCandidateValidLayouts(t *testing.T) {
	t.Run("environment root is direct Claude config root", func(t *testing.T) {
		home := t.TempDir()
		configDir := filepath.Join(home, ".claude-work")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir claude env root: %v", err)
		}

		candidate := ResolveClaudeEnvironmentCandidate(configDir)
		if candidate.Source != ConfigRootEnvironment {
			t.Fatalf("source = %q", candidate.Source)
		}
		if candidate.BasePath != configDir {
			t.Fatalf("base path = %q, want %q", candidate.BasePath, configDir)
		}
		if candidate.ConfigDir != configDir {
			t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, configDir)
		}
		if !candidate.Valid {
			t.Fatalf("candidate should be valid: %s", candidate.ValidationMessage)
		}
	})

	t.Run("legacy nested .claude directory remains valid when selected directly", func(t *testing.T) {
		home := t.TempDir()
		configDir := filepath.Join(home, ".claude-work", ".claude")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir direct .claude: %v", err)
		}

		candidate := ResolveClaudeEnvironmentCandidate(configDir)
		if candidate.BasePath != configDir {
			t.Fatalf("base path = %q, want %q", candidate.BasePath, configDir)
		}
		if candidate.ConfigDir != configDir {
			t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, configDir)
		}
		if !candidate.Valid {
			t.Fatalf("candidate should be valid: %s", candidate.ValidationMessage)
		}
	})
}

func TestResolveClaudeEnvironmentCandidateInvalidPath(t *testing.T) {
	home := t.TempDir()
	invalid := filepath.Join(home, "missing-parent", "not-claude")

	candidate := ResolveClaudeEnvironmentCandidate(invalid)
	if candidate.Valid {
		t.Fatalf("candidate should be invalid")
	}
	if candidate.ValidationMessage == "" {
		t.Fatalf("expected validation message for invalid candidate")
	}
}

func TestResolveCodexManualCandidate(t *testing.T) {
	home := t.TempDir()
	manual := filepath.Join(home, ".codex-work")

	candidate := ResolveCodexManualCandidate(manual)
	if candidate.Source != ConfigRootManual {
		t.Fatalf("source = %q, want %q", candidate.Source, ConfigRootManual)
	}
	if candidate.BasePath != manual {
		t.Fatalf("base path = %q, want %q", candidate.BasePath, manual)
	}
	if candidate.ConfigDir != manual {
		t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, manual)
	}
	if !candidate.Valid {
		t.Fatalf("manual codex candidate should be valid: %s", candidate.ValidationMessage)
	}
}

func TestResolveCodexEnvironmentCandidate(t *testing.T) {
	t.Run("CODEX_HOME is direct Codex config root", func(t *testing.T) {
		home := t.TempDir()
		configDir := filepath.Join(home, ".codex-work")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir codex env root: %v", err)
		}

		candidate := ResolveCodexEnvironmentCandidate(configDir)
		if candidate.Source != ConfigRootEnvironment {
			t.Fatalf("source = %q, want %q", candidate.Source, ConfigRootEnvironment)
		}
		if candidate.Label != codexHomeEnv {
			t.Fatalf("label = %q, want %q", candidate.Label, codexHomeEnv)
		}
		if candidate.BasePath != configDir {
			t.Fatalf("base path = %q, want %q", candidate.BasePath, configDir)
		}
		if candidate.ConfigDir != configDir {
			t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, configDir)
		}
		if !candidate.Valid {
			t.Fatalf("candidate should be valid: %s", candidate.ValidationMessage)
		}
	})

	t.Run("empty CODEX_HOME is invalid", func(t *testing.T) {
		candidate := ResolveCodexEnvironmentCandidate("")
		if candidate.Valid {
			t.Fatalf("empty CODEX_HOME candidate should be invalid")
		}
		if candidate.ValidationMessage == "" {
			t.Fatalf("expected validation message for empty CODEX_HOME")
		}
	})
}

func TestResolveCodexManualCandidateRejectsFilePath(t *testing.T) {
	home := t.TempDir()
	manual := filepath.Join(home, ".codex-work")
	if err := os.WriteFile(manual, []byte("not-a-directory"), 0o644); err != nil {
		t.Fatalf("write manual path file: %v", err)
	}

	candidate := ResolveCodexManualCandidate(manual)
	if candidate.Valid {
		t.Fatalf("manual codex file path should be invalid")
	}
	if candidate.ValidationMessage == "" {
		t.Fatalf("expected validation message for file conflict")
	}
}

func TestResolveGeminiEnvironmentCandidate(t *testing.T) {
	t.Run("GEMINI_CLI_HOME is base path and resolves to nested .gemini", func(t *testing.T) {
		home := t.TempDir()
		base := filepath.Join(home, "gemini-work")
		configDir := filepath.Join(base, ".gemini")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir gemini config dir: %v", err)
		}

		candidate := ResolveGeminiEnvironmentCandidate(base)
		if candidate.Source != ConfigRootEnvironment {
			t.Fatalf("source = %q, want %q", candidate.Source, ConfigRootEnvironment)
		}
		if candidate.Label != geminiCLIHomeEnv {
			t.Fatalf("label = %q, want %q", candidate.Label, geminiCLIHomeEnv)
		}
		if candidate.BasePath != base {
			t.Fatalf("base path = %q, want %q", candidate.BasePath, base)
		}
		if candidate.ConfigDir != configDir {
			t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, configDir)
		}
		if !candidate.Valid {
			t.Fatalf("candidate should be valid: %s", candidate.ValidationMessage)
		}
	})

	t.Run("direct-looking env value still nests .gemini", func(t *testing.T) {
		home := t.TempDir()
		directLooking := filepath.Join(home, ".gemini")
		configDir := filepath.Join(directLooking, ".gemini")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatalf("mkdir nested gemini config dir: %v", err)
		}

		candidate := ResolveGeminiEnvironmentCandidate(directLooking)
		if candidate.BasePath != directLooking {
			t.Fatalf("base path = %q, want %q", candidate.BasePath, directLooking)
		}
		if candidate.ConfigDir != configDir {
			t.Fatalf("config dir = %q, want %q", candidate.ConfigDir, configDir)
		}
		if !candidate.Valid {
			t.Fatalf("candidate should be valid: %s", candidate.ValidationMessage)
		}
	})

	t.Run("empty GEMINI_CLI_HOME is invalid", func(t *testing.T) {
		candidate := ResolveGeminiEnvironmentCandidate("")
		if candidate.Valid {
			t.Fatalf("empty GEMINI_CLI_HOME candidate should be invalid")
		}
		if candidate.ValidationMessage == "" {
			t.Fatalf("expected validation message for empty GEMINI_CLI_HOME")
		}
	})
}
