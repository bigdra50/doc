package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider implements LLMProvider for OpenAI API
type OpenAIProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	apiKey     string
}

// OpenAI API structures
type openAIRequest struct {
	Model       string             `json:"model"`
	Messages    []openAIMessage    `json:"messages"`
	Tools       []openAITool       `json:"tools,omitempty"`
	ToolChoice  string             `json:"tool_choice,omitempty"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature"`
}

type openAIMessage struct {
	Role      string                 `json:"role"`
	Content   string                 `json:"content,omitempty"`
	ToolCalls []openAIToolCall       `json:"tool_calls,omitempty"`
}

type openAITool struct {
	Type     string            `json:"type"`
	Function openAIFunction    `json:"function"`
}

type openAIFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type openAIToolCall struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Function openAIFunctionCall   `json:"function"`
}

type openAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
	Error   *openAIError   `json:"error,omitempty"`
}

type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}


// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	if config.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	
	provider := &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		apiKey: config.OpenAIAPIKey,
	}
	
	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("openai provider configuration invalid: %w", err)
	}
	
	return provider, nil
}

// ValidateConfig validates the OpenAI provider configuration
func (p *OpenAIProvider) ValidateConfig() error {
	if p.apiKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}
	
	// Skip API validation for now - we'll validate when actually making requests
	// This prevents unnecessary API calls during initialization
	
	return nil
}

// GetProviderName returns the name of the provider
func (p *OpenAIProvider) GetProviderName() string {
	return "OpenAI API"
}

// GetSupportedLanguages returns the list of supported language codes
func (p *OpenAIProvider) GetSupportedLanguages() map[string]string {
	return supportedLanguages
}

// Translate translates the given content using OpenAI API with function calling
func (p *OpenAIProvider) Translate(ctx context.Context, content string, options TranslationOptions) (*TranslationResponse, error) {
	if p.config.Verbose {
		log("Using OpenAI provider for translation")
		log("Target language: %s", options.TargetLanguage)
		if options.CustomInstruction != "" {
			log("Custom instruction: %s", options.CustomInstruction)
		}
	}
	
	// Simply use direct prompting without function calling for translation
	// Function calling is not needed for this use case
	
	// Create the system message and user prompt
	systemPrompt := p.createSystemPrompt()
	userPrompt := p.createUserPrompt(options.TargetLanguage, options.CustomInstruction, content)
	
	// Get model from configuration
	model := p.config.OpenAIModel
	if model == "" {
		model = GetDefaultModel(ProviderTypeOpenAI)
	}
	
	if p.config.Verbose {
		log("Using OpenAI model: %s", model)
	}
	
	// Create the API request without function calling
	req := openAIRequest{
		Model: model,
		Messages: []openAIMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens:   4000,
		Temperature: 0.1,
	}
	
	var response openAIResponse
	if err := p.makeAPIRequest(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}
	
	if p.config.Verbose {
		log("OpenAI API response received with %d choices", len(response.Choices))
	}
	
	// Parse the response
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response choices received from OpenAI")
	}
	
	choice := response.Choices[0]
	
	// Use direct content response (no function calling)
	if choice.Message.Content != "" {
		if p.config.Verbose {
			log("Received translation response of length: %d", len(choice.Message.Content))
		}
		
		return &TranslationResponse{
			Content: choice.Message.Content,
			Status:  "success",
			Message: "Translation completed successfully",
		}, nil
	}
	
	return nil, fmt.Errorf("no content received from OpenAI")
}

// createSystemPrompt creates the system prompt for translation
func (p *OpenAIProvider) createSystemPrompt() string {
	return `You are a professional document translator. Your task is to translate documents while preserving their original format perfectly.

CRITICAL RULES:
1. Preserve ALL original formatting (Markdown, HTML, plain text, etc.) EXACTLY
2. Maintain ALL syntax, tags, symbols, and document structure
3. Do NOT translate code blocks, URLs, or technical identifiers
4. Do NOT change the document structure or format in any way
5. Output ONLY the translated document - no explanations, prefixes, or additional text
6. If the document is already in the target language, return it unchanged

Respond with the translated document only.`
}

// createUserPrompt creates the user prompt for translation
func (p *OpenAIProvider) createUserPrompt(targetLang, customInstruction, content string) string {
	langName := supportedLanguages[targetLang]
	
	prompt := fmt.Sprintf(`Translate the following document to %s (%s).`, langName, targetLang)
	
	if customInstruction != "" {
		prompt += fmt.Sprintf("\n\nAdditional instruction: %s", customInstruction)
	}
	
	prompt += fmt.Sprintf("\n\nDocument to translate:\n%s", content)
	
	return prompt
}

// makeAPIRequest makes an HTTP request to the OpenAI API
func (p *OpenAIProvider) makeAPIRequest(ctx context.Context, req openAIRequest, response interface{}) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	
	if p.config.Verbose {
		log("Making OpenAI API request...")
	}
	
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		var apiError openAIResponse
		if json.Unmarshal(body, &apiError) == nil && apiError.Error != nil {
			return fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, apiError.Error.Message)
		}
		return fmt.Errorf("OpenAI API request failed with status %d", resp.StatusCode)
	}
	
	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}
	
	return nil
}

