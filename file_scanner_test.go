package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestScanMarkdownFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "doc_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"chapter1.md":       "# Chapter 1\nContent 1",
		"chapter2.md":       "# Chapter 2\nContent 2", 
		"README.md":         "# README\nReadme content",
		"notes.txt":         "Not a markdown file",
		"subdir/chapter3.md": "# Chapter 3\nContent 3",
		"subdir/notes.md":   "# Notes\nNotes content",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name      string
		directory string
		recursive bool
		includes  []string
		excludes  []string
		expected  []string
		wantErr   bool
	}{
		{
			name:      "Basic scan",
			directory: tempDir,
			recursive: false,
			expected:  []string{"chapter1.md", "chapter2.md", "README.md"},
		},
		{
			name:      "Recursive scan",
			directory: tempDir,
			recursive: true,
			expected:  []string{"chapter1.md", "chapter2.md", "README.md", "subdir/chapter3.md", "subdir/notes.md"},
		},
		{
			name:      "With exclude pattern",
			directory: tempDir,
			recursive: false,
			excludes:  []string{"README.md"},
			expected:  []string{"chapter1.md", "chapter2.md"},
		},
		{
			name:      "With include pattern",
			directory: tempDir,
			recursive: true,
			includes:  []string{"chapter*.md"},
			expected:  []string{"chapter1.md", "chapter2.md", "subdir/chapter3.md"},
		},
		{
			name:      "Non-existent directory",
			directory: "/non/existent/path",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &FileScanner{
				Directory:       tt.directory,
				Recursive:       tt.recursive,
				IncludePatterns: tt.includes,
				ExcludePatterns: tt.excludes,
			}

			files, err := scanner.ScanMarkdownFiles()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Convert absolute paths to relative paths for comparison
			relativePaths := make([]string, len(files))
			for i, file := range files {
				relPath, err := filepath.Rel(tempDir, file.Path)
				if err != nil {
					t.Fatalf("Failed to get relative path: %v", err)
				}
				relativePaths[i] = relPath
			}

			sort.Strings(relativePaths)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(relativePaths, tt.expected) {
				t.Errorf("ScanMarkdownFiles() = %v, want %v", relativePaths, tt.expected)
			}
		})
	}
}

func TestSortMarkdownFiles(t *testing.T) {
	// Create temporary files with different times
	tempDir, err := os.MkdirTemp("", "doc_sort_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with different sizes and modification times
	now := time.Now()
	testFiles := []struct {
		name    string
		content string
		modTime time.Time
	}{
		{"z_file.md", "small", now.Add(-2 * time.Hour)},
		{"a_file.md", "this is a much larger file content to test size sorting", now.Add(-1 * time.Hour)},
		{"m_file.md", "medium", now},
	}

	var markdownFiles []MarkdownFile
	for _, tf := range testFiles {
		path := filepath.Join(tempDir, tf.name)
		if err := os.WriteFile(path, []byte(tf.content), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, tf.modTime, tf.modTime); err != nil {
			t.Fatal(err)
		}

		stat, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}

		markdownFiles = append(markdownFiles, MarkdownFile{
			Path:    path,
			Name:    tf.name,
			ModTime: stat.ModTime(),
			Size:    stat.Size(),
		})
	}

	tests := []struct {
		name     string
		sortType string
		expected []string
	}{
		{
			name:     "Sort by filename",
			sortType: "filename",
			expected: []string{"a_file.md", "m_file.md", "z_file.md"},
		},
		{
			name:     "Sort by modified time",
			sortType: "modified",
			expected: []string{"z_file.md", "a_file.md", "m_file.md"},
		},
		{
			name:     "Sort by size",
			sortType: "size",
			expected: []string{"z_file.md", "m_file.md", "a_file.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := SortMarkdownFiles(markdownFiles, tt.sortType)
			
			result := make([]string, len(sorted))
			for i, file := range sorted {
				result[i] = file.Name
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortMarkdownFiles(%s) = %v, want %v", tt.sortType, result, tt.expected)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		pattern  string
		expected bool
	}{
		{"Exact match", "file.md", "file.md", true},
		{"Wildcard match", "chapter1.md", "chapter*.md", true},
		{"No match", "readme.txt", "*.md", false},
		{"Question mark match", "file1.md", "file?.md", true},
		{"Complex pattern", "chapter01.md", "chapter[0-9][0-9].md", true},
		{"No complex match", "chapterAB.md", "chapter[0-9][0-9].md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPattern(tt.filename, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.filename, tt.pattern, result, tt.expected)
			}
		})
	}
}