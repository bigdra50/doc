package main

import (
	"context"
	"fmt"

	"github.com/bigdra50/doc/internal/config"
)

// TranslationResponse represents the result of a translation operation
type TranslationResponse struct {
	Content   string `json:"content"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

// TranslationOptions holds configuration for translation operations
type TranslationOptions struct {
	TargetLanguage     string
	CustomInstruction  string
	PreserveFormat     bool
	Verbose           bool
}

// LLMProvider defines the interface for different LLM providers
type LLMProvider interface {
	// Translate translates the given content using the specified options
	Translate(ctx context.Context, content string, options TranslationOptions) (*TranslationResponse, error)
	
	// ValidateConfig validates the provider configuration
	ValidateConfig() error
	
	// GetProviderName returns the name of the provider
	GetProviderName() string
	
	// GetSupportedLanguages returns the list of supported language codes
	GetSupportedLanguages() map[string]string
}

// Use config package types
type ProviderConfig = config.Config

// Re-export constants for backward compatibility
const (
	ProviderTypeClaude    = config.ProviderTypeClaude
	ProviderTypeOpenAI    = config.ProviderTypeOpenAI
	ProviderTypeAnthropic = config.ProviderTypeAnthropic
)

// NewLLMProvider creates a new LLM provider based on configuration
func NewLLMProvider(config ProviderConfig) (LLMProvider, error) {
	switch config.ProviderType {
	case ProviderTypeClaude:
		return NewClaudeCodeProvider(config)
	case ProviderTypeOpenAI:
		return NewOpenAIProvider(config)
	case ProviderTypeAnthropic:
		return NewAnthropicProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.ProviderType)
	}
}

// LoadConfig loads provider configuration from config file and environment variables
func LoadConfig() ProviderConfig {
	return config.Load()
}

// LoadConfigFromEnv loads provider configuration from environment variables and .env file (deprecated)
func LoadConfigFromEnv() ProviderConfig {
	return config.LoadFromEnv()
}