package agents

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

const (
	claudeConfigDirEnv = "CLAUDE_CONFIG_DIR"
	codexHomeEnv       = "CODEX_HOME"
	geminiCLIHomeEnv   = "GEMINI_CLI_HOME"
)

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

func DefaultConfigRootCandidate(agentID model.AgentID, homeDir string) ConfigRootCandidate {
	basePath := defaultConfigRootForAgent(agentID, homeDir)
	candidate := ConfigRootCandidate{
		AgentID:  agentID,
		Source:   ConfigRootDefault,
		Label:    "default",
		BasePath: basePath,
	}

	switch agentID {
	case model.AgentClaudeCode:
		candidate.ConfigDir, candidate.Valid, candidate.ValidationMessage = normalizeClaudeRoot(basePath)
	case model.AgentCodex:
		candidate.ConfigDir, candidate.Valid, candidate.ValidationMessage = normalizeCodexRoot(basePath)
	case model.AgentGeminiCLI:
		candidate.ConfigDir, candidate.Valid, candidate.ValidationMessage = normalizeCodexRoot(basePath)
	default:
		candidate.ConfigDir = basePath
		candidate.Valid = false
		candidate.ValidationMessage = fmt.Sprintf("unsupported provider %q", agentID)
	}

	return candidate
}

func ResolveClaudeEnvironmentCandidate(envValue string) ConfigRootCandidate {
	basePath := filepath.Clean(envValue)
	configDir, valid, message := normalizeClaudeRoot(basePath)
	if envValue == "" {
		valid = false
		message = fmt.Sprintf("%s is not set", claudeConfigDirEnv)
	}

	return ConfigRootCandidate{
		AgentID:           model.AgentClaudeCode,
		Source:            ConfigRootEnvironment,
		Label:             claudeConfigDirEnv,
		BasePath:          basePath,
		ConfigDir:         configDir,
		Valid:             valid,
		ValidationMessage: message,
	}
}

func ResolveCodexManualCandidate(manualPath string) ConfigRootCandidate {
	basePath := filepath.Clean(manualPath)
	configDir, valid, message := normalizeCodexRoot(basePath)

	return ConfigRootCandidate{
		AgentID:           model.AgentCodex,
		Source:            ConfigRootManual,
		Label:             "manual",
		BasePath:          basePath,
		ConfigDir:         configDir,
		Valid:             valid,
		ValidationMessage: message,
	}
}

func ResolveCodexEnvironmentCandidate(envValue string) ConfigRootCandidate {
	basePath := filepath.Clean(envValue)
	configDir, valid, message := normalizeCodexRoot(basePath)
	if envValue == "" {
		valid = false
		message = fmt.Sprintf("%s is not set", codexHomeEnv)
	}

	return ConfigRootCandidate{
		AgentID:           model.AgentCodex,
		Source:            ConfigRootEnvironment,
		Label:             codexHomeEnv,
		BasePath:          basePath,
		ConfigDir:         configDir,
		Valid:             valid,
		ValidationMessage: message,
	}
}

func ResolveGeminiEnvironmentCandidate(envValue string) ConfigRootCandidate {
	basePath := filepath.Clean(envValue)
	configDir, valid, message := normalizeGeminiHomeOverride(basePath)
	if envValue == "" {
		valid = false
		message = fmt.Sprintf("%s is not set", geminiCLIHomeEnv)
	}

	return ConfigRootCandidate{
		AgentID:           model.AgentGeminiCLI,
		Source:            ConfigRootEnvironment,
		Label:             geminiCLIHomeEnv,
		BasePath:          basePath,
		ConfigDir:         configDir,
		Valid:             valid,
		ValidationMessage: message,
	}
}

func defaultConfigRootForAgent(agentID model.AgentID, homeDir string) string {
	switch agentID {
	case model.AgentClaudeCode:
		return filepath.Join(homeDir, ".claude")
	case model.AgentCodex:
		return filepath.Join(homeDir, ".codex")
	case model.AgentGeminiCLI:
		return filepath.Join(homeDir, ".gemini")
	default:
		return homeDir
	}
}

func normalizeClaudeRoot(basePath string) (string, bool, string) {
	if msg := ValidateNoFileDirectoryConflict(basePath); msg != "" {
		return basePath, false, msg
	}
	cleanBase := filepath.Clean(basePath)
	if msg := ValidateDirectoryExistsOrParentExists(cleanBase); msg != "" {
		return cleanBase, false, msg
	}

	return cleanBase, true, ""
}

func normalizeCodexRoot(basePath string) (string, bool, string) {
	if msg := ValidateNoFileDirectoryConflict(basePath); msg != "" {
		return basePath, false, msg
	}
	if msg := ValidateDirectoryExistsOrParentExists(basePath); msg != "" {
		return basePath, false, msg
	}

	return basePath, true, ""
}

func normalizeGeminiHomeOverride(basePath string) (string, bool, string) {
	cleanBase := filepath.Clean(basePath)
	configDir := filepath.Join(cleanBase, ".gemini")

	if msg := ValidateNoFileDirectoryConflict(configDir); msg != "" {
		return configDir, false, msg
	}
	if msg := ValidateDirectoryExistsOrParentExists(configDir); msg != "" {
		return configDir, false, msg
	}

	return configDir, true, ""
}

func ValidateDirectoryExists(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "directory does not exist"
		}
		return fmt.Sprintf("cannot stat path: %v", err)
	}
	if !info.IsDir() {
		return "path exists but is not a directory"
	}

	return ""
}

func ValidateParentExists(path string) string {
	parent := filepath.Dir(path)
	if parent == "." || parent == path {
		return ""
	}

	info, err := os.Stat(parent)
	if err != nil {
		if os.IsNotExist(err) {
			return "parent directory does not exist"
		}
		return fmt.Sprintf("cannot stat parent directory: %v", err)
	}
	if !info.IsDir() {
		return "parent exists but is not a directory"
	}

	return ""
}

func ValidateNoFileDirectoryConflict(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		return fmt.Sprintf("cannot stat path: %v", err)
	}
	if !info.IsDir() {
		return "path exists as a file; directory expected"
	}

	return ""
}

func ValidateDirectoryExistsOrParentExists(path string) string {
	if msg := ValidateDirectoryExists(path); msg == "" {
		return ""
	}

	return ValidateParentExists(path)
}
