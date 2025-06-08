package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseMergeArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *CLIArgs
		wantErr  bool
	}{
		{
			name: "Basic merge command",
			args: []string{"./docs"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "merged.md",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Merge with output file",
			args: []string{"./docs", "book.md"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "book.md",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Merge with flags",
			args: []string{"./docs", "-r", "--include-meta", "--dry-run"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "merged.md",
				MergeRecursive:    true,
				MergeIncludeMeta:  true,
				MergeDryRun:       true,
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Merge with output option",
			args: []string{"./docs", "-o", "custom.md"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "custom.md",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Merge with order option",
			args: []string{"./docs", "--order", "modified"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "merged.md",
				MergeOrder:        "modified",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Merge with include/exclude patterns",
			args: []string{"./docs", "--include", "*.md", "--exclude", "README.md"},
			expected: &CLIArgs{
				IsMergeCommand:       true,
				MergeDirectory:       "./docs",
				MergeOutputFile:      "merged.md",
				MergeIncludePatterns: []string{"*.md"},
				MergeExcludePatterns: []string{"README.md"},
				MergeOrder:           "filename",
				MergeSeparator:       "\n\n---\n\n",
				MergeGenerateTOC:     true,
				MergeTOCDepth:        3,
				MergeBaseLevel:       2,
				MergeAdjustHeaders:   true,
			},
			wantErr: false,
		},
		{
			name:    "Merge without directory",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Merge with invalid order",
			args:    []string{"./docs", "--order", "invalid"},
			wantErr: true,
		},
		{
			name:    "Merge with invalid toc-depth",
			args:    []string{"./docs", "--toc-depth", "10"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cliArgs := &CLIArgs{
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			}

			result, err := parseMergeArgs(cliArgs, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseMergeArgs() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestParseArgsWithMergeCommand(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected *CLIArgs
		wantErr  bool
	}{
		{
			name: "Parse merge command",
			args: []string{"doc", "merge", "./docs"},
			expected: &CLIArgs{
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "merged.md",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Parse verbose merge command",
			args: []string{"doc", "-v", "merge", "./docs"},
			expected: &CLIArgs{
				Verbose:           true,
				IsMergeCommand:    true,
				MergeDirectory:    "./docs",
				MergeOutputFile:   "merged.md",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
		{
			name: "Parse regular translation command",
			args: []string{"doc", "ja"},
			expected: &CLIArgs{
				TargetLanguage:    "ja",
				MergeOrder:        "filename",
				MergeSeparator:    "\n\n---\n\n",
				MergeGenerateTOC:  true,
				MergeTOCDepth:     3,
				MergeBaseLevel:    2,
				MergeAdjustHeaders: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			result, err := parseArgs()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseArgs() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestIsValidOrder(t *testing.T) {
	tests := []struct {
		name  string
		order string
		want  bool
	}{
		{"filename", "filename", true},
		{"modified", "modified", true},
		{"size", "size", true},
		{"custom", "custom", true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidOrder(tt.order); got != tt.want {
				t.Errorf("isValidOrder(%q) = %v, want %v", tt.order, got, tt.want)
			}
		})
	}
}