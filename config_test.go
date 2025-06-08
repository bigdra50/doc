package main

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original env vars
	originalProvider := os.Getenv("LLM_PROVIDER")
	originalOpenAI := os.Getenv("OPENAI_API_KEY")
	
	// Temporarily rename .env file to avoid interference
	envExists := false
	if _, err := os.Stat(".env"); err == nil {
		envExists = true
		os.Rename(".env", ".env.backup")
	}
	
	// Clean up after test
	defer func() {
		if envExists {
			os.Rename(".env.backup", ".env")
		}
		os.Setenv("LLM_PROVIDER", originalProvider)
		os.Setenv("OPENAI_API_KEY", originalOpenAI)
	}()
	
	// Test default config
	os.Unsetenv("LLM_PROVIDER")
	os.Unsetenv("OPENAI_API_KEY")
	
	config := LoadConfigFromEnv()
	if config.ProviderType != ProviderTypeClaude {
		t.Errorf("Expected default provider %s, got %s", ProviderTypeClaude, config.ProviderType)
	}
	
	// Test custom config
	os.Setenv("LLM_PROVIDER", "openai")
	os.Setenv("OPENAI_API_KEY", "test-key")
	
	config = LoadConfigFromEnv()
	if config.ProviderType != ProviderTypeOpenAI {
		t.Errorf("Expected provider %s, got %s", ProviderTypeOpenAI, config.ProviderType)
	}
	if config.OpenAIAPIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", config.OpenAIAPIKey)
	}
}