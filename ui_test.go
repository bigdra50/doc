package main

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{100 * time.Millisecond, "100ms"},
		{1500 * time.Millisecond, "1.5s"},
		{30 * time.Second, "30.0s"},
		{90 * time.Second, "1.5m"},
		{3 * time.Minute, "3.0m"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %s, expected %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty key", "", "(not set)"},
		{"Short key", "abc", "***"},
		{"Normal key", "sk-1234567890abcdef", "sk-1...cdef"},
		{"Long key", "sk-proj-1234567890abcdefghijklmnopqrstuvwxyz", "sk-p...wxyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskAPIKey(tt.input)
			if result != tt.expected {
				t.Errorf("maskAPIKey(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}
