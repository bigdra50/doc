package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds configuration for provider creation
type Config struct {
	ProviderType string `toml:"provider"`

	// API Keys
	OpenAIAPIKey    string `toml:"openai_api_key"`
	AnthropicAPIKey string `toml:"anthropic_api_key"`

	// Claude Code CLI path
	ClaudeCodePath string `toml:"claude_code_path"`

	// Model Selection
	OpenAIModel    string `toml:"openai_model"`
	AnthropicModel string `toml:"anthropic_model"`
	ClaudeModel    string `toml:"claude_model"`

	// General settings
	Verbose bool `toml:"verbose"`
}

// ProviderType constants
const (
	ProviderTypeClaude    = "claude-code"
	ProviderTypeOpenAI    = "openai"
	ProviderTypeAnthropic = "anthropic"
)

// GetDefaultModel returns the default model for a provider
func GetDefaultModel(provider string) string {
	switch provider {
	case ProviderTypeOpenAI:
		return "gpt-4o-mini" // Most cost-effective balanced option
	case ProviderTypeAnthropic:
		return "claude-3-5-haiku-20241022" // Most cost-effective recent option
	case ProviderTypeClaude:
		return "sonnet" // Claude Code CLI default
	default:
		return ""
	}
}

// Load loads configuration from config file, then environment variables
func Load() Config {
	// Start with defaults
	config := Config{
		ProviderType:   ProviderTypeClaude,
		ClaudeCodePath: "claude",
		OpenAIModel:    GetDefaultModel(ProviderTypeOpenAI),
		AnthropicModel: GetDefaultModel(ProviderTypeAnthropic),
		ClaudeModel:    GetDefaultModel(ProviderTypeClaude),
		Verbose:        false,
	}

	// Load from config file if it exists
	if configPath := GetConfigPath(); configPath != "" {
		if fileConfig, err := loadFromFile(configPath); err == nil {
			mergeConfig(&config, fileConfig)
		}
	}

	// Override with environment variables and .env file
	loadEnvFile()
	config = overrideWithEnv(config)

	return config
}

// LoadFromEnv loads provider configuration from environment variables and .env file (deprecated, use Load())
func LoadFromEnv() Config {
	return Load()
}

// GetConfigPath returns the path to the config file following XDG Base Directory spec
func GetConfigPath() string {
	configDir := GetConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, "config.toml")
}

// GetConfigDir returns the directory containing the config file following XDG Base Directory spec
func GetConfigDir() string {
	// Check XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, getConfigSubdir())
	}

	// Fallback to ~/.config/bigdra50/doc
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", getConfigSubdir())
}

// getConfigSubdir returns the subdirectory name for config files
// Always uses organization prefix to avoid conflicts
func getConfigSubdir() string {
	return filepath.Join("bigdra50", "doc")
}

// SaveConfig saves the config to the config file
func SaveConfig(config Config) error {
	configDir := GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configPath := GetConfigPath()
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}

	return nil
}

// loadFromFile loads configuration from a TOML file
func loadFromFile(path string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	return config, err
}

// mergeConfig merges fileConfig into config (fileConfig takes precedence for non-empty values)
func mergeConfig(config *Config, fileConfig Config) {
	if fileConfig.ProviderType != "" {
		config.ProviderType = fileConfig.ProviderType
	}
	if fileConfig.OpenAIAPIKey != "" {
		config.OpenAIAPIKey = fileConfig.OpenAIAPIKey
	}
	if fileConfig.AnthropicAPIKey != "" {
		config.AnthropicAPIKey = fileConfig.AnthropicAPIKey
	}
	if fileConfig.ClaudeCodePath != "" {
		config.ClaudeCodePath = fileConfig.ClaudeCodePath
	}
	if fileConfig.OpenAIModel != "" {
		config.OpenAIModel = fileConfig.OpenAIModel
	}
	if fileConfig.AnthropicModel != "" {
		config.AnthropicModel = fileConfig.AnthropicModel
	}
	if fileConfig.ClaudeModel != "" {
		config.ClaudeModel = fileConfig.ClaudeModel
	}
	// Verbose is handled separately by CLI flags
}

// overrideWithEnv overrides config values with environment variables
func overrideWithEnv(config Config) Config {
	config.ProviderType = getEnvOrDefault("LLM_PROVIDER", config.ProviderType)
	config.OpenAIAPIKey = getEnvOrDefault("OPENAI_API_KEY", config.OpenAIAPIKey)
	config.AnthropicAPIKey = getEnvOrDefault("ANTHROPIC_API_KEY", config.AnthropicAPIKey)
	config.ClaudeCodePath = getEnvOrDefault("CLAUDE_CODE_PATH", config.ClaudeCodePath)
	config.OpenAIModel = getEnvOrDefault("OPENAI_MODEL", config.OpenAIModel)
	config.AnthropicModel = getEnvOrDefault("ANTHROPIC_MODEL", config.AnthropicModel)
	config.ClaudeModel = getEnvOrDefault("CLAUDE_MODEL", config.ClaudeModel)

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist or can't be opened, silently continue
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set in environment
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}
