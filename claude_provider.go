package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ClaudeCodeProvider implements LLMProvider for Claude Code CLI
type ClaudeCodeProvider struct {
	config ProviderConfig
}

// NewClaudeCodeProvider creates a new Claude Code provider
func NewClaudeCodeProvider(config ProviderConfig) (*ClaudeCodeProvider, error) {
	provider := &ClaudeCodeProvider{
		config: config,
	}
	
	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("claude code provider configuration invalid: %w", err)
	}
	
	return provider, nil
}

// ValidateConfig validates the Claude Code provider configuration
func (p *ClaudeCodeProvider) ValidateConfig() error {
	// Check if claude command exists
	claudePath := p.config.ClaudeCodePath
	if claudePath == "" {
		claudePath = "claude"
	}
	
	if _, err := exec.LookPath(claudePath); err != nil {
		return fmt.Errorf("claude command not found at %s: %w", claudePath, err)
	}
	
	return nil
}

// GetProviderName returns the name of the provider
func (p *ClaudeCodeProvider) GetProviderName() string {
	return "Claude Code CLI"
}

// GetSupportedLanguages returns the list of supported language codes
func (p *ClaudeCodeProvider) GetSupportedLanguages() map[string]string {
	return supportedLanguages
}

// Translate translates the given content using Claude Code CLI
func (p *ClaudeCodeProvider) Translate(ctx context.Context, content string, options TranslationOptions) (*TranslationResponse, error) {
	if p.config.Verbose {
		log("Using Claude Code provider for translation")
		log("Target language: %s", options.TargetLanguage)
		if options.CustomInstruction != "" {
			log("Custom instruction: %s", options.CustomInstruction)
		}
	}
	
	// Generate prompt using existing logic
	prompt := p.generatePrompt(options.TargetLanguage, options.CustomInstruction, content)
	
	if p.config.Verbose {
		log("Generated prompt length: %d characters", len(prompt))
		// Save prompt to file for debugging
		if err := os.WriteFile("/tmp/xlat_prompt.txt", []byte(prompt), 0644); err == nil {
			log("Prompt saved to /tmp/xlat_prompt.txt for debugging")
		}
	}
	
	// Execute Claude command
	result, err := p.executeClaude(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("claude command execution failed: %w", err)
	}
	
	if p.config.Verbose {
		log("Claude command executed successfully, output length: %d characters", len(result))
		// Save output to file for debugging
		if err := os.WriteFile("/tmp/xlat_output.txt", []byte(result), 0644); err == nil {
			log("Output saved to /tmp/xlat_output.txt for debugging")
		}
	}
	
	// For now, return simple success response
	// TODO: Add structured response parsing if needed
	response := &TranslationResponse{
		Content: result,
		Status:  "success",
		Message: "Translation completed successfully",
	}
	
	return response, nil
}

// generatePrompt generates the translation prompt (migrated from main.go)
func (p *ClaudeCodeProvider) generatePrompt(targetLang, transformInstruction, content string) string {
	langName := supportedLanguages[targetLang]
	
	prompt := fmt.Sprintf(`Translate the following document to %s (%s).

IMPORTANT:
1. Preserve the original document format (Markdown, HTML, plain text, etc.) EXACTLY
2. Maintain ALL syntax, tags, symbols, and structure  
3. Do NOT translate code blocks, URLs, or technical identifiers
4. Do NOT change the document structure or format in any way
5. Output ONLY the translated document - no explanations, prefixes, or additional text

If the document is already in %s, return it unchanged.`, langName, targetLang, langName)

	if transformInstruction != "" {
		prompt += fmt.Sprintf("\n\nAdditional instruction: %s", transformInstruction)
	}

	prompt += fmt.Sprintf("\n\nDocument:\n%s", content)
	
	return prompt
}

// executeClaude executes the Claude command (migrated from main.go)
func (p *ClaudeCodeProvider) executeClaude(ctx context.Context, prompt string) (string, error) {
	claudePath := p.config.ClaudeCodePath
	if claudePath == "" {
		claudePath = "claude"
	}
	
	modelFlag := p.config.ClaudeModel
	if modelFlag == "" {
		modelFlag = "sonnet"
	}
	
	if p.config.Verbose {
		log("Creating claude command: %s -p --model %s", claudePath, modelFlag)
	}
	
	cmd := exec.CommandContext(ctx, claudePath, "-p", "--model", modelFlag)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stderr = os.Stderr
	
	if p.config.Verbose {
		log("Starting claude command execution...")
	}
	
	output, err := cmd.Output()
	if err != nil {
		if p.config.Verbose {
			log("Claude command failed with error: %v", err)
		}
		return "", fmt.Errorf("claude command execution failed: %w", err)
	}
	
	result := strings.TrimSpace(string(output))
	
	if result == "" {
		return "", fmt.Errorf("claude returned empty response")
	}
	
	return result, nil
}