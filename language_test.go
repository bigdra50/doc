package main

import (
	"testing"
)

func TestValidateLanguageCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErr  bool
	}{
		{"Valid Japanese", "ja", false},
		{"Valid English", "en", false},
		{"Valid Russian", "ru", false},
		{"Valid Chinese", "zh", false},
		{"Invalid code", "xyz", true},
		{"Empty code", "", true},
		{"Case sensitive", "JA", true},
		{"Long code", "japanese", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLanguageCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLanguageCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSupportedLanguagesMap(t *testing.T) {
	// Test that common language codes exist
	requiredCodes := []string{"ja", "en", "ru", "zh", "es", "fr", "de"}
	
	for _, code := range requiredCodes {
		if _, exists := supportedLanguages[code]; !exists {
			t.Errorf("Required language code %s not found in supportedLanguages", code)
		}
	}
	
	// Test that all values are non-empty
	for code, name := range supportedLanguages {
		if name == "" {
			t.Errorf("Language code %s has empty name", code)
		}
		if len(code) < 2 || len(code) > 3 {
			t.Errorf("Language code %s has invalid length (expected 2-3 chars)", code)
		}
	}
}

func TestGetSimilarLanguageCodes(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"j", []string{"ja"}},
		{"e", []string{"el", "en", "es", "et"}},
		{"xyz", []string{}},
		{"", []string{"am", "ar", "bg", "cs", "da", "de", "el", "en", "es", "et", "fi", "fr", "he", "hi", "hr", "hu", "id", "it", "ja", "ko", "lt", "lv", "ms", "mt", "nl", "no", "pl", "pt", "ro", "ru", "sk", "sl", "sv", "sw", "th", "tl", "tr", "vi", "zh"}},
	}

	for _, tt := range tests {
		t.Run("input_"+tt.input, func(t *testing.T) {
			result := getSimilarLanguageCodes(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("getSimilarLanguageCodes(%s) = %v, expected %v", tt.input, result, tt.expected)
				return
			}
			for i, code := range result {
				if code != tt.expected[i] {
					t.Errorf("getSimilarLanguageCodes(%s) = %v, expected %v", tt.input, result, tt.expected)
					break
				}
			}
		})
	}
}