package main

import "github.com/bigdra50/doc/internal/config"

// Model represents an LLM model with its characteristics
type Model struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	InputCostPer1M  float64  `json:"input_cost_per_1m"`
	OutputCostPer1M float64  `json:"output_cost_per_1m"`
	ContextWindow   int      `json:"context_window"`
	Tier            string   `json:"tier"`
	RecommendedFor  []string `json:"recommended_for"`
}

// ModelCatalog holds all available models by provider
type ModelCatalog struct {
	OpenAI    []Model `json:"openai"`
	Anthropic []Model `json:"anthropic"`
}

// GetModelCatalog returns the complete model catalog
func GetModelCatalog() ModelCatalog {
	return ModelCatalog{
		OpenAI: []Model{
			{
				ID:              "gpt-4",
				Name:            "GPT-4",
				InputCostPer1M:  30.00,
				OutputCostPer1M: 60.00,
				ContextWindow:   8192,
				Tier:            "premium",
				RecommendedFor:  []string{"complex_reasoning", "code_generation"},
			},
			{
				ID:              "gpt-4-turbo",
				Name:            "GPT-4 Turbo",
				InputCostPer1M:  10.00,
				OutputCostPer1M: 30.00,
				ContextWindow:   128000,
				Tier:            "balanced",
				RecommendedFor:  []string{"general_translation", "balanced_performance"},
			},
			{
				ID:              "gpt-4o",
				Name:            "GPT-4o",
				InputCostPer1M:  2.50,
				OutputCostPer1M: 10.00,
				ContextWindow:   128000,
				Tier:            "balanced",
				RecommendedFor:  []string{"document_with_images", "complex_formatting"},
			},
			{
				ID:              "gpt-4o-mini",
				Name:            "GPT-4o Mini",
				InputCostPer1M:  0.15,
				OutputCostPer1M: 0.60,
				ContextWindow:   128000,
				Tier:            "economy",
				RecommendedFor:  []string{"simple_translation", "high_volume"},
			},
			{
				ID:              "gpt-3.5-turbo",
				Name:            "GPT-3.5 Turbo",
				InputCostPer1M:  0.50,
				OutputCostPer1M: 1.50,
				ContextWindow:   16000,
				Tier:            "economy",
				RecommendedFor:  []string{"budget_translation"},
			},
		},
		Anthropic: []Model{
			{
				ID:              "claude-3-opus-20240229",
				Name:            "Claude 3 Opus",
				InputCostPer1M:  15.00,
				OutputCostPer1M: 75.00,
				ContextWindow:   200000,
				Tier:            "premium",
				RecommendedFor:  []string{"complex_reasoning", "code_generation"},
			},
			{
				ID:              "claude-3-sonnet-20240229",
				Name:            "Claude 3 Sonnet",
				InputCostPer1M:  3.00,
				OutputCostPer1M: 15.00,
				ContextWindow:   200000,
				Tier:            "balanced",
				RecommendedFor:  []string{"general_translation", "balanced_performance"},
			},
			{
				ID:              "claude-3-5-sonnet-20241022",
				Name:            "Claude 3.5 Sonnet",
				InputCostPer1M:  3.00,
				OutputCostPer1M: 15.00,
				ContextWindow:   200000,
				Tier:            "balanced",
				RecommendedFor:  []string{"general_translation", "advanced_reasoning"},
			},
			{
				ID:              "claude-3-haiku-20240307",
				Name:            "Claude 3 Haiku",
				InputCostPer1M:  0.25,
				OutputCostPer1M: 1.25,
				ContextWindow:   200000,
				Tier:            "economy",
				RecommendedFor:  []string{"simple_translation", "high_volume"},
			},
			{
				ID:              "claude-3-5-haiku-20241022",
				Name:            "Claude 3.5 Haiku",
				InputCostPer1M:  0.80,
				OutputCostPer1M: 4.00,
				ContextWindow:   200000,
				Tier:            "economy",
				RecommendedFor:  []string{"simple_translation", "high_volume"},
			},
		},
	}
}

// GetModelsByProvider returns models for a specific provider
func GetModelsByProvider(provider string) []Model {
	catalog := GetModelCatalog()
	switch provider {
	case ProviderTypeOpenAI:
		return catalog.OpenAI
	case ProviderTypeAnthropic:
		return catalog.Anthropic
	default:
		return []Model{}
	}
}

// FindModel finds a model by ID within a provider
func FindModel(provider, modelID string) *Model {
	models := GetModelsByProvider(provider)
	for _, model := range models {
		if model.ID == modelID {
			return &model
		}
	}
	return nil
}

// GetDefaultModel returns the default model for a provider
// Delegates to config package to avoid duplication
func GetDefaultModel(provider string) string {
	// Use the config package's implementation
	return config.GetDefaultModel(provider)
}

// GetModelsByTier returns models filtered by tier
func GetModelsByTier(provider, tier string) []Model {
	models := GetModelsByProvider(provider)
	var filtered []Model
	for _, model := range models {
		if model.Tier == tier {
			filtered = append(filtered, model)
		}
	}
	return filtered
}

// EstimateCost estimates the cost for a translation request
func EstimateCost(model Model, inputLength, outputLength int) float64 {
	// Rough estimation: 1 token ≈ 4 characters
	inputTokens := float64(inputLength) / 4.0
	outputTokens := float64(outputLength) / 4.0
	
	inputCost := (inputTokens / 1000000.0) * model.InputCostPer1M
	outputCost := (outputTokens / 1000000.0) * model.OutputCostPer1M
	
	return inputCost + outputCost
}