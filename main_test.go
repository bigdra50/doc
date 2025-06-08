package main

import (
	"testing"
)

func TestValidateLanguageCodeWithMap(t *testing.T) {
	testMap := map[string]string{
		"ja": "Japanese",
		"en": "English",
		"fr": "French",
	}
	
	tests := []struct {
		name     string
		code     string
		wantErr  bool
	}{
		{"Valid Japanese", "ja", false},
		{"Valid English", "en", false},
		{"Invalid code", "xyz", true},
		{"Empty code", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLanguageCodeWithMap(tt.code, testMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLanguageCodeWithMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSimilarLanguageCodesWithMap(t *testing.T) {
	testMap := map[string]string{
		"ja": "Japanese",
		"je": "Test Language",
		"en": "English",
		"es": "Spanish",
	}
	
	tests := []struct {
		input    string
		expected []string
	}{
		{"j", []string{"ja", "je"}},
		{"e", []string{"en", "es"}},
		{"xyz", []string{}},
	}

	for _, tt := range tests {
		t.Run("input_"+tt.input, func(t *testing.T) {
			result := getSimilarLanguageCodesWithMap(tt.input, testMap)
			if len(result) != len(tt.expected) {
				t.Errorf("getSimilarLanguageCodesWithMap(%s) = %v, expected %v", tt.input, result, tt.expected)
				return
			}
			for i, code := range result {
				if code != tt.expected[i] {
					t.Errorf("getSimilarLanguageCodesWithMap(%s) = %v, expected %v", tt.input, result, tt.expected)
					break
				}
			}
		})
	}
}